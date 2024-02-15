package main

import (
	"context"
	"database/sql"
	"strings"
	"log"
	"sync"
	"time"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"github.com/google/uuid"
)

const (
	FEEDS_TO_FETCH = 3;
	POST_ALREADY_PRESENT = "duplicate key"
)
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
	for _, item := range feedData.Items {
		description := sql.NullString{}
		if item.Description != ""{
			description.String = item.Description
			description.Valid = true
		}
		
		pubAt, err := parsePubTime(item.PublishedAt);
		if err != nil {
			log.Printf("ERR: %v. Cannot parse pub time\n", err)
			// continue
		}
		// DO I WANT TO LOG THE POST THAT WAS CREATED?
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title: feed.Name,
			Url: feed.Url,
			Description: description,
			PublishedAt: pubAt,
			FeedID: feed.ID,
			})
		if err != nil {
			// ignore error if post already present
			if strings.Contains(err.Error(), POST_ALREADY_PRESENT){
				continue
			}
			log.Printf("ERR: %v. Could not create post.\n", err)
		}
		//log.Printf("Found post: %v\n", item.Title)
	}
	log.Printf("Scraped feed: %s, found: %d posts.\n", feed.Name, len(feedData.Items))
}

func parsePubTime(pubAtTime string) (time.Time, error){
	// TODO handle multiple different formats
	t, err := time.Parse(time.RFC3339,pubAtTime)
	if err != nil{
		return time.Time{}, err
	}
	return t,nil
}
