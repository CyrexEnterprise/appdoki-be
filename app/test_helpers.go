package app

import (
	"github.com/gorilla/mux"
	"net/http"
	"testing"
)

func prepareRouter(method string, path string, h func(http.ResponseWriter, *http.Request)) *mux.Router {
	router := mux.NewRouter()
	router.Methods(method).
		Path(path).
		HandlerFunc(h)
	return router
}

func assertJSONContentType(t *testing.T, r *http.Response) {
	if r.Header.Get("Content-Type") != "application/json; charset=utf-8" {
		t.Fatalf("expected 'application/json', got '%s'", r.Header.Get("Content-Type"))
	}
}

func assertStatusCode(t *testing.T, r *http.Response, status int) {
	if r.StatusCode != status {
		t.Fatalf("expected %d response, got %d", status, r.StatusCode)
	}
}
