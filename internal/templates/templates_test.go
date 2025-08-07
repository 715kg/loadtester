package templates

import (
	"strings"
	"testing"
)

func TestHTMLTemplate(t *testing.T) {
	if len(HTMLTemplate) == 0 {
		t.Error("HTMLTemplate should not be empty")
	}

	// Check for essential HTML structure
	requiredElements := []string{
		"<!DOCTYPE html>",
		"<html",
		"<head>",
		"<body>",
		"</html>",
	}

	for _, element := range requiredElements {
		if !strings.Contains(HTMLTemplate, element) {
			t.Errorf("HTMLTemplate should contain %q", element)
		}
	}
}

func TestHTMLTemplateTitle(t *testing.T) {
	if !strings.Contains(HTMLTemplate, "<title>Load Tester</title>") {
		t.Error("HTMLTemplate should contain proper title")
	}
}

func TestHTMLTemplateNavigation(t *testing.T) {
	navigationElements := []string{
		"Главная",
		"Инструкция",
		"Соглашение",
	}

	for _, element := range navigationElements {
		if !strings.Contains(HTMLTemplate, element) {
			t.Errorf("HTMLTemplate should contain navigation element %q", element)
		}
	}
}

func TestHTMLTemplateFormElements(t *testing.T) {
	formElements := []string{
		"targetURL",
		"method",
		"totalRequests",
		"maxConcurrency",
		"timeout",
		"maxMemory",
	}

	for _, element := range formElements {
		if !strings.Contains(HTMLTemplate, element) {
			t.Errorf("HTMLTemplate should contain form element %q", element)
		}
	}
}

func TestHTMLTemplateJavaScript(t *testing.T) {
	jsElements := []string{
		"startTest()",
		"stopTest()",
		"updateStats()",
		"toggleTheme()",
		"showPage(",
	}

	for _, element := range jsElements {
		if !strings.Contains(HTMLTemplate, element) {
			t.Errorf("HTMLTemplate should contain JavaScript function %q", element)
		}
	}
}

func TestHTMLTemplateCSS(t *testing.T) {
	cssElements := []string{
		":root",
		"--bg-primary",
		"--accent-color",
		".nav-container",
		".main-card",
		".form-group",
	}

	for _, element := range cssElements {
		if !strings.Contains(HTMLTemplate, element) {
			t.Errorf("HTMLTemplate should contain CSS element %q", element)
		}
	}
}

func TestHTMLTemplateThemeSupport(t *testing.T) {
	themeElements := []string{
		"[data-theme=\"dark\"]",
		"toggleTheme()",
		"theme-toggle",
	}

	for _, element := range themeElements {
		if !strings.Contains(HTMLTemplate, element) {
			t.Errorf("HTMLTemplate should contain theme element %q", element)
		}
	}
}

func TestHTMLTemplateVersionPlaceholder(t *testing.T) {
	if !strings.Contains(HTMLTemplate, "{{.Version}}") {
		t.Error("HTMLTemplate should contain version placeholder {{.Version}}")
	}
}

func TestHTMLTemplateStatistics(t *testing.T) {
	statsElements := []string{
		"totalStat",
		"successStat",
		"failsStat",
		"errorsStat",
		"rpsStat",
	}

	for _, element := range statsElements {
		if !strings.Contains(HTMLTemplate, element) {
			t.Errorf("HTMLTemplate should contain statistics element %q", element)
		}
	}
}

func TestHTMLTemplateAgreement(t *testing.T) {
	agreementElements := []string{
		"agreementCheckbox",
		"пользовательское соглашение",
		"Принимаю",
	}

	for _, element := range agreementElements {
		if !strings.Contains(HTMLTemplate, element) {
			t.Errorf("HTMLTemplate should contain agreement element %q", element)
		}
	}
}
