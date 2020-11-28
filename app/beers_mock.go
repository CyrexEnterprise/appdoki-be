package app

import (
	repos "appdoki-be/app/repositories"
	"context"
	"github.com/brianvoe/gofakeit/v5"
	"time"
)

type mockBeersRepository struct {
	getBeerTransferImpl  func(ctx context.Context, id int) (*repos.BeerTransferFeedItem, error)
	getBeerTransfersImpl func(ctx context.Context, options *repos.BeerFeedPaginationOptions) ([]repos.BeerTransferFeedItem, error)
}

func (r *mockBeersRepository) GetBeerTransfer(ctx context.Context, id int) (*repos.BeerTransferFeedItem, error) {
	return r.getBeerTransferImpl(ctx, id)
}

func (r *mockBeersRepository) GetBeerTransfers(ctx context.Context, options *repos.BeerFeedPaginationOptions) ([]repos.BeerTransferFeedItem, error) {
	return r.getBeerTransfersImpl(ctx, options)
}

func getDefaultMockBeersRepository() *mockBeersRepository {
	return &mockBeersRepository{
		getBeerTransferImpl: func(ctx context.Context, id int) (*repos.BeerTransferFeedItem, error) {
			return generateRandomBeerTransferMock(), nil
		},
		getBeerTransfersImpl: func(ctx context.Context, options *repos.BeerFeedPaginationOptions) ([]repos.BeerTransferFeedItem, error) {
			return []repos.BeerTransferFeedItem{
				*generateRandomBeerTransferMock(),
				*generateRandomBeerTransferMock(),
				*generateRandomBeerTransferMock(),
			}, nil
		},
	}
}

func generateRandomBeerTransferMock() *repos.BeerTransferFeedItem {
	return &repos.BeerTransferFeedItem{
		Beers:    gofakeit.Number(1, 100),
		GivenAt:  time.Now().String(),
		Giver:    *generateRandomUserMock(),
		Receiver: *generateRandomUserMock(),
	}
}
