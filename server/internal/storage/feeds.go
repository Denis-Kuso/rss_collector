package storage

import (
	"context"
	"database/sql"

	"github.com/Denis-Kuso/rss_collector/server/internal/database"
	"github.com/google/uuid"
)

type FeedStore interface {
	Create(ctx context.Context, name, URL string) error
	Get(ctx context.Context, userID uuid.UUID) error
	Follow(ctx context.Context, userID, feedID uuid.UUID) error
	ShowAvailable(ctx context.Context) error
	Delete(ctx context.Context) error
	GetPosts(ctx context.Context, userID uuid.UUID, limit int) error
}

type Feed struct {
	// TODO
}

type Post struct {
	// TODO
}

// no better name?
type FeedsModel struct {
	DB *database.Queries
}

func NewFeedsModel(db *sql.DB) *FeedsModel {
	return &FeedsModel{DB: database.New(db)}
}

// TODO establish return values
func (f *FeedsModel) Create(ctx context.Context, name, URL string) error {
	return nil
}
func (f *FeedsModel) Get(ctx context.Context, userID uuid.UUID) error {
	return nil
}
func (f *FeedsModel) Follow(ctx context.Context, userID, feedID uuid.UUID) error { // good name?
	return nil
}
func (f *FeedsModel) ShowAvailable(ctx context.Context) error {
	return nil
}
func (f *FeedsModel) Delete(ctx context.Context) error {
	return nil
}
func (f *FeedsModel) GetPosts(ctx context.Context, userID uuid.UUID, limit int) error { // should this be in its own file
	return nil
}
