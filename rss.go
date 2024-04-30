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
	PRIVACY_GUDIES = "https://blog.privacyguides.org/rss/"
)

// version 3
type TempFeed struct {
	Title   string   `xml:"title"`
	Link    string `xml:"link"`
	Feeds []RSSEntry `xml:"entry,omitempty"`
	VersionChannel struct {
		Title string `xml:"title"`
		Link string `xml:"link"`
		Feeds []RSSEntry `xml:"item"`
	} `xml:"channel,omitempty"`
}


// inidvidual items 
type RSSEntry struct {
		Title string `xml:"title"`
		PublishedAt string `xml:"published"`
		PubDate string `xml:"pubDate"`
		Description string `xml:"description"`
		Link string `xml:"link"`
		Id string `xml:"id"`
}


func (rss *RSSEntry) getPubTime()(string){
	if len(rss.PublishedAt) > len(rss.PubDate){
		return rss.PublishedAt
	}else { return rss.PubDate}
}


func (rss *RSSEntry) getLink()(string){
	if rss.Link == ""{
		return rss.Id
	}else{return rss.Link}
}
// discriminates between different versions of nametags (e.g. id vs link),
// chooses non-emtpy one, populates Feed.
func (rss *TempFeed) pruneFeeds()(Feed){
	// will there always be nonzero long list of entries in only one of them?
	var tempFeeds []RSSEntry
	var feed Feed
	if len(rss.VersionChannel.Feeds) > len(rss.Feeds) {
		tempFeeds = rss.VersionChannel.Feeds
		feed.Title = rss.VersionChannel.Title
		feed.Link = rss.VersionChannel.Link
	}else {
		tempFeeds = rss.Feeds
		feed.Title = rss.Title 
		feed.Link = rss.Link 
	}
	// process items
	var tempItems []RSSitem
	for _, feed := range tempFeeds{
		pubAt := feed.getPubTime()
		link := feed.getLink()
		tempItems = append(tempItems, RSSitem{
						Title: feed.Title,
						Description: feed.Description,
						Link: link,
						PublishedAt: pubAt,
		})
	}
	feed.Items = tempItems
	return feed
}

func fetchFeed(url string) ([]byte, error){

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("ERR:%v. failed request!\n",err)
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed request: %v\n", err)
		return nil, err
	}
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Successfull response from %v, parsing response...\n", url)
	return resp, nil
}

// returns feed given a valid url
func URLtoFeed(url string) (Feed, error){
	// make request
	resp, err := fetchFeed(url)
	if err != nil{
		log.Printf("ERR during fetching URL: %v\n", err)
		return Feed{}, nil
	}
	rss := TempFeed{}
	err = xml.Unmarshal(resp, &rss)
	if err != nil {
		log.Fatalf("ERR during unmarshaling: %v\n", err)
		return Feed{}, err
	}
	feeds := rss.pruneFeeds()
	return feeds, nil
}


// Feed representing a RSS feed
type Feed struct {
	Title   string
	Description string 
	Link    string 
	Items []RSSitem
} 

// item provided by a RSS feed
type RSSitem struct {
	Title string
	Description string
	Link string
	PublishedAt string
}
