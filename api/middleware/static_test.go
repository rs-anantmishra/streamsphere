package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// the entrypoint rewrites main-*.js with the environment's addresses at every start, under the same file name, and
// index.html names that file; a browser must ask again for both instead of keeping them for an hour
func TestStaticFilesThatChangeWithTheEnvironmentAreRevalidated(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "assets"), 0o755)
	for name, body := range map[string]string{"index.html": "<html></html>", "main-ABC123.js": "var a", "chunk-XYZ.js": "var b", "assets/x.gif": "GIF"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	app := fiber.New()
	serveStatic(app, http.Dir(dir))

	cases := []struct{ path, want string }{
		{"/", "no-cache"},
		{"/index.html", "no-cache"},
		{"/main-ABC123.js", "no-cache"},
		{"/chunk-XYZ.js", "public, max-age=3600"},
		{"/assets/x.gif", "public, max-age=3600"},
	}
	for _, c := range cases {
		resp, err := app.Test(httptest.NewRequest("GET", c.path, nil))
		if err != nil {
			t.Fatalf("GET %s: %v", c.path, err)
		}
		if resp.StatusCode != 200 {
			t.Errorf("GET %s: status %d, want 200", c.path, resp.StatusCode)
		}
		if got := resp.Header.Get("Cache-Control"); got != c.want {
			t.Errorf("GET %s: Cache-Control = %q, want %q", c.path, got, c.want)
		}
	}
}
