package rss

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/geophpherie/boot-dev-rss-feed-agg/internal/database"
	"github.com/google/uuid"
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
				id, err := uuid.NewV7()
				if err != nil {
					continue
				}

				// Wed, 28 Feb 2024 00:00:00 +0000
				t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", post.PubDate)
				if err != nil {
					fmt.Print("Unable to parse time", post.PubDate)
					continue
				}
				now := time.Now()
				postParams := database.CreatePostParams{
					ID:          id,
					CreatedAt:   now,
					UpdatedAt:   now,
					Title:       post.Title,
					Url:         post.Link,
					Description: sql.NullString{String: post.Description, Valid: true},
					PublishedAt: sql.NullTime{Time: t, Valid: true},
					FeedID:      feed.ID,
				}
				_, err = db.CreatePost(context.Background(), postParams)
				if err != nil {
					fmt.Println("Unable to create post ", err)
				}
			}
		}(feed.Url)

	}

	fmt.Println("Waiting on waitgroup")
	wg.Wait()
	fmt.Println("Done")
	// need wait group
	return nil
}
