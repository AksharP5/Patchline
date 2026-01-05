package npm

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchPackageInfoRequiresName(t *testing.T) {
	_, err := fetchPackageInfo(context.Background(), &http.Client{}, "http://example.test", "")
	if err == nil {
		t.Fatalf("expected error for empty name")
	}
}

func TestFetchPackageInfoNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := fetchPackageInfo(context.Background(), server.Client(), server.URL, "missing")
	if !errors.Is(err, ErrPackageNotFound) {
		t.Fatalf("expected ErrPackageNotFound, got %v", err)
	}
}

func TestFetchPackageInfoNonSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	defer server.Close()

	_, err := fetchPackageInfo(context.Background(), server.Client(), server.URL, "bad")
	if err == nil {
		t.Fatalf("expected error for non-2xx response")
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected error to include response body, got %v", err)
	}
}

func TestFetchPackageInfoValidPayload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"pkg","dist-tags":{"latest":"1.2.3"},"versions":{"1.2.3":{},"1.0.0":{}}}`))
	}))
	defer server.Close()

	info, err := fetchPackageInfo(context.Background(), server.Client(), server.URL, "pkg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Name != "pkg" {
		t.Fatalf("expected name to be pkg, got %q", info.Name)
	}
	if info.Latest != "1.2.3" {
		t.Fatalf("expected latest to be 1.2.3, got %q", info.Latest)
	}
	if len(info.Versions) != 2 || info.Versions[0] != "1.0.0" || info.Versions[1] != "1.2.3" {
		t.Fatalf("unexpected versions list: %#v", info.Versions)
	}
}

func TestFetchPackageInfoInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{"))
	}))
	defer server.Close()

	_, err := fetchPackageInfo(context.Background(), server.Client(), server.URL, "broken")
	if err == nil {
		t.Fatalf("expected error for invalid json")
	}
}
