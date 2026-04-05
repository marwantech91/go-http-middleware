package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChain_AppliesMiddlewareInOrder(t *testing.T) {
	var order []string

	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw1-before")
			next.ServeHTTP(w, r)
			order = append(order, "mw1-after")
		})
	}

	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw2-before")
			next.ServeHTTP(w, r)
			order = append(order, "mw2-after")
		})
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
		w.WriteHeader(http.StatusOK)
	})

	chained := Chain(handler, mw1, mw2)

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	chained.ServeHTTP(rec, req)

	expected := []string{"mw1-before", "mw2-before", "handler", "mw2-after", "mw1-after"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d calls, got %d", len(expected), len(order))
	}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("position %d: expected %q, got %q", i, v, order[i])
		}
	}
}

func TestChain_EmptyMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	chained := Chain(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	chained.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
