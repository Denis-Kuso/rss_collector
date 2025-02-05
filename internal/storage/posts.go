package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Denis-Kuso/rss_collector/internal/database"
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
	// uneccessarily complicated (there could be more posts than there are feeds)
	feedIDs := make([]uuid.UUID, n)
	for i, pi := range privatePosts {
		feedIDs[i] = pi.FeedID
	}
	feeds, err := p.DB.GetBasicInfoFeed(ctx, feedIDs)
	if err != nil {
		return []Post{}, err
	}
	feedInfo := make(map[uuid.UUID]string)
	for _, f := range feeds {
		feedInfo[f.ID] = f.Name
	}
	posts := make([]Post, n)

	for i, pi := range privatePosts {
		posts[i] = Post{URL: pi.Url, Title: pi.Title, FeedName: feedInfo[pi.FeedID]}
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
