package app

import (
	"appdoki-be/app/repositories"
	"net/http"
)

// BeersHandler holds handler dependencies
type BeersHandler struct {
	beersRepo repositories.BeersRepositoryInterface
}

// NewBeersHandler returns an initialized beers handler with the required dependencies
func NewBeersHandler(beersRepo repositories.BeersRepositoryInterface) *BeersHandler {
	return &BeersHandler{
		beersRepo: beersRepo,
	}
}

// Get gets all the beer transfers
func (h *BeersHandler) Get(w http.ResponseWriter, r *http.Request) {
	feed, err := h.beersRepo.GetBeerTransfers(r.Context())
	if err != nil {
		respondInternalError(w)
		return
	}

	respondJSON(w, feed, http.StatusOK)
}
