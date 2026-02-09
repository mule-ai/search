package errors

import (
	"errors"
	"net"
	"testing"
)

func TestSearchError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *SearchError
		wantMsg string
	}{
		{
			name: "basic error",
			err: &SearchError{
				Code:    ErrCodeEmptyQuery,
				Message: "Query is empty",
			},
			wantMsg: "[EMPTY_QUERY] Query is empty",
		},
		{
			name: "error with suggestion",
			err: &SearchError{
				Code:       ErrCodeInvalidFormat,
				Message:    "Invalid format",
				Suggestion: "Use json, markdown, or text",
			},
			wantMsg: "[INVALID_FORMAT] Invalid format\nSuggestion: Use json, markdown, or text",
		},
		{
			name: "error with wrapped error",
			err: &SearchError{
				Code:    ErrCodeConfigInvalid,
				Message: "Config is invalid",
				Err:     errors.New("yaml error"),
			},
			wantMsg: "[CONFIG_INVALID] Config is invalid\nDetails: yaml error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.wantMsg {
				t.Errorf("SearchError.Error() = %q, want %q", got, tt.wantMsg)
			}
		})
	}
}

func TestSearchError_Unwrap(t *testing.T) {
	wrapped := errors.New("wrapped error")
	err := &SearchError{
		Code:    ErrCodeConfigInvalid,
		Message: "Config error",
		Err:     wrapped,
	}

	if !errors.Is(err.Unwrap(), wrapped) {
		t.Error("Unwrap() should return the wrapped error")
	}
}

func TestNew(t *testing.T) {
	err := New(ErrCodeEmptyQuery, "Test message")
	if err.Code != ErrCodeEmptyQuery {
		t.Errorf("New() code = %v, want %v", err.Code, ErrCodeEmptyQuery)
	}
	if err.Message != "Test message" {
		t.Errorf("New() message = %v, want 'Test message'", err.Message)
	}
}

func TestWrap(t *testing.T) {
	wrapped := errors.New("wrapped")
	err := Wrap(ErrCodeConfigInvalid, "Config failed", wrapped)

	if err.Code != ErrCodeConfigInvalid {
		t.Errorf("Wrap() code = %v, want %v", err.Code, ErrCodeConfigInvalid)
	}
	if err.Message != "Config failed" {
		t.Errorf("Wrap() message = %v, want 'Config failed'", err.Message)
	}
	if err.Err != wrapped {
		t.Error("Wrap() should wrap the provided error")
	}
}

func TestWithSuggestion(t *testing.T) {
	err := New(ErrCodeEmptyQuery, "Test").
		WithSuggestion("Try adding a query")

	if err.Suggestion != "Try adding a query" {
		t.Errorf("WithSuggestion() = %v, want 'Try adding a query'", err.Suggestion)
	}
}

func TestWithVerbose(t *testing.T) {
	err := New(ErrCodeEmptyQuery, "Test").
		WithVerbose("Detailed debug info")

	if err.Verbose != "Detailed debug info" {
		t.Errorf("WithVerbose() = %v, want 'Detailed debug info'", err.Verbose)
	}
}

func TestConfigNotFound(t *testing.T) {
	path := "/test/config.yaml"
	err := ConfigNotFound(path)

	if err.Code != ErrCodeConfigNotFound {
		t.Errorf("ConfigNotFound() code = %v, want %v", err.Code, ErrCodeConfigNotFound)
	}
	if err.Suggestion == "" {
		t.Error("ConfigNotFound() should provide a suggestion")
	}
}

func TestConfigInvalid(t *testing.T) {
	wrapped := errors.New("parse error")
	err := ConfigInvalid(wrapped)

	if err.Code != ErrCodeConfigInvalid {
		t.Errorf("ConfigInvalid() code = %v, want %v", err.Code, ErrCodeConfigInvalid)
	}
	if err.Err != wrapped {
		t.Error("ConfigInvalid() should wrap the provided error")
	}
	if err.Suggestion == "" {
		t.Error("ConfigInvalid() should provide a suggestion")
	}
}

func TestNetworkError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		wantCode    ErrorCode
		wantSuggest bool
	}{
		{
			name:        "timeout error",
			err:         &timeoutError{},
			wantCode:    ErrCodeNetworkTimeout,
			wantSuggest: true,
		},
		{
			name:        "connection refused",
			err:         errors.New("connection refused"),
			wantCode:    ErrCodeConnectionRefused,
			wantSuggest: true,
		},
		{
			name:        "DNS error",
			err:         errors.New("no such host"),
			wantCode:    ErrCodeDNSFailed,
			wantSuggest: true,
		},
		{
			name:        "generic network error",
			err:         errors.New("some network error"),
			wantCode:    ErrCodeNetworkUnreachable,
			wantSuggest: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NetworkError(tt.err)
			if err.Code != tt.wantCode {
				t.Errorf("NetworkError() code = %v, want %v", err.Code, tt.wantCode)
			}
			hasSuggestion := err.Suggestion != ""
			if hasSuggestion != tt.wantSuggest {
				t.Errorf("NetworkError() has suggestion = %v, want %v", hasSuggestion, tt.wantSuggest)
			}
		})
	}
}

func TestHTTPError(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		status      string
		wantCode    ErrorCode
	}{
		{
			name:       "server error",
			statusCode: 500,
			status:     "Internal Server Error",
			wantCode:   ErrCodeAPIUnavailable,
		},
		{
			name:       "not found",
			statusCode: 404,
			status:     "Not Found",
			wantCode:   ErrCodeAPIError,
		},
		{
			name:       "unauthorized",
			statusCode: 401,
			status:     "Unauthorized",
			wantCode:   ErrCodeAPIError,
		},
		{
			name:       "bad request",
			statusCode: 400,
			status:     "Bad Request",
			wantCode:   ErrCodeAPIError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := HTTPStatusError(tt.statusCode, tt.status)
			if err.Code != tt.wantCode {
				t.Errorf("HTTPStatusError() code = %v, want %v", err.Code, tt.wantCode)
			}
		})
	}
}

func TestInvalidResponse(t *testing.T) {
	wrapped := errors.New("invalid JSON")
	err := InvalidResponse(wrapped)

	if err.Code != ErrCodeInvalidResponse {
		t.Errorf("InvalidResponse() code = %v, want %v", err.Code, ErrCodeInvalidResponse)
	}
	if err.Suggestion == "" {
		t.Error("InvalidResponse() should provide a suggestion")
	}
}

func TestEmptyQuery(t *testing.T) {
	err := EmptyQuery()

	if err.Code != ErrCodeEmptyQuery {
		t.Errorf("EmptyQuery() code = %v, want %v", err.Code, ErrCodeEmptyQuery)
	}
	if err.Suggestion == "" {
		t.Error("EmptyQuery() should provide a suggestion")
	}
}

func TestInvalidFormat(t *testing.T) {
	format := "xml"
	err := InvalidFormat(format)

	if err.Code != ErrCodeInvalidFormat {
		t.Errorf("InvalidFormat() code = %v, want %v", err.Code, ErrCodeInvalidFormat)
	}
	if err.Suggestion == "" {
		t.Error("InvalidFormat() should provide a suggestion")
	}
}

func TestInvalidURL(t *testing.T) {
	url := "not-a-url"
	err := InvalidURL(url)

	if err.Code != ErrCodeInvalidURL {
		t.Errorf("InvalidURL() code = %v, want %v", err.Code, ErrCodeInvalidURL)
	}
	if err.Suggestion == "" {
		t.Error("InvalidURL() should provide a suggestion")
	}
}

func TestInvalidRange(t *testing.T) {
	param := "results"
	min, max, value := 1, 100, 150
	err := InvalidRange(param, min, max, value)

	if err.Code != ErrCodeInvalidRange {
		t.Errorf("InvalidRange() code = %v, want %v", err.Code, ErrCodeInvalidRange)
	}
	if err.Suggestion == "" {
		t.Error("InvalidRange() should provide a suggestion")
	}
}

func TestNoResults(t *testing.T) {
	query := "test query"
	err := NoResults(query)

	if err.Code != ErrCodeNoResults {
		t.Errorf("NoResults() code = %v, want %v", err.Code, ErrCodeNoResults)
	}
	if err.Suggestion == "" {
		t.Error("NoResults() should provide a suggestion")
	}
}

func TestIsSearchError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantBool bool
	}{
		{
			name:     "SearchError",
			err:      EmptyQuery(),
			wantBool: true,
		},
		{
			name:     "standard error",
			err:      errors.New("standard error"),
			wantBool: false,
		},
		{
			name:     "nil error",
			err:      nil,
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := IsSearchError(tt.err)
			if ok != tt.wantBool {
				t.Errorf("IsSearchError() = %v, want %v", ok, tt.wantBool)
			}
		})
	}
}

func TestGetErrorCode(t *testing.T) {
	tests := []struct {
		name  string
		err   error
		want  ErrorCode
	}{
		{
			name: "SearchError with code",
			err:  EmptyQuery(),
			want: ErrCodeEmptyQuery,
		},
		{
			name: "standard error",
			err:  errors.New("standard error"),
			want: "",
		},
		{
			name: "nil error",
			err:  nil,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetErrorCode(tt.err)
			if got != tt.want {
				t.Errorf("GetErrorCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPErrorType(t *testing.T) {
	err := &HTTPResponseError{
		StatusCode: 404,
		Status:     "Not Found",
		Body:       "404 page not found",
	}

	got := err.Error()
	want := "HTTP 404: Not Found"
	if got != want {
		t.Errorf("HTTPResponseError.Error() = %v, want %v", got, want)
	}
}

// Helper type for timeout testing
type timeoutError struct{}

func (e *timeoutError) Error() string   { return "timeout" }
func (e *timeoutError) Timeout() bool   { return true }
func (e *timeoutError) Temporary() bool { return true }

// Test that we can detect net.Error
func TestNetErrorDetection(t *testing.T) {
	err := &timeoutError{}
	var netErr net.Error = err

	if !netErr.Timeout() {
		t.Error("Expected timeout error to have Timeout() = true")
	}
}

func TestWrapNetworkError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{
			name:    "nil error",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "standard error",
			err:     errors.New("network error"),
			wantErr: true,
		},
		{
			name:    "timeout error",
			err:     &timeoutError{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WrapNetworkError(tt.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("WrapNetworkError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Benchmark tests
func BenchmarkSearchError_Error(b *testing.B) {
	err := &SearchError{
		Code:       ErrCodeEmptyQuery,
		Message:    "Test message",
		Suggestion: "Test suggestion",
	}
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkNetworkError(b *testing.B) {
	wrapped := &timeoutError{}
	for i := 0; i < b.N; i++ {
		_ = NetworkError(wrapped)
	}
}
