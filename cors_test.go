package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS_SetsHeadersForAllowedOrigin(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := CORS(CORSOptions{AllowOrigins: []string{"https://example.com"}})
	wrapped := mw(handler)

	req := httptest.NewRequest("GET", "/api/data", nil)
	req.Header.Set("Origin", "https://example.com")
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://example.com" {
		t.Errorf("expected origin https://example.com, got %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Error("expected Allow-Methods header to be set")
	}
}

func TestCORS_WildcardAllowsAnyOrigin(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := CORS(CORSOptions{AllowOrigins: []string{"*"}})
	wrapped := mw(handler)

	req := httptest.NewRequest("GET", "/api/data", nil)
	req.Header.Set("Origin", "https://anything.com")
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://anything.com" {
		t.Errorf("expected origin https://anything.com, got %q", got)
	}
}

func TestCORS_RejectsDisallowedOrigin(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := CORS(CORSOptions{AllowOrigins: []string{"https://allowed.com"}})
	wrapped := mw(handler)

	req := httptest.NewRequest("GET", "/api/data", nil)
	req.Header.Set("Origin", "https://evil.com")
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("expected no CORS header, got %q", got)
	}
}

func TestCORS_PreflightReturnsNoContent(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called for preflight")
	})

	mw := CORS(CORSOptions{AllowOrigins: []string{"*"}})
	wrapped := mw(handler)

	req := httptest.NewRequest("OPTIONS", "/api/data", nil)
	req.Header.Set("Origin", "https://example.com")
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
}

func TestCORS_SetsMaxAgeHeader(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := CORS(CORSOptions{AllowOrigins: []string{"*"}, MaxAge: 3600})
	wrapped := mw(handler)

	req := httptest.NewRequest("GET", "/api/data", nil)
	req.Header.Set("Origin", "https://example.com")
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Max-Age"); got != "3600" {
		t.Errorf("expected max-age 3600, got %q", got)
	}
}

func TestFormatInt(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{42, "42"},
		{3600, "3600"},
	}

	for _, tc := range tests {
		got := formatInt(tc.input)
		if got != tc.expected {
			t.Errorf("formatInt(%d) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}
