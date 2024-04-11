// convert database models to public api representation
package main

import (
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"github.com/google/uuid"
)

type PublicUser struct {
	Name   string       `json:"name"`
	ApiKey string       `json:"apiKey"`
	Feeds  []PublicFeed `json:"followedFeeds,omitempty"`
}

func dbUserToPublicUser(user database.User, feeds []database.Feed) PublicUser {
	f := dbFeedToPublicFeeds(feeds)
	return PublicUser{
		Name:   user.Name,
		ApiKey: user.ApiKey,
		Feeds:  f,
	}
}

type PublicFeed struct {
	Name string    `json:"name"`
	URL  string    `json:"url"`
	ID   uuid.UUID `json:"id"`
}

func dbFeedToPublicFeed(feed database.Feed) PublicFeed {
	return PublicFeed{
		Name: feed.Name,
		URL:  feed.Url,
		ID:   feed.ID,
	}
}
func dbFeedToPublicFeeds(feeds []database.Feed) []PublicFeed {
	SIZE := len(feeds)
	f := make([]PublicFeed, SIZE)
	for i, feed := range feeds {
		f[i] = dbFeedToPublicFeed(feed)
	}
	return f
}

type PublicPost struct {
	FeedName string `json:"feedName"`
	Title    string `json:"title"`
	URL      string `json:"url"`
}

func dbPostToPublicPost(post database.Post, feed database.Feed) PublicPost {
	f := dbFeedToPublicFeed(feed)
	return PublicPost{
		FeedName: f.Name,
		Title:    post.Title,
		URL:      post.Url,
	}
}

func dbPostsToPublicPosts(posts []database.Post, feeds []database.Feed) []PublicPost{
	SIZE := len(posts)
	p := make([]PublicPost, SIZE)
	for i := 0; i < SIZE; i++ {
		p[i] = dbPostToPublicPost(posts[i], feeds[i])
	}
	return p
}
