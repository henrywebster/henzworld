package internal

import (
	"fmt"
	"net/http"

	"github.com/mmcdole/gofeed"
)

type RssClient struct {
	parser *gofeed.Parser
	url    string
}

func NewRssClient(httpClient *http.Client, url string) *RssClient {
	fp := gofeed.NewParser()
	fp.Client = httpClient

	return &RssClient{
		parser: fp,
		url:    url,
	}
}

func (c *RssClient) GetRssFeed() ([]*gofeed.Item, error) {
	defer TimeFunction(fmt.Sprintf("get_rss_feed_%s", c.url))()
	feed, err := c.parser.ParseURL(c.url)
	if err != nil {
		return nil, err
	}

	return feed.Items, nil
}
