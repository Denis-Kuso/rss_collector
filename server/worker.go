package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/Denis-Kuso/rss_collector/server/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	FEEDS_TO_FETCH = 3
)

// Periodically:
// fetch from DB feeds that need fetching
// fetch feeds from their URLS (concurently)
// mark feed as fetched
func worker(db *database.Queries, interRequestInterval time.Duration, workers int) {
	fetch_ticker := time.NewTicker(interRequestInterval)
	for ; ; <-fetch_ticker.C {
		// fetch from DB
		feeds, err := db.GetNextFeedsToFetch(context.Background(), FEEDS_TO_FETCH)
		if err != nil {
			log.Printf("ERR: %v durring retrieval of feeds from db\n", err)
		}
		var wg sync.WaitGroup
		for _, feed := range feeds {
			wg.Add(1)
			go processFeed(db, &wg, feed)
		}
		wg.Wait()
	}
}

func processFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
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
		log.Printf("Scraped feed: %s, found: %d posts.\n", feed.Name, len(feedData.Items))
	}
}

func transformPubTime(pubTime string) (time.Time, error) {
	const DESIRED_FORMAT = time.RFC3339
	FORMATS := []string{time.RFC822, time.RFC822Z, time.RFC1123, time.RFC850, time.RFC1123Z,
		time.DateTime, time.DateOnly, time.Stamp, "Mon, 2 Jan 2006 15:04:05 MST"} // custom format found in one of the feeds
	var t_pub time.Time
	var err error

	if t_pub, err = time.Parse(DESIRED_FORMAT, pubTime); err != nil {
		// try other formats
		for _, format := range FORMATS {
			if t_pub, err = time.Parse(format, pubTime); err != nil {
				continue
			} else {
				t_str := t_pub.Format(DESIRED_FORMAT)
				_, err = time.Parse(DESIRED_FORMAT, t_str)
				if err != nil {
					log.Printf("failed to transform time: %s, %v\n", err, t_str)
				}
				break
			}
		}
	}
	return t_pub, err
}
