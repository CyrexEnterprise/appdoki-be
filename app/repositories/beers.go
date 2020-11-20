package repositories

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type BeerTransferFeedItem struct {
	GiverId      string `json:"giverId"`
	GiverName    string `json:"giverName"`
	ReceiverId   string `json:"receiverId"`
	ReceiverName string `json:"receiverName"`
	Beers        int    `json:"beers"`
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
	beerFeed := []BeerTransferFeedItem{}

	query := `SELECT ugiv.id AS giverId, ugiv.name AS giverName, urec.id AS receiverId, urec.name AS receiverName, btf.beers 
		FROM beer_transfers btf 
		JOIN users ugiv ON ugiv.id = btf.giver_id 
		JOIN users urec ON urec.id = btf.taker_id
		ORDER BY btf.given_at DESC;`
	err := r.db.SelectContext(ctx, &beerFeed, query)
	if err != nil {
		return nil, parseError(err)
	}

	return beerFeed, nil
}
