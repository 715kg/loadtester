package dnstest

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type DNSTester struct {
	Domain      string
	RecordTypes []string
	DNSServers  []string
	Timeout     time.Duration
	Results     *DNSResults
	isActive    bool
	stopChan    chan bool
	mu          sync.RWMutex
}

type DNSResults struct {
	Domain      string          `json:"domain"`
	IsRunning   bool            `json:"isRunning"`
	IsCompleted bool            `json:"isCompleted"`
	StartTime   time.Time       `json:"startTime"`
	EndTime     time.Time       `json:"endTime"`
	TotalTime   time.Duration   `json:"totalTime"`
	Error       string          `json:"error,omitempty"`
	Records     []DNSRecord     `json:"records"`
	ServerStats []DNSServerStat `json:"serverStats"`
	Summary     DNSSummary      `json:"summary"`
}

type DNSRecord struct {
	Type         string        `json:"type"`
	Name         string        `json:"name"`
	Value        string        `json:"value"`
	TTL          uint32        `json:"ttl"`
	Server       string        `json:"server"`
	ResponseTime time.Duration `json:"responseTime"`
	Success      bool          `json:"success"`
	Error        string        `json:"error,omitempty"`
}

type DNSServerStat struct {
	Server       string        `json:"server"`
	TotalQueries int           `json:"totalQueries"`
	Successful   int           `json:"successful"`
	Failed       int           `json:"failed"`
	AvgTime      time.Duration `json:"avgTime"`
	MinTime      time.Duration `json:"minTime"`
	MaxTime      time.Duration `json:"maxTime"`
	Available    bool          `json:"available"`
}

type DNSSummary struct {
	TotalRecords      int     `json:"totalRecords"`
	SuccessfulLookups int     `json:"successfulLookups"`
	FailedLookups     int     `json:"failedLookups"`
	SuccessRate       float64 `json:"successRate"`
	FastestServer     string  `json:"fastestServer"`
	SlowestServer     string  `json:"slowestServer"`
}

func NewDNSTester() *DNSTester {
	return &DNSTester{
		RecordTypes: []string{"A", "AAAA", "MX", "CNAME", "TXT", "NS"},
		DNSServers:  []string{"8.8.8.8", "1.1.1.1", "208.67.222.222"},
		Timeout:     5 * time.Second,
		Results:     &DNSResults{},
		stopChan:    make(chan bool, 1),
	}
}

func (d *DNSTester) Start(domain string, recordTypes []string, dnsServers []string, timeout int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.isActive {
		return fmt.Errorf("DNS тест уже выполняется")
	}

	// Валидация входных данных
	if domain == "" {
		return fmt.Errorf("не указан домен для проверки")
	}

	d.Domain = domain
	if len(recordTypes) > 0 {
		d.RecordTypes = recordTypes
	}
	if len(dnsServers) > 0 {
		d.DNSServers = dnsServers
	}
	d.Timeout = time.Duration(timeout) * time.Second
	d.isActive = true
	d.stopChan = make(chan bool, 1)

	// Инициализируем результаты
	d.Results = &DNSResults{
		Domain:      domain,
		IsRunning:   true,
		IsCompleted: false,
		StartTime:   time.Now(),
		Records:     make([]DNSRecord, 0),
		ServerStats: make([]DNSServerStat, 0),
	}

	// Запускаем DNS тест в горутине
	go d.runDNSTest()

	return nil
}

func (d *DNSTester) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.isActive {
		d.isActive = false
		select {
		case d.stopChan <- true:
		default:
		}
	}
}

func (d *DNSTester) IsActive() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.isActive
}

func (d *DNSTester) GetResults() *DNSResults {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.Results
}

func (d *DNSTester) runDNSTest() {
	defer func() {
		d.mu.Lock()
		d.isActive = false
		d.Results.IsRunning = false
		d.Results.IsCompleted = true
		d.Results.EndTime = time.Now()
		d.Results.TotalTime = d.Results.EndTime.Sub(d.Results.StartTime)
		d.calculateSummary()
		d.mu.Unlock()
	}()

	// Инициализируем статистику серверов
	serverStats := make(map[string]*DNSServerStat)
	for _, server := range d.DNSServers {
		serverStats[server] = &DNSServerStat{
			Server:    server,
			Available: true,
		}
	}

	// Тестируем каждый тип записи на каждом DNS сервере
	for _, recordType := range d.RecordTypes {
		for _, server := range d.DNSServers {
			// Проверяем, нужно ли остановиться
			select {
			case <-d.stopChan:
				return
			default:
			}

			record := d.queryDNS(d.Domain, recordType, server)

			d.mu.Lock()
			d.Results.Records = append(d.Results.Records, record)

			// Обновляем статистику сервера
			stat := serverStats[server]
			stat.TotalQueries++

			if record.Success {
				stat.Successful++
				if stat.MinTime == 0 || record.ResponseTime < stat.MinTime {
					stat.MinTime = record.ResponseTime
				}
				if record.ResponseTime > stat.MaxTime {
					stat.MaxTime = record.ResponseTime
				}
			} else {
				stat.Failed++
				if strings.Contains(record.Error, "timeout") || strings.Contains(record.Error, "no such host") {
					stat.Available = false
				}
			}
			d.mu.Unlock()
		}
	}

	// Вычисляем среднее время для каждого сервера
	d.mu.Lock()
	for _, stat := range serverStats {
		if stat.Successful > 0 {
			var totalTime time.Duration
			for _, record := range d.Results.Records {
				if record.Server == stat.Server && record.Success {
					totalTime += record.ResponseTime
				}
			}
			stat.AvgTime = totalTime / time.Duration(stat.Successful)
		}
		d.Results.ServerStats = append(d.Results.ServerStats, *stat)
	}
	d.mu.Unlock()
}

func (d *DNSTester) queryDNS(domain, recordType, server string) DNSRecord {
	start := time.Now()

	record := DNSRecord{
		Type:   recordType,
		Name:   domain,
		Server: server,
	}

	// Создаем кастомный резолвер
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: d.Timeout,
			}
			return d.DialContext(ctx, network, server+":53")
		},
	}

	ctx := context.Background()

	switch recordType {
	case "A":
		ips, err := resolver.LookupIPAddr(ctx, domain)
		if err != nil {
			record.Error = err.Error()
		} else {
			var ipv4s []string
			for _, ip := range ips {
				if ip.IP.To4() != nil {
					ipv4s = append(ipv4s, ip.IP.String())
				}
			}
			if len(ipv4s) > 0 {
				record.Value = strings.Join(ipv4s, ", ")
				record.Success = true
			} else {
				record.Error = "No A records found"
			}
		}

	case "AAAA":
		ips, err := resolver.LookupIPAddr(ctx, domain)
		if err != nil {
			record.Error = err.Error()
		} else {
			var ipv6s []string
			for _, ip := range ips {
				if ip.IP.To4() == nil {
					ipv6s = append(ipv6s, ip.IP.String())
				}
			}
			if len(ipv6s) > 0 {
				record.Value = strings.Join(ipv6s, ", ")
				record.Success = true
			} else {
				record.Error = "No AAAA records found"
			}
		}

	case "MX":
		mxs, err := resolver.LookupMX(ctx, domain)
		if err != nil {
			record.Error = err.Error()
		} else if len(mxs) > 0 {
			var mxList []string
			for _, mx := range mxs {
				mxList = append(mxList, fmt.Sprintf("%d %s", mx.Pref, mx.Host))
			}
			record.Value = strings.Join(mxList, ", ")
			record.Success = true
		} else {
			record.Error = "No MX records found"
		}

	case "CNAME":
		cname, err := resolver.LookupCNAME(ctx, domain)
		if err != nil {
			record.Error = err.Error()
		} else {
			record.Value = cname
			record.Success = true
		}

	case "TXT":
		txts, err := resolver.LookupTXT(ctx, domain)
		if err != nil {
			record.Error = err.Error()
		} else if len(txts) > 0 {
			record.Value = strings.Join(txts, " | ")
			record.Success = true
		} else {
			record.Error = "No TXT records found"
		}

	case "NS":
		nss, err := resolver.LookupNS(ctx, domain)
		if err != nil {
			record.Error = err.Error()
		} else if len(nss) > 0 {
			var nsList []string
			for _, ns := range nss {
				nsList = append(nsList, ns.Host)
			}
			record.Value = strings.Join(nsList, ", ")
			record.Success = true
		} else {
			record.Error = "No NS records found"
		}

	default:
		record.Error = "Unsupported record type"
	}

	record.ResponseTime = time.Since(start)
	return record
}

func (d *DNSTester) calculateSummary() {
	summary := DNSSummary{}

	summary.TotalRecords = len(d.Results.Records)

	for _, record := range d.Results.Records {
		if record.Success {
			summary.SuccessfulLookups++
		} else {
			summary.FailedLookups++
		}
	}

	if summary.TotalRecords > 0 {
		summary.SuccessRate = float64(summary.SuccessfulLookups) / float64(summary.TotalRecords) * 100
	}

	// Находим самый быстрый и медленный сервер
	var fastestTime, slowestTime time.Duration
	for _, stat := range d.Results.ServerStats {
		if stat.Successful > 0 {
			if summary.FastestServer == "" || stat.AvgTime < fastestTime {
				summary.FastestServer = stat.Server
				fastestTime = stat.AvgTime
			}
			if summary.SlowestServer == "" || stat.AvgTime > slowestTime {
				summary.SlowestServer = stat.Server
				slowestTime = stat.AvgTime
			}
		}
	}

	d.Results.Summary = summary
}
