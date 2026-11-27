// Package server serves the ui and exposes the REST endpoints
package server

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// spa serves the embedded UI. Real files are served as they are. Unknown
// paths without a file extension get index.html, so SvelteKit's client-side
// routing works when the page is refreshed on a route like /functions/abc.
func spa(fsys fs.FS) http.Handler {
	files := http.FileServerFS(fsys)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if name == "" {
			name = "index.html"
		}

		if _, err := fs.Stat(fsys, name); err != nil {
			if path.Ext(name) != "" { // missing asset (.js, .css, .png): real 404
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Cache-Control", "no-cache")
			http.ServeFileFS(w, r, fsys, "index.html")
			return
		}

		if strings.HasPrefix(name, "_app/immutable/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache")
		}
		files.ServeHTTP(w, r)
	})
}
