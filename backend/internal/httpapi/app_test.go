package httpapi

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppHandler(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("<h1>Calculator</h1>"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "assets"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "assets", "app.js"), []byte("console.log('calculator')"), 0600); err != nil {
		t.Fatal(err)
	}
	handler := NewAppHandler(dir)
	for _, tt := range []struct {
		method, path, body string
		status             int
		content            string
	}{
		{"GET", "/", "", 200, "Calculator"},
		{"HEAD", "/", "", 200, ""},
		{"GET", "/assets/app.js", "", 200, "calculator"},
		{"GET", "/assets/", "", 404, ""},
		{"GET", "/missing.js", "", 404, ""},
		{"POST", "/", "", 405, "METHOD_NOT_ALLOWED"},
		{"POST", "/api/calculate", `{"expression":"sqrt(9)+2^3"}`, 200, `"result":"11"`},
		{"GET", "/api/missing", "", 404, "NOT_FOUND"},
		{"GET", "/api", "", 404, "NOT_FOUND"},
	} {
		t.Run(tt.method+tt.path, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != tt.status || !strings.Contains(response.Body.String(), tt.content) {
				t.Fatalf("got %d: %s", response.Code, response.Body)
			}
			if tt.method == "HEAD" && response.Body.Len() != 0 {
				t.Fatal("HEAD returned a body")
			}
		})
	}
	response := httptest.NewRecorder()
	NewAppHandler("").ServeHTTP(response, httptest.NewRequest("GET", "/", nil))
	if response.Code != 404 {
		t.Fatal("API-only mode should not serve frontend")
	}
}
