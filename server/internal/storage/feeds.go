package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/Denis-Kuso/rss_collector/server/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type FeedStore interface {
	Create(ctx context.Context, userID uuid.UUID, name, URL string) error
	Get(ctx context.Context, userID uuid.UUID) error
	Follow(ctx context.Context, userID, feedID uuid.UUID) error
	ShowAvailable(ctx context.Context) error
	GetLastFetched(ctx context.Context, numFeeds int) ([]Feed, error)
	Delete(ctx context.Context) error
	GetPosts(ctx context.Context, userID uuid.UUID, limit int) error
	SavePost(ctx context.Context, title, URI, desc string, feedID uuid.UUID, pubAt string) error
}

type Feed struct {
	Name string
	URL  string
	ID   uuid.UUID
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
func (f *FeedsModel) GetLastFetched(ctx context.Context, numFeeds int) ([]Feed, error) {
	return []Feed{}, nil
}

func (f *FeedsModel) GetPosts(ctx context.Context, userID uuid.UUID, limit int) error { // should this be in its own file
	return nil
}
func (f *FeedsModel) SavePost(ctx context.Context, title, URI, desc string, feedID uuid.UUID, pubAt string) error { // should this be in its own file
	// TODO should pubAt be time.Time?
	// should this be here?(mark fetched)
	_, err := f.DB.MarkFeedFetched(context.Background(), feedID)
	if err != nil {
		//log.Printf("Cannot't make feed %s fetched: %v\n", feed.Name, err)
		return err
	}
	description := sql.NullString{}
	if desc != "" {
		description.String = desc
		description.Valid = true
	}

	tPublished, err := transformPubTime(pubAt)
	// well this is no good
	if err != nil {
		//log.Printf("ERR: %v. Post: %s. Pub time: %s\n", err)
	}
	_, err = f.DB.CreatePost(context.Background(), database.CreatePostParams{
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
		// ignore error if post already present
		if err, ok := err.(*pq.Error); ok {
			// unique key violation https://www.postgresql.org/docs/current/errcodes-appendix.html
			if err.Code == "23505" {
				return err // TODO duplicate
			}
		}
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
