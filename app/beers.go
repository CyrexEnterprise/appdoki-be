package app

import (
	"appdoki-be/app/repositories"
	"net/http"
	"strconv"
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
	options := &repositories.BeerFeedPaginationOptions{
		Limit:   20,
		GivenAt: "",
	}

	givenAt := r.URL.Query().Get("givenAt")

	if len(givenAt) > 0 {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			respondJSON(w, &appError{
				Errors: []string{"invalid limit param"},
			}, http.StatusBadRequest)
			return
		}

		options.Limit = limit
		options.GivenAt = givenAt

		switch r.URL.Query().Get("op") {
		case "gt": options.SetGtOperator()
		case "lt": options.SetLtOperator()
		default: options.SetGtOperator()
		}
	}
	
	feed, err := h.beersRepo.GetBeerTransfers(r.Context(), options)
	if err != nil {
		respondInternalError(w)
		return
	}

	respondJSON(w, feed, http.StatusOK)
}
