package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

// User model
type User struct {
	ID      string `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Email   string `json:"email" db:"email"`
	Picture string `json:"picture" db:"picture"`
}

type UserBeerLog struct {
	Given    int `json:"given" db:"given"`
	Received int `json:"received" db:"received"`
}

// UsersRepositoryInterface defines the set of User related methods available
type UsersRepositoryInterface interface {
	GetAll(ctx context.Context) ([]*User, error)
	FindByID(ctx context.Context, ID string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindOrCreateUser(ctx context.Context, userData *User) (*User, bool, error)
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, ID string) (bool, error)
	AddBeerTransfer(ctx context.Context, giverID string, takerID string, beers int) (int, error)
	GetBeerTransfersSummary(ctx context.Context, userID string) (*UserBeerLog, error)
	ClearTokens(ctx context.Context, ID string) error
}

// UsersRepository implements UsersRepositoryInterface
type UsersRepository struct {
	db *sqlx.DB
}

// NewUsersRepository returns a configured UsersRepository object
func NewUsersRepository(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

// GetAll fetches all users, returns an empty slice if no user exists
func (r *UsersRepository) GetAll(ctx context.Context) ([]*User, error) {
	users := []*User{}
	err := r.db.SelectContext(ctx, &users, "SELECT id, name, email, picture FROM users")
	if err != nil {
		return nil, err
	}

	return users, nil
}

// FindByID finds a user by ID, returns nil if not found
func (r *UsersRepository) FindByID(ctx context.Context, ID string) (*User, error) {
	user := &User{}
	err := r.db.GetContext(ctx, user, "SELECT id, name, email, picture FROM users WHERE id = $1", ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// FindByEmail finds a user by email, returns nil if not found
func (r *UsersRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	stmt := "SELECT id, name, email, picture FROM users WHERE email = $1"
	err := r.db.GetContext(ctx, user, stmt, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// FindOrCreateUser finds a user by ID and creates it if not found
// returns a boolean indicating if the user was created
// TODO deal with passing txn around
func (r *UsersRepository) FindOrCreateUser(ctx context.Context, userData *User) (*User, bool, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, false, err
	}
	defer tx.Rollback()

	user := &User{}
	selectStmt := "SELECT id, name, email, picture FROM users WHERE id = $1"
	err = tx.GetContext(ctx, user, selectStmt, userData.ID)
	if err == nil {
		return user, false, nil
	}
	if err != sql.ErrNoRows {
		return nil, false, parseError(err)
	}

	insertStmt := "INSERT INTO users (id, name, email, picture) VALUES ($1, $2, $3, $4)"
	res, err := tx.ExecContext(ctx, insertStmt, userData.ID, userData.Name, userData.Email, userData.Picture)
	if err != nil {
		return nil, false, parseError(err)
	}

	if rows, err := res.RowsAffected(); err != nil {
		if rows == 0 {
			return nil, false, errors.New("could not create user")
		}
		return nil, false, parseError(err)
	}

	err = tx.GetContext(ctx, user, selectStmt, userData.ID)
	if err != nil {
		return nil, false, parseError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, false, parseError(err)
	}

	return user, true, nil
}

// Create creates a new user, returning the full model
func (r *UsersRepository) Create(ctx context.Context, user *User) (*User, error) {
	stmt := "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id"
	row := r.db.QueryRowxContext(ctx, stmt, user.Name, user.Email)
	err := row.Scan(&user.ID)
	if err != nil {
		return nil, parseError(err)
	}
	return user, nil
}

// Update updates a user, returning the updated model or nil if no rows were affected
func (r *UsersRepository) Update(ctx context.Context, user *User) (*User, error) {
	stmt := "UPDATE users SET name = $1, email = $2 WHERE id = $3"
	res, err := r.db.ExecContext(ctx, stmt, user.Name, user.Email, user.ID)
	if err != nil {
		return nil, parseError(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, nil
	}
	return user, nil
}

// Delete deletes a user, only returns error if action fails
func (r *UsersRepository) Delete(ctx context.Context, ID string) (bool, error) {
	stmt := "DELETE FROM users WHERE id = $1 RETURNING id"
	res, err := r.db.ExecContext(ctx, stmt, ID)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (r *UsersRepository) AddBeerTransfer(ctx context.Context, giverID string, takerID string, beers int) (int, error) {
	stmt := "INSERT INTO beer_transfers (giver_id, taker_id, beers) VALUES ($1, $2, $3) RETURNING id"
	var newID int
	err := r.db.GetContext(ctx, &newID, stmt, giverID, takerID, beers)
	if err != nil {
		return 0, parseError(err)
	}

	return newID, nil
}

func (r *UsersRepository) GetBeerTransfersSummary(ctx context.Context, userID string) (*UserBeerLog, error) {
	beerLog := &UserBeerLog{}

	giverQuery := "SELECT COALESCE(SUM(beers), 0) AS given FROM beer_transfers WHERE giver_id = $1"
	err := r.db.GetContext(ctx, beerLog, giverQuery, userID)
	if err != nil {
		return nil, parseError(err)
	}

	receivedQuery := "SELECT COALESCE(SUM(beers), 0) AS received FROM beer_transfers WHERE taker_id = $1"
	err = r.db.GetContext(ctx, beerLog, receivedQuery, userID)
	if err != nil {
		return nil, parseError(err)
	}

	return beerLog, nil
}

func (r *UsersRepository) ClearTokens(ctx context.Context, userID string) error {
	stmt := "UPDATE users SET oidc_refresh_token = NULL WHERE id = $1"
	_, err := r.db.ExecContext(ctx, stmt, userID)
	if err != nil {
		return parseError(err)
	}

	return nil
}
