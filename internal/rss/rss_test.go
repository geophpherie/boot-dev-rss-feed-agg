package rss

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {

	feed, _ := parseRssFeed("https://blog.boot.dev/index.xml")
	fmt.Print(feed)
}
