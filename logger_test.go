package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogger_RecordsStatusCode(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	mw := Logger()
	wrapped := mw(handler)

	req := httptest.NewRequest("POST", "/api/users", nil)
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}
}

func TestLogger_DefaultsTo200(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	mw := Logger()
	wrapped := mw(handler)

	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
