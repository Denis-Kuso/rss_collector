package storage

import (
	"context"
	"database/sql"
	"github.com/Denis-Kuso/rss_collector/server/internal/database"
	"time"

	"github.com/google/uuid"
)

// TODO - some custom error types?

// poor name
type UserStore interface {
	Create(ctx context.Context, username string) (User, error)
	Get() (User, error)
}

type UsersModel struct {
	DB *database.Queries
}

func NewUsersModel(db *sql.DB) *UsersModel {
	u := new(UsersModel)
	u.DB = database.New(db)
	return u
}

// public instance of User
type User struct {
	// what ever fields we promise in public API
	//e.g.
	Name  string `json:"username"`
	Feeds []Feed `json:"followedFeeds,omitempty"`
}

// create user
// no "validation" here
// db messing about
// ...
// ...
// return public instance of model
func (m *UsersModel) Create(ctx context.Context, username string) (User, error) {
	// e.g.
	u, err := m.DB.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      username,
	})
	if err != nil {
		// check what it might be
	}
	user := User{Name: u.Name}

	return user, nil
}

// returns info about user
func (m *UsersModel) Get() (User, error) {
	return User{}, nil
}
