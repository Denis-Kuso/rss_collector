package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Denis-Kuso/rss_collector/server/internal/storage"
	"github.com/google/uuid"
)

const (
	numFeeds = 3 // TODO pass around differently
)

// Periodically:
// fetch from DB feeds that need fetching
// fetch feeds from their URLS (concurently)
// mark feed as fetched
func worker(done <-chan struct{}, fs storage.FeedStore, ps storage.PostStore, interRequestInterval time.Duration) {
	fetchTicker := time.NewTicker(interRequestInterval)
	ctx := context.Background() // TODO I don't like this...
	defer fetchTicker.Stop()
	var wg sync.WaitGroup
	for {
		select {
		case <-fetchTicker.C:
			feeds, err := fs.GetLastFetched(ctx, numFeeds)
			slog.Info("received numFeeds", "n", len(feeds))
			if err != nil {
				slog.Error("cannot retrieve feeds", "error", err)
				continue
				// TODO could check if there are any feeds at all
			}
			for _, feed := range feeds {
				wg.Add(1)
				go func(fs storage.FeedStore, ps storage.PostStore, URL string, feedID uuid.UUID) {
					defer wg.Done()
					defer func(URL string) {
						if r := recover(); r != nil {
							slog.Warn("panic recovered", "error", r, "URL", URL)
						}
					}(URL)
					processFeed(fs, ps, URL, feedID) // TODO also need to receive err
				}(fs, ps, feed.URL, feed.ID)
			}
		case <-done: // server initiatied shutdown
			wg.Wait()
			slog.Info("worker stopped")
			return
		}
	}
}

func processFeed(fs storage.FeedStore, ps storage.PostStore, URL string, feedID uuid.UUID) error {
	feed, err := URLtoFeed(URL) // TODO better name could be scrapeFeed
	if err != nil {
		return fmt.Errorf("cannot process feed: %v: %w", URL, err)
	}
	ctx := context.Background() // TODO will probably pass from worker
	for _, item := range feed.Items {
		err := ps.SavePost(ctx, item.Title, item.Link, item.Description, feedID, item.PublishedAt)
		if err != nil {
			if errors.Is(err, storage.ErrDuplicate) {
				continue
			}
			return fmt.Errorf("cannot process feed: %v: %w", URL, err)
		}
	}
	err = fs.FeedFetched(ctx, feedID)
	return err
}
