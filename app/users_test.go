package app

import (
	"appdoki-be/app/repositories"
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
	t.Run("expect GET /users to return 200 and a list of users", func(t *testing.T) {
		uh := NewUsersHandler(
			getDefaultMockUsersRepository(),
			getDefaultMockBeersRepository(),
			getMockNotifier())

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

		var users []*repositories.User
		err = json.Unmarshal(body, &users)
		if err != nil {
			t.Fatal("failed to parse response body")
		}
	})

	t.Run("expect GET /users to return 200 and an empty list of users ", func(t *testing.T) {
		mock := getDefaultMockUsersRepository()
		mock.getAllImpl = func(_ context.Context) ([]*repositories.User, error) {
			return []*repositories.User{}, nil
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

		var users []*repositories.User
		err = json.Unmarshal(body, &users)
		if err != nil {
			t.Fatal("failed to parse response body")
		}
		if len(users) > 0 {
			t.Fatal("response should have no records")
		}
	})
}
