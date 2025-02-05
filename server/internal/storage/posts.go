package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Denis-Kuso/rss_collector/server/internal/database"
	"github.com/google/uuid"
)

type PostStore interface {
	Get(ctx context.Context, userID uuid.UUID, limit int) ([]Post, error)
	SavePost(ctx context.Context, title, URI, desc string, feedID uuid.UUID, pubAt string) error
}

type PostModel struct {
	DB *database.Queries
}

type Post struct {
	FeedName string `json:"feedName"`
	Title    string `json:"title"`
	URL      string `json:"url"`
}

func NewPostsModel(db *sql.DB) *PostModel {
	return &PostModel{DB: database.New(db)}
}

func (p *PostModel) Get(ctx context.Context, userID uuid.UUID, limit int) ([]Post, error) {
	privatePosts, err := p.DB.GetPostsFromUser(ctx, database.GetPostsFromUserParams{
		UserID: userID,
		Limit:  int32(limit),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Post{}, ErrNotFound
		}
		return []Post{}, err
	}
	n := len(privatePosts)
	feedIDs := make([]uuid.UUID, n)
	feeds, err := p.DB.GetBasicInfoFeed(ctx, feedIDs)
	if err != nil {
		return []Post{}, err
	}

	posts := make([]Post, n)
	for i, pi := range privatePosts {
		posts[i] = Post{URL: pi.Url, Title: pi.Title, FeedName: feeds[i].Name}
	}
	return posts, nil
}

func (p *PostModel) SavePost(ctx context.Context, title, URI, desc string, feedID uuid.UUID, pubAt string) error {
	description := sql.NullString{}
	if desc != "" {
		description.String = desc
		description.Valid = true
	}

	tPublished, err := transformPubTime(pubAt)
	if err != nil {
		return fmt.Errorf("cannot transform publication time: %q", pubAt)
	}
	_, err = p.DB.CreatePost(ctx, database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Title:       title,
		Url:         URI,
		Description: description,
		PublishedAt: tPublished,
		FeedID:      feedID,
	})
	if err != nil {
		if pqDuplicate(err) {
			return ErrDuplicate
		}
	}
	return err
}
