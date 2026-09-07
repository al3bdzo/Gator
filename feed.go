package main 

import (
	"net/http"
	"encoding/xml"
	"context"
	"time"
	"io"
	"html"
	"fmt"
	"database/sql"
	"github.com/al3bdzo/Gator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title 		string    `xml:"title"`
		Link 		string 	  `xml:"link"`
		Description string 	  `xml:"description"`
		Item 		[]RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title 		string `xml:"title"`
	Link 		string `xml:"link"`
	Description string `xml:"description"`
	PubDate		string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error){
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var RssFeed RSSFeed

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err = xml.Unmarshal(data, &RssFeed); err != nil {
		return nil, err
	}

	RssFeed.Channel.Title = html.UnescapeString(RssFeed.Channel.Title)
	RssFeed.Channel.Description = html.UnescapeString(RssFeed.Channel.Description)

	for i, item := range RssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		RssFeed.Channel.Item[i] = item
	}

	return &RssFeed, nil
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	
	now := time.Now()
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		UpdatedAt: now,
		LastFetchedAt: sql.NullTime{
			Time: now,
			Valid: true,
		},
		ID: nextFeed.ID,
	})
	if err != nil {
		return err
	}

	rssFeed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}

	if len(rssFeed.Channel.Item) == 0 {
		fmt.Println("No Items in this RSS Feed")
		return nil
	}

	fmt.Printf("* Feed: %s\n\n", nextFeed.Name)
	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("* Title:	%s\n", item.Title)
	}
	fmt.Println("===============================================================")
	return nil
}