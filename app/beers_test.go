package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBeersHandler_Get(t *testing.T) {
	defaultHandler := NewBeersHandler(getDefaultMockBeersRepository())

	t.Run("expect GET /beers to return 200", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/beers", nil)
		w := httptest.NewRecorder()
		router := prepareRouter(http.MethodGet, "/beers", defaultHandler.Get)
		router.ServeHTTP(w, r)

		resp := w.Result()

		assertStatusCode(t, resp, http.StatusOK)
	})

	t.Run("expect GET /beers to return 404 when limit param is invalid", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/beers?givenAt=2020-10-01&limit=bonk", nil)
		w := httptest.NewRecorder()
		router := prepareRouter(http.MethodGet, "/beers", defaultHandler.Get)
		router.ServeHTTP(w, r)

		resp := w.Result()

		assertStatusCode(t, resp, http.StatusBadRequest)
	})

	t.Run("expect GET /beers to return 200 with filter params", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/beers?givenAt=2020-10-01&limit=10&op=gt", nil)
		w := httptest.NewRecorder()
		router := prepareRouter(http.MethodGet, "/beers", defaultHandler.Get)
		router.ServeHTTP(w, r)

		resp := w.Result()

		assertStatusCode(t, resp, http.StatusOK)
	})
}
