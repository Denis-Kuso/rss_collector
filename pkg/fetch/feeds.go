package fetch

import "fmt"

func CreateFeed(URL string) {
	fmt.Println("Creating feed")
}

func FollowFeed(feed string) {
	fmt.Println("Printing feed")
}

func UnfollowFeed(feed string) {
	fmt.Println("Stopped following:", feed)
}

func GetAllFollowedFeeds() {
	fmt.Println("ALL FEEDS")
}
