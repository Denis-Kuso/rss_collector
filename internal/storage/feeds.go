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

// TODO a thought...-> if ctx is passed, do I need to provide userID (and other request-scoped vals) as separate parameter(s)?
type FeedStore interface {
	Create(ctx context.Context, userID uuid.UUID, name, URL string) error
	Get(ctx context.Context, userID uuid.UUID) ([]Feed, error)
	Follow(ctx context.Context, userID, feedID uuid.UUID) error
	ShowAvailable(ctx context.Context) ([]Feed, error)
	GetLastFetched(ctx context.Context, numFeeds int) ([]Feed, error)
	Delete(ctx context.Context, feedID, userID uuid.UUID) error
	FeedFetched(ctx context.Context, feedID uuid.UUID) error
}

type Feed struct {
	Name string
	URL  string
	ID   uuid.UUID
}

// no better name?
type FeedsModel struct {
	DB *database.Queries
}

func NewFeedsModel(db *sql.DB) *FeedsModel {
	return &FeedsModel{DB: database.New(db)}
}

func (f *FeedsModel) Create(ctx context.Context, userID uuid.UUID, name, URL string) error {
	feed, err := f.DB.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    userID,
		Name:      name,
		Url:       URL,
	})
	if err != nil {
		if pqDuplicate(err) {
			return ErrDuplicate
		}
		return err
	}
	_, err = f.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    userID,
		FeedID:    feed.ID,
	})
	if err != nil {
		if pqDuplicate(err) {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

func (f *FeedsModel) Get(ctx context.Context, userID uuid.UUID) ([]Feed, error) {
	feedFollows, err := f.DB.GetFeedFollowsForUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Feed{}, ErrNotFound
		}
		return nil, err
	}
	feedIDs := make([]uuid.UUID, len(feedFollows))
	for i, f := range feedFollows {
		feedIDs[i] = f.FeedID
	}
	dbFeeds, err := f.DB.GetBasicInfoFeed(ctx, feedIDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Feed{}, ErrNotFound
		}
		return nil, err
	}
	feeds := make([]Feed, len(dbFeeds))
	for i, fi := range dbFeeds {
		feeds[i] = Feed{
			Name: fi.Name,
			URL:  fi.Url,
			ID:   fi.ID,
		}
	}
	return feeds, nil
}

func (f *FeedsModel) Follow(ctx context.Context, userID, feedID uuid.UUID) error { // good name?
	// does desired feed even exist?
	_, err := f.DB.GetBasicInfoFeed(ctx, []uuid.UUID{feedID})
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return ErrNotFound
		}
		return err
	}
	_, err = f.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    userID,
		FeedID:    feedID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		if pqDuplicate(err) {
			return ErrDuplicate
		}
	}
	return nil
}

func (f *FeedsModel) ShowAvailable(ctx context.Context) ([]Feed, error) {
	feeds, err := f.DB.GetFeeds(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	publicFeeds := make([]Feed, len(feeds))
	for i, fs := range feeds {
		publicFeeds[i] = Feed{
			URL:  fs.Url,
			ID:   fs.ID,
			Name: fs.Name,
		}
	}
	return publicFeeds, nil
}

func (f *FeedsModel) Delete(ctx context.Context, feedID, userID uuid.UUID) error {
	err := f.DB.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		FeedID: feedID,
		UserID: userID,
	})
	if err != nil {
		// under the current "implementation" with sqlc an err should not happen
		// because if there is no response that is not an error, and the
		// sqlc-generated code ignores the sql.Result return value
		return err
	}
	return nil
}

func (f *FeedsModel) GetLastFetched(ctx context.Context, numFeeds int) ([]Feed, error) {
	feeds, err := f.DB.GetNextFeedsToFetch(ctx, int32(numFeeds))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Feed{}, nil // ok to return empty slice
		}
		return nil, err
	}

	fs := make([]Feed, len(feeds))
	for i, fi := range feeds {
		fs[i] = Feed{
			URL: fi.Url,
			ID:  fi.ID,
		}
	}

	return fs, nil
}

func (f *FeedsModel) FeedFetched(ctx context.Context, feedID uuid.UUID) error {
	_, err := f.DB.MarkFeedFetched(ctx, feedID)
	if err != nil {
		return fmt.Errorf("cannot mark feed fetched: %w", err)
	}
	return nil
}

// TODO where should I put this? in RSS? IN storage?
func transformPubTime(pubTime string) (time.Time, error) {
	const desiredFormat = time.RFC3339
	formats := []string{time.RFC822, time.RFC822Z, time.RFC1123, time.RFC850, time.RFC1123Z,
		time.DateTime, time.DateOnly, time.Stamp, "Mon, 2 Jan 2006 15:04:05 MST"} // custom format found in one of the feeds
	var timeOfPub time.Time
	var err error

	if timeOfPub, err = time.Parse(desiredFormat, pubTime); err != nil {
		// try other formats
		for _, format := range formats {
			if timeOfPub, err = time.Parse(format, pubTime); err != nil {
				continue
			}
			t := timeOfPub.Format(desiredFormat)
			_, err = time.Parse(desiredFormat, t)
			if err != nil {
				//log.Printf("failed to transform time: %s, %v\n", err, t)
			}
			break
		}
	}
	return timeOfPub, err
}
