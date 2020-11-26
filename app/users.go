package app

import (
	"appdoki-be/app/repositories"
	"encoding/json"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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
	notifier notifier
}

// NewUsersHandler returns an initialized users handler with the required dependencies
func NewUsersHandler(userRepo repositories.UsersRepositoryInterface, notifierSrv notifier) *UsersHandler {
	return &UsersHandler{
		userRepo: userRepo,
		notifier: notifierSrv,
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

func (h *UsersHandler) GiveBeers(w http.ResponseWriter, r *http.Request) {
	userID := fmt.Sprintf("%v", r.Context().Value("userID"))

	vars := mux.Vars(r)
	takerUserId, ok := vars["id"]
	if !ok {
		respondJSON(w, &appError{
			Errors: []string{"invalid id param"},
		}, http.StatusBadRequest)
		return
	}

	if userID == takerUserId {
		respondJSON(w, &appError{
			Errors: []string{"oi, cheeky bastard, give beers to others"},
		}, http.StatusForbidden)
		return
	}

	beersParam, ok := vars["beers"]
	if !ok {
		respondJSON(w, &appError{
			Errors: []string{"invalid beers param"},
		}, http.StatusBadRequest)
		return
	}

	beers, err := strconv.Atoi(beersParam)
	if err != nil {
		respondInternalError(w)
		return
	}

	if beers <= 0 {
		respondJSON(w, &appError{
			Errors: []string{"invalid amount of beers: don't be a cheap bastard!"},
		}, http.StatusBadRequest)
		return
	}

	err = h.userRepo.AddBeerTransfer(r.Context(), userID, takerUserId, beers)
	if err != nil {
		respondInternalError(w)
		return
	}

	go func() {
		giver, _ := h.userRepo.FindByID(r.Context(), userID)
		receiver, _ := h.userRepo.FindByID(r.Context(), takerUserId)

		if giver == nil || receiver == nil {
			log.Warn("could not fetch users for beers notification")
			return
		}

		h.notifier.notifyAll(beersTopic, messaging.Notification{
			Title: "BeerTab event",
			Body:  fmt.Sprintf("%s just rewarded %s with %n beers!", giver.Name, receiver.Name, beers),
		})
	}()

	respondNoContent(w, http.StatusNoContent)
}

func (h *UsersHandler) BeersSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["id"]
	if !ok {
		respondJSON(w, &appError{
			Errors: []string{"invalid id param"},
		}, http.StatusBadRequest)
		return
	}

	beerLog, err := h.userRepo.GetBeerTransfersSummary(r.Context(), userID)
	if err != nil {
		respondInternalError(w)
	}

	respondJSON(w, beerLog, http.StatusOK)
}
