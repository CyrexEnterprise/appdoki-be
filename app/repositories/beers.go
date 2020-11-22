package repositories

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type BeerTransferFeedItem struct {
	Beers    int    `json:"beers"`
	GivenAt  string `json:"givenAt" db:"given_at"`
	Giver    User   `json:"giver"`
	Receiver User   `json:"receiver"`
}

type BeerFeedPaginationOptions struct {
	Limit   int
	GivenAt string
	op      string
}

func (o *BeerFeedPaginationOptions) SetGtOperator() {
	o.op = ">"
}

func (o *BeerFeedPaginationOptions) SetLtOperator() {
	o.op = "<"
}

// BeersRepositoryInterface defines the set of User related methods available
type BeersRepositoryInterface interface {
	GetBeerTransfers(ctx context.Context, options *BeerFeedPaginationOptions) ([]BeerTransferFeedItem, error)
}

// BeersRepository implements UsersRepositoryInterface
type BeersRepository struct {
	db *sqlx.DB
}

// NewBeersRepository returns a configured BeersRepository object
func NewBeersRepository(db *sqlx.DB) *BeersRepository {
	return &BeersRepository{db: db}
}

func (r *BeersRepository) GetBeerTransfers(ctx context.Context, options *BeerFeedPaginationOptions) ([]BeerTransferFeedItem, error) {
	baseQuery := `
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
	`

	var whereClause string
	var limitClause string

	if len(options.GivenAt) > 0 {
		limitClause = fmt.Sprintf(" LIMIT %d", options.Limit)
		whereClause = fmt.Sprintf(" WHERE given_at %s $1", options.op)
	}

	query := fmt.Sprintf("%s %s ORDER BY btf.given_at DESC %s;", baseQuery, whereClause, limitClause)

	var rows *sqlx.Rows
	var err error
	if len(options.GivenAt) > 0 {
		rows, err = r.db.QueryxContext(ctx, query, options.GivenAt)
	} else {
		rows, err = r.db.QueryxContext(ctx, query)
	}
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
