package app

import (
	"appdoki-be/app/repositories"
	"context"
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
	userRepo  repositories.UsersRepositoryInterface
	beersRepo repositories.BeersRepositoryInterface
	notifier  notifier
}

// NewUsersHandler returns an initialized users handler with the required dependencies
func NewUsersHandler(
	userRepo repositories.UsersRepositoryInterface,
	beersRepo repositories.BeersRepositoryInterface,
	notifierSrv notifier) *UsersHandler {
	return &UsersHandler{
		userRepo:  userRepo,
		beersRepo: beersRepo,
		notifier:  notifierSrv,
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
		log.Error("could not read id param in UsersHandler.GetByID")
		respondInternalError(w)
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

// GiveBeers creates a beer transaction between two users
func (h *UsersHandler) GiveBeers(w http.ResponseWriter, r *http.Request) {
	userID := fmt.Sprintf("%v", r.Context().Value("userID"))

	vars := mux.Vars(r)
	takerUserId, ok := vars["id"]
	if !ok {
		log.Error("could not read id param in UsersHandler.GiveBeers")
		respondInternalError(w)
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
		log.Error("could not read beers param in UsersHandler.GiveBeers")
		respondInternalError(w)
		return
	}

	beers, err := strconv.Atoi(beersParam)
	if err != nil {
		respondJSON(w, &appError{
			Errors: []string{"invalid beers param: number expected"},
		}, http.StatusBadRequest)
		return
	}

	if beers <= 0 {
		respondJSON(w, &appError{
			Errors: []string{"invalid amount of beers: don't be a cheap bastard!"},
		}, http.StatusBadRequest)
		return
	}

	transferID, err := h.userRepo.AddBeerTransfer(r.Context(), userID, takerUserId, beers)
	if err != nil {
		respondInternalError(w)
		return
	}

	go func() {
		backgroundCtx := context.Background()
		transfer, err := h.beersRepo.GetBeerTransfer(backgroundCtx, transferID)
		if err != nil {
			log.Error("failed to get beer transfer")
			return
		}

		giverJSON, _ := json.Marshal(transfer.Giver)
		receiverJSON, _ := json.Marshal(transfer.Receiver)
		notification := &messaging.Notification{
			Title: "BeerTab event",
			Body:  fmt.Sprintf("%s just rewarded %s with %d beers!", transfer.Giver.Name, transfer.Receiver.Name, beers),
		}
		data := map[string]string{
			"id":    	string(transferID),
			"giver":    string(giverJSON),
			"receiver": string(receiverJSON),
			"beers":    strconv.Itoa(beers),
			"givenAt":  transfer.GivenAt,
		}

		h.notifier.notifyAll(beersTopic, notification, data)
	}()

	respondNoContent(w, http.StatusNoContent)
}

// BeersSummary generates a short beer transfer summary for a user
func (h *UsersHandler) BeersSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["id"]
	if !ok {
		log.Error("could not read id param in UsersHandler.BeersSummary")
		respondInternalError(w)
		return
	}

	beerLog, err := h.userRepo.GetBeerTransfersSummary(r.Context(), userID)
	if err != nil {
		respondInternalError(w)
	}

	respondJSON(w, beerLog, http.StatusOK)
}
