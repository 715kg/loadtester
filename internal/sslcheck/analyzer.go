package sslcheck

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type SSLAnalyzer struct {
	Target   string
	Port     int
	Timeout  time.Duration
	Results  *SSLResults
	isActive bool
	mu       sync.RWMutex
}

type SSLResults struct {
	Target          string             `json:"target"`
	Port            int                `json:"port"`
	IsSecure        bool               `json:"isSecure"`
	Certificate     *CertificateInfo   `json:"certificate,omitempty"`
	Chain           []*CertificateInfo `json:"chain,omitempty"`
	TLSVersion      string             `json:"tlsVersion"`
	CipherSuite     string             `json:"cipherSuite"`
	Protocols       []string           `json:"protocols"`
	SecurityGrade   string             `json:"securityGrade"`
	Vulnerabilities []string           `json:"vulnerabilities"`
	CheckTime       time.Time          `json:"checkTime"`
	ResponseTime    time.Duration      `json:"responseTime"`
	Error           string             `json:"error,omitempty"`
}

type CertificateInfo struct {
	Subject            string    `json:"subject"`
	Issuer             string    `json:"issuer"`
	CommonName         string    `json:"commonName"`
	SANs               []string  `json:"sans"`
	NotBefore          time.Time `json:"notBefore"`
	NotAfter           time.Time `json:"notAfter"`
	DaysUntilExpiry    int       `json:"daysUntilExpiry"`
	IsExpired          bool      `json:"isExpired"`
	IsExpiringSoon     bool      `json:"isExpiringSoon"`
	SignatureAlgorithm string    `json:"signatureAlgorithm"`
	KeySize            int       `json:"keySize"`
	SerialNumber       string    `json:"serialNumber"`
	Fingerprint        string    `json:"fingerprint"`
}

func NewSSLAnalyzer() *SSLAnalyzer {
	return &SSLAnalyzer{
		Port:    443,
		Timeout: 10 * time.Second,
		Results: &SSLResults{},
	}
}

func (s *SSLAnalyzer) Start(target string, port int, timeout int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isActive {
		return fmt.Errorf("SSL анализ уже выполняется")
	}

	s.Target = target
	s.Port = port
	s.Timeout = time.Duration(timeout) * time.Second
	s.isActive = true

	// Запускаем анализ в горутине
	go s.analyze()

	return nil
}

func (s *SSLAnalyzer) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isActive = false
}

func (s *SSLAnalyzer) IsActive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isActive
}

func (s *SSLAnalyzer) GetResults() *SSLResults {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Results
}

func (s *SSLAnalyzer) analyze() {
	defer func() {
		s.mu.Lock()
		s.isActive = false
		s.mu.Unlock()
	}()

	startTime := time.Now()

	s.Results = &SSLResults{
		Target:    s.Target,
		Port:      s.Port,
		CheckTime: startTime,
	}

	// Подготавливаем адрес
	address := fmt.Sprintf("%s:%d", s.Target, s.Port)

	// Настраиваем TLS конфигурацию
	config := &tls.Config{
		ServerName:         s.Target,
		InsecureSkipVerify: true, // Нам нужно проанализировать даже невалидные сертификаты
	}

	// Устанавливаем соединение с таймаутом
	dialer := &net.Dialer{
		Timeout: s.Timeout,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", address, config)
	if err != nil {
		s.Results.Error = fmt.Sprintf("Ошибка подключения: %v", err)
		s.Results.IsSecure = false
		return
	}
	defer conn.Close()

	s.Results.ResponseTime = time.Since(startTime)
	s.Results.IsSecure = true

	// Получаем информацию о TLS соединении
	state := conn.ConnectionState()
	s.Results.TLSVersion = s.getTLSVersionString(state.Version)
	s.Results.CipherSuite = s.getCipherSuiteString(state.CipherSuite)

	// Анализируем сертификаты
	if len(state.PeerCertificates) > 0 {
		// Основной сертификат
		cert := state.PeerCertificates[0]
		s.Results.Certificate = s.analyzeCertificate(cert)

		// Цепочка сертификатов
		for _, cert := range state.PeerCertificates {
			certInfo := s.analyzeCertificate(cert)
			s.Results.Chain = append(s.Results.Chain, certInfo)
		}
	}

	// Проверяем поддерживаемые протоколы
	s.Results.Protocols = s.checkSupportedProtocols()

	// Оцениваем безопасность
	s.Results.SecurityGrade = s.calculateSecurityGrade()

	// Проверяем уязвимости
	s.Results.Vulnerabilities = s.checkVulnerabilities()
}

func (s *SSLAnalyzer) analyzeCertificate(cert *x509.Certificate) *CertificateInfo {
	now := time.Now()
	daysUntilExpiry := int(cert.NotAfter.Sub(now).Hours() / 24)

	certInfo := &CertificateInfo{
		Subject:            cert.Subject.String(),
		Issuer:             cert.Issuer.String(),
		CommonName:         cert.Subject.CommonName,
		SANs:               cert.DNSNames,
		NotBefore:          cert.NotBefore,
		NotAfter:           cert.NotAfter,
		DaysUntilExpiry:    daysUntilExpiry,
		IsExpired:          now.After(cert.NotAfter),
		IsExpiringSoon:     daysUntilExpiry <= 30 && daysUntilExpiry > 0,
		SignatureAlgorithm: cert.SignatureAlgorithm.String(),
		SerialNumber:       cert.SerialNumber.String(),
	}

	// Определяем размер ключа
	switch pub := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		certInfo.KeySize = pub.N.BitLen()
	case *ecdsa.PublicKey:
		certInfo.KeySize = pub.Curve.Params().BitSize
	}

	// Вычисляем отпечаток
	certInfo.Fingerprint = fmt.Sprintf("%x", cert.Raw)

	return certInfo
}

func (s *SSLAnalyzer) getTLSVersionString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("Unknown (0x%04x)", version)
	}
}

func (s *SSLAnalyzer) getCipherSuiteString(suite uint16) string {
	switch suite {
	case tls.TLS_RSA_WITH_RC4_128_SHA:
		return "TLS_RSA_WITH_RC4_128_SHA"
	case tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA:
		return "TLS_RSA_WITH_3DES_EDE_CBC_SHA"
	case tls.TLS_RSA_WITH_AES_128_CBC_SHA:
		return "TLS_RSA_WITH_AES_128_CBC_SHA"
	case tls.TLS_RSA_WITH_AES_256_CBC_SHA:
		return "TLS_RSA_WITH_AES_256_CBC_SHA"
	case tls.TLS_RSA_WITH_AES_128_CBC_SHA256:
		return "TLS_RSA_WITH_AES_128_CBC_SHA256"
	case tls.TLS_RSA_WITH_AES_128_GCM_SHA256:
		return "TLS_RSA_WITH_AES_128_GCM_SHA256"
	case tls.TLS_RSA_WITH_AES_256_GCM_SHA384:
		return "TLS_RSA_WITH_AES_256_GCM_SHA384"
	case tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA:
		return "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA"
	case tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA:
		return "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA"
	case tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA:
		return "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA"
	case tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA:
		return "TLS_ECDHE_RSA_WITH_RC4_128_SHA"
	case tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA:
		return "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA"
	case tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA:
		return "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA"
	case tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA:
		return "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA"
	case tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256:
		return "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256"
	case tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256:
		return "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256"
	case tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256:
		return "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"
	case tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256:
		return "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"
	case tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384:
		return "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
	case tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384:
		return "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384"
	case tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256:
		return "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256"
	case tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256:
		return "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256"
	default:
		return fmt.Sprintf("Unknown (0x%04x)", suite)
	}
}

func (s *SSLAnalyzer) checkSupportedProtocols() []string {
	protocols := []string{}

	// Проверяем различные версии TLS
	versions := []struct {
		name    string
		version uint16
	}{
		{"TLS 1.0", tls.VersionTLS10},
		{"TLS 1.1", tls.VersionTLS11},
		{"TLS 1.2", tls.VersionTLS12},
		{"TLS 1.3", tls.VersionTLS13},
	}

	for _, v := range versions {
		if s.testTLSVersion(v.version) {
			protocols = append(protocols, v.name)
		}
	}

	return protocols
}

func (s *SSLAnalyzer) testTLSVersion(version uint16) bool {
	address := fmt.Sprintf("%s:%d", s.Target, s.Port)

	config := &tls.Config{
		ServerName:         s.Target,
		InsecureSkipVerify: true,
		MinVersion:         version,
		MaxVersion:         version,
	}

	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", address, config)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}

func (s *SSLAnalyzer) calculateSecurityGrade() string {
	if s.Results.Certificate == nil {
		return "F"
	}

	score := 100

	// Проверяем версию TLS
	if strings.Contains(s.Results.TLSVersion, "1.0") || strings.Contains(s.Results.TLSVersion, "1.1") {
		score -= 20
	}

	// Проверяем сертификат
	if s.Results.Certificate.IsExpired {
		score -= 50
	} else if s.Results.Certificate.IsExpiringSoon {
		score -= 10
	}

	// Проверяем размер ключа
	if s.Results.Certificate.KeySize < 2048 {
		score -= 30
	}

	// Проверяем алгоритм подписи
	if strings.Contains(strings.ToLower(s.Results.Certificate.SignatureAlgorithm), "sha1") {
		score -= 20
	}

	// Определяем оценку
	if score >= 90 {
		return "A+"
	} else if score >= 80 {
		return "A"
	} else if score >= 70 {
		return "B"
	} else if score >= 60 {
		return "C"
	} else if score >= 50 {
		return "D"
	} else {
		return "F"
	}
}

func (s *SSLAnalyzer) checkVulnerabilities() []string {
	vulnerabilities := []string{}

	// Проверяем устаревшие протоколы
	for _, protocol := range s.Results.Protocols {
		if protocol == "TLS 1.0" {
			vulnerabilities = append(vulnerabilities, "Поддержка устаревшего TLS 1.0")
		}
		if protocol == "TLS 1.1" {
			vulnerabilities = append(vulnerabilities, "Поддержка устаревшего TLS 1.1")
		}
	}

	// Проверяем сертификат
	if s.Results.Certificate != nil {
		if s.Results.Certificate.IsExpired {
			vulnerabilities = append(vulnerabilities, "Сертификат истек")
		}
		if s.Results.Certificate.IsExpiringSoon {
			vulnerabilities = append(vulnerabilities, "Сертификат скоро истечет")
		}
		if s.Results.Certificate.KeySize < 2048 {
			vulnerabilities = append(vulnerabilities, "Слабый ключ (< 2048 бит)")
		}
		if strings.Contains(strings.ToLower(s.Results.Certificate.SignatureAlgorithm), "sha1") {
			vulnerabilities = append(vulnerabilities, "Использование устаревшего SHA-1")
		}
	}

	// Проверяем слабые шифры
	weakCiphers := []string{"RC4", "3DES", "DES"}
	for _, weak := range weakCiphers {
		if strings.Contains(s.Results.CipherSuite, weak) {
			vulnerabilities = append(vulnerabilities, fmt.Sprintf("Слабый шифр: %s", weak))
		}
	}

	return vulnerabilities
}
