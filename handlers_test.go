package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
)

func TestIndex(t *testing.T) {
	mr := miniredis.RunT(t)

	// MIRRORS hash: mirror_id -> mirror_name
	mr.HSet("MIRRORS", "berlin", "Berlin")

	// Mirror details
	mr.HSet("MIRROR_berlin",
		"http", "https://mirror.example.com/",
		"enabled", "1",
		"continentCode", "EU",
		"countryCodes", "de",
	)

	ctx := &appContext{redisPool: NewPool(mr.Addr())}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := appHandler{ctx, Index}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `"Status":{"status":"ok"},"MirrorList":[{"ID":"Berlin"`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
