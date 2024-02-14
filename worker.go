package main

import (
	"context"
	"log"
	"sync"
	"time"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
)

const FEEDS_TO_FETCH = 3;
// Periodically:
// fetch from DB feeds that need fetching (already have function)
// fetch feeds from their URLS (concurently)
// mark feed as fetched 
// what if time.Duration provided is "too low"?
func worker(db *database.Queries, interRequestInterval time.Duration, workers int) {
	fetch_ticker := time.NewTicker(interRequestInterval)
	for ; ; <- fetch_ticker.C {
		// fetch from DB 
		log.Println("HJELLOO from worker")
		feeds, err := db.GetNextFeedsToFetch(context.Background(),FEEDS_TO_FETCH)
		if err != nil {
			log.Printf("ERR: %v durring retrieval of feeds from db\n",err)
		}
		// get url from feeds and show it
		var wg sync.WaitGroup
		for _, feed := range feeds {
			wg.Add(1)
			go processFeed(db, &wg, feed)
		}
		wg.Wait()
	}
}


//perhaps don't expose database.Feed (althouth this is a private method)
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
	for _, item := range feedData.Feed1.Feeds {
		log.Printf("Found post: %v\n", item.Title)
	}
	log.Printf("Feed %s gathered, found %v posts\n", feed.Name, len(feedData.Feed1.Feeds))
}
