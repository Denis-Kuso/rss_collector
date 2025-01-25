package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/Denis-Kuso/rss_collector/server/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	numFeeds = 3 // TODO pass around differently
)

// Periodically:
// fetch from DB feeds that need fetching
// fetch feeds from their URLS (concurently)
// mark feed as fetched
func worker(done <-chan struct{}, db *database.Queries, interRequestInterval time.Duration) {
	fetchTicker := time.NewTicker(interRequestInterval)
	defer fetchTicker.Stop()
	var wg sync.WaitGroup
	for {
		select {
		case <-fetchTicker.C:
			feeds, err := db.GetNextFeedsToFetch(context.Background(), numFeeds)
			if err != nil {
				slog.Error("cannot retrieve feeds", "error", err)
				continue
			}
			for _, feed := range feeds {
				wg.Add(1)
				go func(f database.Feed) {
					//
					defer wg.Done()
					defer func(cf database.Feed) {
						if r := recover(); r != nil {
							slog.Warn("panic recovered: %v", "error", r)
						}
					}(f)
					processFeed(db, f)
				}(feed)
			}
		case <-done: // server initiatied shutdown
			wg.Wait()
			slog.Info("worker stopped")
			return
		}
	}
}

func processFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Cannot't make feed %s fetched: %v\n", feed.Name, err)
		return
	}

	feedData, err := URLtoFeed(feed.Url)
	if err != nil {
		log.Printf("ERR: %v. Couldn't gather feed %s\n", err, feed.Name)
		return
	}
	for _, item := range feedData.Items {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubAt, err := transformPubTime(item.PublishedAt)
		if err != nil {
			log.Printf("ERR: %v. Post: %s. Pub time: %s\n", err, item.Link, item.PublishedAt)
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: pubAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			// ignore error if post already present
			if err, ok := err.(*pq.Error); ok {
				// unique key violation https://www.postgresql.org/docs/current/errcodes-appendix.html
				if err.Code == "23505" {
					continue
				}
			}
			log.Printf("ERR: %v. Could not create post from :%v\n", err, feed.Url)
		}
	}
	log.Printf("Scraped feed: %s, found: %d posts.\n", feed.Name, len(feedData.Items))
}

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
				log.Printf("failed to transform time: %s, %v\n", err, t)
			}
			break
		}
	}
	return timeOfPub, err
}
