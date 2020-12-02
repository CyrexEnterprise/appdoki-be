package repositories

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type BeerTransferFeedItem struct {
	ID       int    `json:"id"`
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
	GetBeerTransfer(ctx context.Context, id int) (*BeerTransferFeedItem, error)
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

const baseBeerTransferQuery = `
	SELECT giver.id,
			giver.name,
			giver.email,
			giver.picture,
			receiver.id,
			receiver.name,
			receiver.email,
			receiver.picture,
			btf.beers,
			btf.given_at,
			btf.id
	FROM beer_transfers btf 
	JOIN users giver ON giver.id = btf.giver_id 
	JOIN users receiver ON receiver.id = btf.taker_id
`

func (r *BeersRepository) GetBeerTransfer(ctx context.Context, id int) (*BeerTransferFeedItem, error) {
	query := baseBeerTransferQuery + " WHERE btf.id = $1;"
	row := r.db.QueryRowxContext(ctx, query, id)

	var t BeerTransferFeedItem
	err := row.Scan(
		&t.Giver.ID,
		&t.Giver.Name,
		&t.Giver.Email,
		&t.Giver.Picture,
		&t.Receiver.ID,
		&t.Receiver.Name,
		&t.Receiver.Email,
		&t.Receiver.Picture,
		&t.Beers,
		&t.GivenAt,
		&t.ID)

	if err != nil {
		return nil, parseError(err)
	}

	return &t, nil
}

func (r *BeersRepository) GetBeerTransfers(ctx context.Context, options *BeerFeedPaginationOptions) ([]BeerTransferFeedItem, error) {
	var whereClause string
	var limitClause string

	if len(options.GivenAt) > 0 {
		limitClause = fmt.Sprintf(" LIMIT %d", options.Limit)
		whereClause = fmt.Sprintf(" WHERE given_at %s $1", options.op)
	}

	query := fmt.Sprintf("%s %s ORDER BY btf.given_at DESC %s;", baseBeerTransferQuery, whereClause, limitClause)

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
		var t BeerTransferFeedItem

		err = rows.Scan(
			&t.Giver.ID,
			&t.Giver.Name,
			&t.Giver.Email,
			&t.Giver.Picture,
			&t.Receiver.ID,
			&t.Receiver.Name,
			&t.Receiver.Email,
			&t.Receiver.Picture,
			&t.Beers,
			&t.GivenAt,
			&t.ID)
		beerFeed = append(beerFeed, t)
	}

	return beerFeed, nil
}
