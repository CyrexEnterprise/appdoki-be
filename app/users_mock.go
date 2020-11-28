package app

import (
	repos "appdoki-be/app/repositories"
	"context"
	"github.com/brianvoe/gofakeit/v5"
	"strconv"
)

type mockUsersRepository struct {
	getAllImpl             func(ctx context.Context) ([]*repos.User, error)
	findByIDImpl           func(ctx context.Context, ID string) (*repos.User, error)
	findByEmailImpl        func(ctx context.Context, email string) (*repos.User, error)
	findOrCreateUserImpl   func(ctx context.Context, userData *repos.User) (*repos.User, bool, error)
	createImpl             func(ctx context.Context, user *repos.User) (*repos.User, error)
	updateImpl             func(ctx context.Context, user *repos.User) (*repos.User, error)
	deleteImpl             func(ctx context.Context, ID string) (bool, error)
	addBeerTransferImpl    func(ctx context.Context, giverID string, takerID string, beers int) (int, error)
	getBeerTransferLogImpl func(ctx context.Context, userID string) (*repos.UserBeerLog, error)
}

func (r *mockUsersRepository) GetAll(ctx context.Context) ([]*repos.User, error) {
	return r.getAllImpl(ctx)
}

func (r *mockUsersRepository) FindByID(ctx context.Context, ID string) (*repos.User, error) {
	return r.findByIDImpl(ctx, ID)
}

func (r *mockUsersRepository) FindByEmail(ctx context.Context, email string) (*repos.User, error) {
	return r.findByEmailImpl(ctx, email)
}

func (r *mockUsersRepository) Create(ctx context.Context, user *repos.User) (*repos.User, error) {
	return r.createImpl(ctx, user)
}

func (r *mockUsersRepository) FindOrCreateUser(ctx context.Context, userData *repos.User) (*repos.User, bool, error) {
	return r.findOrCreateUserImpl(ctx, userData)
}

func (r *mockUsersRepository) Update(ctx context.Context, user *repos.User) (*repos.User, error) {
	return r.updateImpl(ctx, user)
}

func (r *mockUsersRepository) Delete(ctx context.Context, ID string) (bool, error) {
	return r.deleteImpl(ctx, ID)
}

func (r *mockUsersRepository) AddBeerTransfer(ctx context.Context, giverID string, takerID string, beers int) (int, error) {
	return r.addBeerTransferImpl(ctx, giverID, takerID, beers)
}

func (r *mockUsersRepository) GetBeerTransfersSummary(ctx context.Context, userID string) (*repos.UserBeerLog, error) {
	return r.getBeerTransferLogImpl(ctx, userID)
}

func getDefaultMockUsersRepository() *mockUsersRepository {
	return &mockUsersRepository{
		getAllImpl: func(_ context.Context) ([]*repos.User, error) {
			return []*repos.User{generateRandomUserMock()}, nil
		},
		findByIDImpl: func(ctx context.Context, ID string) (*repos.User, error) {
			return generateRandomUserMock(), nil
		},
		findByEmailImpl: func(ctx context.Context, email string) (*repos.User, error) {
			return generateRandomUserMock(), nil
		},
		findOrCreateUserImpl: func(ctx context.Context, user *repos.User) (*repos.User, bool, error) {
			return generateRandomUserMock(), true, nil
		},
		createImpl: func(ctx context.Context, user *repos.User) (*repos.User, error) {
			return generateRandomUserMock(), nil
		},
		updateImpl: func(ctx context.Context, user *repos.User) (*repos.User, error) {
			return generateRandomUserMock(), nil
		},
		deleteImpl: func(ctx context.Context, ID string) (bool, error) {
			return true, nil
		},
		addBeerTransferImpl: func(ctx context.Context, giverID string, takerID string, beers int) (int, error) {
			return 1, nil
		},
		getBeerTransferLogImpl: func(ctx context.Context, userID string) (*repos.UserBeerLog, error) {
			return &repos.UserBeerLog{
				Given:    1,
				Received: 2,
			}, nil
		},
	}
}

func generateRandomUserMock() *repos.User {
	gofakeit.Seed(0)
	return &repos.User{
		ID:      strconv.Itoa(gofakeit.Number(0, 1000000)),
		Name:    gofakeit.Name(),
		Email:   gofakeit.Email(),
		Picture: gofakeit.URL(),
	}
}

func generateRandomUserMockWithID(id string) *repos.User {
	gofakeit.Seed(0)
	return &repos.User{
		ID:      id,
		Name:    gofakeit.Name(),
		Email:   gofakeit.Email(),
		Picture: gofakeit.URL(),
	}
}