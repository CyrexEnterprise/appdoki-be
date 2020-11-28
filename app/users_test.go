package app

import (
	repos "appdoki-be/app/repositories"
	"context"
	"encoding/json"
	"firebase.google.com/go/v4/messaging"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockNotifier struct{}

func (n *mockNotifier) notifyAll(_ string, _ *messaging.Notification, _ map[string]string) {}
func (n *mockNotifier) messageAll(_ string, _ map[string]string)                           {}

func getMockNotifier() *mockNotifier {
	return &mockNotifier{}
}

func TestUsersHandler_Get(t *testing.T) {
	defaultHandler := NewUsersHandler(
		getDefaultMockUsersRepository(),
		getDefaultMockBeersRepository(),
		getMockNotifier())

	t.Run("expect GET /users to return 200 and a list of users", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()
		router := prepareRouter(http.MethodGet, "/users", defaultHandler.Get)
		router.ServeHTTP(w, r)

		resp := w.Result()

		assertStatusCode(t, resp, http.StatusOK)
		assertJSONContentType(t, resp)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal("failed to read response body")
		}

		var users []*repos.User
		err = json.Unmarshal(body, &users)
		if err != nil {
			t.Fatal("failed to parse response body")
		}
	})

	t.Run("expect GET /users to return 200 and an empty list of users ", func(t *testing.T) {
		mock := getDefaultMockUsersRepository()
		mock.getAllImpl = func(_ context.Context) ([]*repos.User, error) {
			return []*repos.User{}, nil
		}
		uh := NewUsersHandler(mock, getDefaultMockBeersRepository(), getMockNotifier())

		r := httptest.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()
		router := prepareRouter(http.MethodGet, "/users", uh.Get)
		router.ServeHTTP(w, r)

		resp := w.Result()

		assertStatusCode(t, resp, http.StatusOK)
		assertJSONContentType(t, resp)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal("failed to read response body")
		}

		var users []*repos.User
		err = json.Unmarshal(body, &users)
		if err != nil {
			t.Fatal("failed to parse response body")
		}
		if len(users) > 0 {
			t.Fatal("response should have no records")
		}
	})
}

func TestUsersHandler_GetByID(t *testing.T) {
	t.Run("expect GET /users/{id} to return 200 and a user", func(t *testing.T) {
		mock := getDefaultMockUsersRepository()
		mock.findByIDImpl = func(ctx context.Context, ID string) (*repos.User, error) {
			return generateRandomUserMockWithID("1"), nil
		}
		uh := NewUsersHandler(mock, getDefaultMockBeersRepository(), getMockNotifier())

		r := httptest.NewRequest("GET", "/users/1", nil)
		w := httptest.NewRecorder()
		router := prepareRouter(http.MethodGet, "/users/{id}", uh.GetByID)
		router.ServeHTTP(w, r)

		resp := w.Result()

		assertStatusCode(t, resp, http.StatusOK)
		assertJSONContentType(t, resp)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal("failed to read response body")
		}

		var users repos.User
		err = json.Unmarshal(body, &users)
		if err != nil {
			t.Fatal("failed to parse response body")
		}
	})

	t.Run("expect GET /users/{id} to return 404 when param is invalid", func(t *testing.T) {
		mock := getDefaultMockUsersRepository()
		mock.findByIDImpl = func(ctx context.Context, ID string) (*repos.User, error) {
			return nil, nil
		}
		uh := NewUsersHandler(mock, getDefaultMockBeersRepository(), getMockNotifier())

		r := httptest.NewRequest("GET", "/users/1", nil)
		w := httptest.NewRecorder()
		router := prepareRouter(http.MethodGet, "/users/{id}", uh.GetByID)
		router.ServeHTTP(w, r)

		resp := w.Result()

		assertStatusCode(t, resp, http.StatusNotFound)
	})
}

func TestUsersHandler_GiveBeers(t *testing.T) {
	defaultHandler := NewUsersHandler(
		getDefaultMockUsersRepository(),
		getDefaultMockBeersRepository(),
		getMockNotifier())

	t.Run("expect POST /users/{id}/beers/{beers} to return 403 when a user gives beers to self", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/users/1/beers/10", nil)
		r = r.WithContext(context.WithValue(r.Context(), "userID", "1"))
		w := httptest.NewRecorder()
		router := prepareRouter(http.MethodPost, "/users/{id}/beers/{beers}", defaultHandler.GiveBeers)
		router.ServeHTTP(w, r)
		resp := w.Result()

		assertStatusCode(t, resp, http.StatusForbidden)
	})

	t.Run("expect POST /users/{id}/beers/{beers} to return 404 when beers param invalid", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/users/999/beers/bonk", nil)
		r = r.WithContext(context.WithValue(r.Context(), "userID", "1"))
		w := httptest.NewRecorder()
		router := prepareRouter(http.MethodPost, "/users/{id}/beers/{beers}", defaultHandler.GiveBeers)
		router.ServeHTTP(w, r)
		resp := w.Result()

		assertStatusCode(t, resp, http.StatusBadRequest)
	})

	t.Run("expect POST /users/{id}/beers/{beers} to return 404 when beers param negative", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/users/999/beers/-22", nil)
		r = r.WithContext(context.WithValue(r.Context(), "userID", "1"))
		w := httptest.NewRecorder()
		router := prepareRouter(http.MethodPost, "/users/{id}/beers/{beers}", defaultHandler.GiveBeers)
		router.ServeHTTP(w, r)
		resp := w.Result()

		assertStatusCode(t, resp, http.StatusBadRequest)
	})

	t.Run("expect POST /users/{id}/beers/{beers} to return 200", func(t *testing.T) {
		urMock := getDefaultMockUsersRepository()
		urMock.addBeerTransferImpl = func(ctx context.Context, giverID string, takerID string, beers int) (int, error) {
			return 5, nil
		}
		brMock := getDefaultMockBeersRepository()
		brMock.getBeerTransferImpl = func(ctx context.Context, id int) (*repos.BeerTransferFeedItem, error) {
			return generateRandomBeerTransferMock(), nil
		}
		uh := NewUsersHandler(urMock, brMock, getMockNotifier())

		r := httptest.NewRequest("POST", "/users/999/beers/10", nil)
		r = r.WithContext(context.WithValue(r.Context(), "userID", "1"))
		w := httptest.NewRecorder()
		router := prepareRouter(http.MethodPost, "/users/{id}/beers/{beers}", uh.GiveBeers)
		router.ServeHTTP(w, r)
		resp := w.Result()

		assertStatusCode(t, resp, http.StatusNoContent)
	})
}

func TestUsersHandler_BeersSummary(t *testing.T) {
	defaultHandler := NewUsersHandler(
		getDefaultMockUsersRepository(),
		getDefaultMockBeersRepository(),
		getMockNotifier())

	t.Run("expect GET /users/{id}/beers to return 200", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/users/1/beers", nil)
		w := httptest.NewRecorder()
		router := prepareRouter(http.MethodGet, "/users/{id}/beers", defaultHandler.BeersSummary)
		router.ServeHTTP(w, r)

		resp := w.Result()

		assertStatusCode(t, resp, http.StatusOK)
		assertJSONContentType(t, resp)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal("failed to read response body")
		}

		var log *repos.UserBeerLog
		err = json.Unmarshal(body, &log)
		if err != nil {
			t.Fatal("failed to parse response body")
		}
	})
}
