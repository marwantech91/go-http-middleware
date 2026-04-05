package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimit_AllowsRequestsUnderLimit(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := RateLimit(5, 1)
	wrapped := mw(handler)

	req := httptest.NewRequest("GET", "/api/data", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestRateLimit_BlocksAfterLimit(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := RateLimit(3, 0)
	wrapped := mw(handler)

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/api/data", nil)
		req.RemoteAddr = "10.0.0.1:12345"
		rec := httptest.NewRecorder()
		wrapped.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i, rec.Code)
		}
	}

	// 4th request should be blocked
	req := httptest.NewRequest("GET", "/api/data", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", rec.Code)
	}
}

func TestRateLimit_TracksSeparateIPs(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := RateLimit(1, 0)
	wrapped := mw(handler)

	// Exhaust limit for IP A
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	// IP B should still work
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.RemoteAddr = "10.0.0.2:1234"
	rec2 := httptest.NewRecorder()
	wrapped.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Errorf("expected 200 for different IP, got %d", rec2.Code)
	}
}

func TestRecovery_CatchesPanics(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	})

	mw := Recovery()
	wrapped := mw(handler)

	req := httptest.NewRequest("GET", "/panic", nil)
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
}
