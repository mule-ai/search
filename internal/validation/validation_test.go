package validation

import (
	"testing"
)

func TestValidateQuery(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{"valid query", "golang tutorials", false},
		{"valid query with spaces", "  machine learning  ", false},
		{"empty query", "", true},
		{"whitespace only", "   ", true},
		{"too long query", string(make([]byte, 1001)), true},
		{"exactly max length", string(make([]byte, 1000)) + "a", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuery(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateResultCount(t *testing.T) {
	tests := []struct {
		name    string
		count   int
		wantErr bool
	}{
		{"minimum valid", 1, false},
		{"valid count", 10, false},
		{"maximum valid", 100, false},
		{"zero", 0, true},
		{"negative", -1, true},
		{"exceeds maximum", 101, true},
		{"large number", 1000, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateResultCount(tt.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateResultCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout int
		wantErr bool
	}{
		{"minimum valid", 1, false},
		{"valid timeout", 30, false},
		{"maximum valid", 300, false},
		{"zero", 0, true},
		{"negative", -1, true},
		{"exceeds maximum", 301, true},
		{"large number", 600, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTimeout(tt.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTimeout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateInstanceURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid https URL", "https://search.butler.ooo", false},
		{"valid http URL - localhost allowed", "http://localhost:8080", false},
		{"valid http URL - 127.0.0.1 allowed", "http://127.0.0.1:8080", false},
		{"valid http URL - ::1 allowed", "http://[::1]:8080", false},
		{"http URL rejected for non-localhost", "http://searx.example.com", true},
		{"valid URL with path", "https://search.example.com/search", false},
		{"valid URL with port", "https://search.example.com:8080", false},
		{"empty string", "", true},
		{"no scheme", "search.example.com", true},
		{"invalid scheme", "ftp://search.example.com", true},
		{"no host", "https://", true},
		{"just path", "/search", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInstanceURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInstanceURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFormat(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"text format", "text", false},
		{"json format", "json", false},
		{"markdown format", "markdown", false},
		{"uppercase JSON", "JSON", false},
		{"mixed case TEXT", "TeXt", false},
		{"with spaces", " json ", false},
		{"invalid format", "xml", true},
		{"empty string", "", true},
		{"csv format", "csv", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFormat(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSafeSearch(t *testing.T) {
	tests := []struct {
		name       string
		safeSearch int
		wantErr    bool
	}{
		{"level 0 - off", 0, false},
		{"level 1 - moderate", 1, false},
		{"level 2 - strict", 2, false},
		{"negative level", -1, true},
		{"level 3", 3, true},
		{"large number", 10, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSafeSearch(tt.safeSearch)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSafeSearch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateLanguage(t *testing.T) {
	tests := []struct {
		name    string
		language string
		wantErr bool
	}{
		{"valid 2-letter code", "en", false},
		{"valid 2-letter code", "de", false},
		{"valid 5-letter code", "en_US", false},
		{"valid 5-letter code", "de_DE", false},
		{"uppercase", "EN", false},
		{"with spaces", " en ", false},
		{"empty string", "", true},
		{"too short", "e", true},
		{"too long", "english", true},
		{"invalid length 3", "eng", true},
		{"invalid length 4", "engl", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLanguage(tt.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLanguage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePageNumber(t *testing.T) {
	tests := []struct {
		name    string
		page    int
		wantErr bool
	}{
		{"page 1", 1, false},
		{"valid page", 5, false},
		{"maximum valid", 50, false},
		{"zero", 0, true},
		{"negative", -1, true},
		{"exceeds maximum", 51, true},
		{"large number", 100, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePageNumber(tt.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePageNumber() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTimeRange(t *testing.T) {
	tests := []struct {
		name      string
		timeRange string
		wantErr   bool
	}{
		{"day", "day", false},
		{"week", "week", false},
		{"month", "month", false},
		{"year", "year", false},
		{"uppercase", "DAY", false},
		{"with spaces", " week ", false},
		{"empty - optional", "", false},
		{"invalid range", "hour", true},
		{"invalid range", "decade", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTimeRange(tt.timeRange)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTimeRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCategory(t *testing.T) {
	tests := []struct {
		name    string
		category string
		wantErr bool
	}{
		{"general", "general", false},
		{"images", "images", false},
		{"videos", "videos", false},
		{"news", "news", false},
		{"uppercase", "GENERAL", false},
		{"with spaces", " images ", false},
		{"custom category - allowed", "custom", false},
		{"empty string", "", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCategory(tt.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCategory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	err := ValidationError{
		Field:   "test_field",
		Value:   "test_value",
		Message: "test message",
	}
	
	expected := "validation error for field 'test_field': test message (value: test_value)"
	if err.Error() != expected {
		t.Errorf("ValidationError.Error() = %v, want %v", err.Error(), expected)
	}
}