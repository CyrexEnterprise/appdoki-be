package app

import (
	"appdoki-be/app/repositories"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUsersRepository struct {
	getAllImpl           func(ctx context.Context) ([]*repositories.User, error)
	findByIDImpl         func(ctx context.Context, ID string) (*repositories.User, error)
	findByEmailImpl      func(ctx context.Context, email string) (*repositories.User, error)
	findOrCreateUserImpl func(ctx context.Context, userData *repositories.User) (*repositories.User, error)
	createImpl           func(ctx context.Context, user *repositories.User) (*repositories.User, error)
	updateImpl           func(ctx context.Context, user *repositories.User) (*repositories.User, error)
	deleteImpl           func(ctx context.Context, ID string) (bool, error)
	addBeerTransferImpl  func(ctx context.Context, giverID string, takerID string, beers int) error
}

func (r *mockUsersRepository) GetAll(ctx context.Context) ([]*repositories.User, error) {
	return r.getAllImpl(ctx)
}

func (r *mockUsersRepository) FindByID(ctx context.Context, ID string) (*repositories.User, error) {
	return r.findByIDImpl(ctx, ID)
}

func (r *mockUsersRepository) FindByEmail(ctx context.Context, email string) (*repositories.User, error) {
	return r.findByEmailImpl(ctx, email)
}

func (r *mockUsersRepository) Create(ctx context.Context, user *repositories.User) (*repositories.User, error) {
	return r.createImpl(ctx, user)
}

func (r *mockUsersRepository) FindOrCreateUser(ctx context.Context, userData *repositories.User) (*repositories.User, error) {
	return r.findOrCreateUserImpl(ctx, userData)
}

func (r *mockUsersRepository) Update(ctx context.Context, user *repositories.User) (*repositories.User, error) {
	return r.updateImpl(ctx, user)
}

func (r *mockUsersRepository) Delete(ctx context.Context, ID string) (bool, error) {
	return r.deleteImpl(ctx, ID)
}

func (r *mockUsersRepository) AddBeerTransfer(ctx context.Context, giverID string, takerID string, beers int) error {
	return nil
}

func getDefaultMockUsersRepository() *mockUsersRepository {
	mockUser := &repositories.User{
		ID:    "1",
		Name:  "Roger Rabbit",
		Email: "rrabbit@acme.com",
	}

	return &mockUsersRepository{
		getAllImpl: func(_ context.Context) ([]*repositories.User, error) {
			mockUsers := []*repositories.User{mockUser}
			return mockUsers, nil
		},
		findByIDImpl: func(ctx context.Context, ID string) (*repositories.User, error) {
			return mockUser, nil
		},
		findByEmailImpl: func(ctx context.Context, email string) (*repositories.User, error) {
			return mockUser, nil
		},
		findOrCreateUserImpl: func(ctx context.Context, user *repositories.User) (*repositories.User, error) {
			return mockUser, nil
		},
		createImpl: func(ctx context.Context, user *repositories.User) (*repositories.User, error) {
			return mockUser, nil
		},
		updateImpl: func(ctx context.Context, user *repositories.User) (*repositories.User, error) {
			return mockUser, nil
		},
		deleteImpl: func(ctx context.Context, ID string) (bool, error) {
			return true, nil
		},
	}
}

func TestUsersHandler_Get(t *testing.T) {
	t.Run("expect GET /users to return 200 and a list of users", func(t *testing.T) {
		uh := NewUsersHandler(getDefaultMockUsersRepository())

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
		uh := NewUsersHandler(mock)

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
