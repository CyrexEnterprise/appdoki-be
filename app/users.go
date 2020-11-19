package app

import (
	"appdoki-be/app/repositories"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type CreateUserPayload struct {
	Name  string
	Email string
}

func (p *CreateUserPayload) validate() []string {
	var errs []string

	if len(p.Name) < 3 {
		errs = append(errs, "name: invalid length")
	}

	if len(p.Email) < 5 {
		errs = append(errs, "email: invalid length")
	}

	return errs
}

// UsersHandler holds handler dependencies
type UsersHandler struct {
	userRepo repositories.UsersRepositoryInterface
}

// NewUsersHandler returns an initialized users handler with the required dependencies
func NewUsersHandler(userRepo repositories.UsersRepositoryInterface) *UsersHandler {
	return &UsersHandler{
		userRepo: userRepo,
	}
}

// Get gets all users
func (h *UsersHandler) Get(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepo.GetAll(r.Context())
	if err != nil {
		respondInternalError(w)
		return
	}

	respondJSON(w, users, http.StatusOK)
}

// GetByID tries to get a user by ID
func (h *UsersHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, ok := vars["id"]
	if !ok {
		respondJSON(w, &appError{
			Errors: []string{"invalid id param"},
		}, http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.FindByID(r.Context(), uid)
	if err != nil {
		respondInternalError(w)
		return
	}

	if user == nil {
		respondJSON(w, &appError{
			Errors: []string{"user not found"},
		}, http.StatusNotFound)
		return
	}

	respondJSON(w, user, http.StatusOK)
}

// Create creates a new user
func (h *UsersHandler) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var userPayload CreateUserPayload
	err := decoder.Decode(&userPayload)
	if err != nil {
		respondInternalError(w)
		return
	}

	errs := userPayload.validate()
	if errs != nil {
		respondJSON(w, &appError{
			Errors: errs,
		}, http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.Create(r.Context(), &repositories.User{
		Name:  userPayload.Name,
		Email: userPayload.Email,
	})
	if err != nil {
		if e, ok := err.(*repositories.ConflictError); ok {
			respondJSON(w, &appError{
				Errors: []string{e.Message},
			}, http.StatusConflict)
			return
		}

		respondInternalError(w)
		return
	}

	respondJSON(w, user, http.StatusCreated)
}

// Update updates a user
func (h *UsersHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, ok := vars["id"]
	if !ok {
		respondJSON(w, &appError{
			Errors: []string{"invalid id param"},
		}, http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var userPayload CreateUserPayload
	err := decoder.Decode(&userPayload)
	if err != nil {
		respondInternalError(w)
		return
	}

	errs := userPayload.validate()
	if errs != nil {
		respondJSON(w, &appError{
			Errors: errs,
		}, http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.Update(r.Context(), &repositories.User{
		ID:    uid,
		Name:  userPayload.Name,
		Email: userPayload.Email,
	})
	if err != nil {
		if e, ok := err.(*repositories.ConflictError); ok {
			respondJSON(w, &appError{
				Errors: []string{e.Message},
			}, http.StatusConflict)
			return
		}

		respondInternalError(w)
		return
	}

	if user == nil {
		respondJSON(w, &appError{
			Errors: []string{"user not found"},
		}, http.StatusNotFound)
		return
	}

	respondJSON(w, user, http.StatusOK)
}
