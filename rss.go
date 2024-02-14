package main 

import (
	"net/http"
	"encoding/xml"
	"log"
	"io"
)

const (
	GEOHOT_BLOG = "https://geohot.github.io/blog/feed.xml"
	LANES_BLOG = "https://wagslane.dev/index.xml"
	BOOTDEV_BLOG = "https://blog.boot.dev/index.xml"
	LESSWRONG_BLOG = "https://www.lesswrong.com/feed.xml?view=curated-rss"
	BURRITOS = "https://sideofburritos.com/blog/atom.xml"
	NEW_NEWS_WIRE = "https://netnewswire.blog/feed.xml"
	HARDCORE_HISTORY = "https://feeds.feedburner.com/dancarlin/history"
	AI_UNRAVELED = "https://media.rss.com/djamgatech/feed.xml"
	LEX_FRIDMAN = "https://podcastaddict.com/podcast/lex-fridman-podcast/3041340#"
	XKCD = "https://xkcd.com/rss.xml"
	BEST_PRACTICES = "https://www.bestpractices.dev/en/feed"
)



// Problem - there seem to be at least two "formats" for blogs
// Define more types or have xml on null conditon within your struct
// definition
// but then I cannot acces feeds, because I am referenceing field feed 
// or if it is empty I need to reference feed1
// I could write a method that checks which one is empty and return the
// non empty one
type RSSfeed1 struct {
	Title   string   `xml:"title"`
	ID    string `xml:"id"`
	Feeds []RSSEntry `xml:"entry,omitempty"`
//	Feed struct {
//		Title string `xml:"title"`
//		Id string `xml:"id"`
//		Feeds []RSSEntry `xml:"entry"`
//	} `xml:"feed,omitempty"`
	Feed1 struct {
		Title string `xml:"title"`
		Id string `xml:"id"`
		Feeds []RSSEntry `xml:"item"`
	} `xml:"channel,omitempty"`
}

// Structure for inidvidual items (articles) on a blog
// content string
type RSSEntry struct {
		Title string `xml:"title"`
		PublishedAt string `xml:"published"`
		UpdatedAt string `xml:"updated"`
		Id string `xml:"id"`
		Link string `xml:"link"`
		Content struct {
			Text string `xml:"chardata"`
			Type string `xml:"type"`
		} `xml:"content,omiteempty"`
}

//func (r *RSSfeed1) isAlternateFormat() bool {
//	if len(r.Feed.Feeds) > 0 {
//		return true
//	}else
//	{return false
//	}
//}

//type RSSfeed struct {
//	Feed struct {
//		Title string `json:"title"`
//		Link string `json:"link"`
//		Items []RSSItem `xml:"item"`
//	} `xml:"channel"`
//}
//
//type RSSItem struct {
//	Title string `xml:"title"`
//	Link string `xml:"link"`
//	PublishedAt string `xml:"published"`
//	UpdatedAt string `xml:"updated"`
//	Content string `xml:"description"`
//}


// returns feed given a valid xml url
// reconsider what this functions returns - perhaps call
// a helper function that depending on which RSS structure was unmarshaled
// return the proper go structure
func URLtoFeed(url string) (RSSfeed1, error){
	// make request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return RSSfeed1{}, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed response: %v\n", err)
		return RSSfeed1{}, err
	}
	defer res.Body.Close()
	resp, err := io.ReadAll(res.Body)
	log.Printf("Successfull response from %v, parsing response...\n", url)
	rss := RSSfeed1{}
	err = xml.Unmarshal(resp, &rss)
	if err != nil {
		log.Fatalf("ERR during unmarshaling: %v\n", err)
		return RSSfeed1{}, err
	}
	return rss, nil
}

// 
type Feed struct {
	Title   string   `xml:"title"`
	ID    string `xml:"id"`
	Entry []struct {
		Text      string `xml:",chardata"`
		Title     string `xml:"title"`
		ID      string `xml:"id"`
		Content struct {
			Text string `xml:",chardata"`
		} `xml:"content"`
	} `xml:"entry"`
} 
