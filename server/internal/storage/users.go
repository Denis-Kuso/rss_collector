package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Denis-Kuso/rss_collector/server/internal/database"
	"github.com/lib/pq"

	"github.com/google/uuid"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("already exists")
)

type UserStore interface {
	Create(context.Context, string) (User, error)
	Get(context.Context, string) (User, error)
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
	Name   string `json:"username"`
	Feeds  []Feed `json:"followedFeeds,omitempty"`
	APIkey string `json:"APIkey,omitempty"`
}

// create user
// return public instance of model
func (m *UsersModel) Create(ctx context.Context, username string) (User, error) {
	u, err := m.DB.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      username,
	})

	if err != nil {
		if pqDuplicate(err) {
			return User{}, ErrDuplicate
		}
		return User{}, err
	}
	return User{Name: u.Name, APIkey: u.ApiKey}, nil
}

// returns info about user
func (m *UsersModel) Get(ctx context.Context, APIkey string) (User, error) {

	user, err := m.DB.GetUserByAPIKey(ctx, APIkey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, nil
	}

	feedFollows, err := m.DB.GetFeedFollowsForUser(ctx, user.ID)
	// ErrNoRows is acceptable, since the user might not yet follow anything
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("cannot retrieve followed feeds: user: %v: %w", user.ID, err)
		return User{}, err
	}

	numFeeds := len(feedFollows)
	// does not follow anything yet
	if numFeeds == 0 {
		return User{Name: user.Name}, nil
	}

	feedIDs := make([]uuid.UUID, numFeeds)
	for i, f := range feedFollows {
		feedIDs[i] = f.FeedID
	}
	feeds, err := m.DB.GetBasicInfoFeed(ctx, feedIDs)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("cannot retrieve feed info: %v", err)
		return User{}, nil
	}
	var userFeeds []Feed
	for _, f := range feeds {
		userFeeds = append(userFeeds, Feed{Name: f.Name,
			ID:  f.ID,
			URL: f.Url, //TODO rename field to all capital
		})
	}

	return User{Name: user.Name, Feeds: userFeeds}, nil
}

// pqDuplicate returns true if an error corresponds
// to uniqueKey violation (duplicate) and false otherwise.
// See more: unique key violation https://www.postgresql.org/docs/current/errcodes-appendix.html
func pqDuplicate(err error) bool {
	if err, ok := err.(*pq.Error); ok {
		if err.Code == "23505" {
			return true
		}
	}
	return false
}
