package repositories

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type BeerTransferFeedItem struct {
	Beers        int    `json:"beers"`
	GivenAt      string `json:"givenAt" db:"given_at"`
	Giver 		 User   `json:"giver"`
	Receiver	 User   `json:"receiver"`
}

// BeersRepositoryInterface defines the set of User related methods available
type BeersRepositoryInterface interface {
	GetBeerTransferLog(ctx context.Context) ([]BeerTransferFeedItem, error)
}

// BeersRepository implements UsersRepositoryInterface
type BeersRepository struct {
	db *sqlx.DB
}

// NewBeersRepository returns a configured BeersRepository object
func NewBeersRepository(db *sqlx.DB) *BeersRepository {
	return &BeersRepository{db: db}
}

func (r *BeersRepository) GetBeerTransferLog(ctx context.Context) ([]BeerTransferFeedItem, error) {
	query := `
		SELECT giver.id,
				giver.name,
				giver.email,
				giver.picture,
				receiver.id,
				receiver.name,
				receiver.email,
				receiver.picture,
				btf.beers,
				btf.given_at
		FROM beer_transfers btf 
		JOIN users giver ON giver.id = btf.giver_id 
		JOIN users receiver ON receiver.id = btf.taker_id
		ORDER BY btf.given_at DESC;`
	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, parseError(err)
	}

	var beerFeed []BeerTransferFeedItem
	for rows.Next() {
		var l BeerTransferFeedItem

		err = rows.Scan(
			&l.Giver.ID,
			&l.Giver.Name,
			&l.Giver.Email,
			&l.Giver.Picture,
			&l.Receiver.ID,
			&l.Receiver.Name,
			&l.Receiver.Email,
			&l.Receiver.Picture,
			&l.Beers,
			&l.GivenAt)
		beerFeed = append(beerFeed, l)
	}

	return beerFeed, nil
}
