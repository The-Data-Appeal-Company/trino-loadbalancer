package ui

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type UI struct {
	path      string
	indexPath string
}

func New(path string) UI {
	return UI{
		path:      path,
		indexPath: "index.html",
	}
}

func (a UI) Router() *mux.Router {
	r := mux.NewRouter()

	r.PathPrefix("/ui").HandlerFunc(a.ServeUI)

	return r
}

// ServeUI inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h UI) ServeUI(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	path = strings.Replace(path, "/ui", "", 1)

	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.path, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		indexPage := filepath.Join(h.path, h.indexPath)
		http.ServeFile(w, r, indexPage)
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, path)
}

func (a UI) Serve(addr string) error {
	return http.ListenAndServe(addr, a.Router())
}
