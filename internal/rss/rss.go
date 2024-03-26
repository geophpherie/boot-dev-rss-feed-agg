package rss

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/geophpherie/boot-dev-rss-feed-agg/internal/database"
)

type RssFeed struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Description   string `xml:"description"`
		Generator     string `xml:"generator"`
		Language      string `xml:"language"`
		LastBuildDate string `xml:"lastBuildDate"`
		Item          []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Guid        string `xml:"guid"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

func ParseRssFeed(url string) (*RssFeed, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("bad request")
	}

	rssXml := RssFeed{}
	decoder := xml.NewDecoder(resp.Body)
	decoder.Decode(&rssXml)

	return &rssXml, nil
}

func ScrapeFeeds(db *database.Queries) {
	c := time.Tick(time.Second * 10)
	for {
		fmt.Println("waiting")
		<-c
		fmt.Println("Starting Call")
		go ProcessRssFeeds(db)
		fmt.Println("Ending Call")
	}
}
func ProcessRssFeeds(db *database.Queries) error {
	feeds, err := db.GetNextFeedsToFetch(context.Background(), 10)
	if err != nil {
		return err
	}
	fmt.Println("Starting waitgroup")
	var wg sync.WaitGroup
	// this should be a go routine
	for _, feed := range feeds {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			rssFeed, err := ParseRssFeed(feed.Url)
			if err != nil {
				fmt.Printf("Error parsing %v", feed.Url)
				return
			}

			params := database.MarkFeedFetchedParams{
				CreatedAt: time.Now().UTC(),
				ID:        feed.ID,
			}

			db.MarkFeedFetched(context.Background(), params)
			for _, post := range rssFeed.Channel.Item {
				fmt.Println(post.Title)
			}
		}(feed.Url)

	}

	fmt.Println("Waiting on waitgroup")
	wg.Wait()
	fmt.Println("Done")
	// need wait group
	return nil
}
