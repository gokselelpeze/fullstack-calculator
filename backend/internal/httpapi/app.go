package httpapi

import (
	"net/http"
	"os"
	"path"
	"strings"
)

// NewAppHandler serves the API and, when configured, the built frontend.
// There are no client-side routes, so missing files stay 404s.
func NewAppHandler(staticDir string) http.Handler {
	api := NewHandler()
	if staticDir == "" {
		return api
	}
	files := http.FileServer(http.Dir(staticDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api" || strings.HasPrefix(r.URL.Path, "/api/") {
			api.ServeHTTP(w, r)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			fail(w, http.StatusMethodNotAllowed, codeMethodNotAllowed, "Use GET or HEAD for frontend files.")
			return
		}
		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if name == "" {
			name = "index.html"
		}
		// Avoid directory listings. FileServer handles the actual file response.
		info, err := os.Stat(staticDir + "/" + name)
		if err != nil || info.IsDir() {
			http.NotFound(w, r)
			return
		}
		files.ServeHTTP(w, r)
	})
}
