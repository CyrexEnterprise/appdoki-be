package app

import (
	"appdoki-be/app/repositories"
	"net/http"
	"strconv"
	"time"
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

	limitParam := r.URL.Query().Get("limit")
	if len(limitParam) > 0 {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			respondJSON(w, &appError{
				Errors: []string{"invalid limit param"},
			}, http.StatusBadRequest)
			return
		}

		options.Limit = limit
	}

	options.GivenAt = r.URL.Query().Get("givenAt")
	if len(options.GivenAt) == 0 {
		options.GivenAt = time.Now().Format(time.RFC3339)
	}

	switch r.URL.Query().Get("op") {
	case "gt":
		options.SetGtOperator()
	case "lt":
		options.SetLtOperator()
	default:
		options.SetLtOperator()
	}

	feed, err := h.beersRepo.GetBeerTransfers(r.Context(), options)
	if err != nil {
		respondInternalError(w)
		return
	}

	respondJSON(w, feed, http.StatusOK)
}
