package main

import (
	"fmt"
	"time"
	"context"
	"github.com/al3bdzo/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Usage: %s <name> <url>", cmd.name)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	now := time.Now()
	feedParams := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name: cmd.args[0],
		Url: cmd.args[1],
		UserID: user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}
	fmt.Println("Feed created successfully!")
	printFeed(feed)

	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("too many arguments!")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	for _, feed := range feeds {
		fmt.Printf("* Name:		%s\n", feed.Name)
		fmt.Printf("* URL: 		%s\n", feed.Url)
		fmt.Printf("* UserName: 	%s\n", feed.UserName)
		fmt.Println()
	}
	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}