package main 

import (
	"net/http"
	"encoding/xml"
	"context"
	"time"
	"io"
	"html"
	"fmt"
	"log"
	"database/sql"
	"github.com/al3bdzo/Gator/internal/database"
	"github.com/google/uuid"
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
	
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		LastFetchedAt: sql.NullTime{
			Time: time.Now(),
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

	for _, item := range rssFeed.Channel.Item {
		var publishedAt sql.NullTime

		parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err == nil {
			publishedAt = sql.NullTime{
				Time:  parsedTime,
				Valid: true,
			}
		}

		err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: item.Title,
			Url: item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid: true,
			},
			PublishedAt: publishedAt,
			FeedID: nextFeed.ID,
		})
		
		if err != nil {
			log.Printf("Couldn't Create Post: %v", err)
			continue
		}
	}

	fmt.Printf("Feed %s Collected, %v posts found\n", nextFeed.Name, len(rssFeed.Channel.Item))
	return nil
}