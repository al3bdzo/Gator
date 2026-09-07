package main

import (
	"fmt"
	"time"
	"context"
	"github.com/al3bdzo/Gator/internal/database"
	"strconv"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <time_between_requests>", cmd.name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every: %v\n", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int64
	if len(cmd.args) == 1 {
		limit, _ = strconv.ParseInt(cmd.args[0], 10, 32)
	} else {
		limit = 2
	}

	lim := int32(limit)
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: lim,
	})
	if err != nil {
		return err
	}
	
	for _, post := range posts {
		printPost(post)
		fmt.Println("========================================================")
	}
	fmt.Println()
	return nil 
}

func printPost(post database.Post) {
	fmt.Printf("* ID:	        %v\n", post.ID)
	fmt.Printf("* Title:        %v\n", post.Title)
	fmt.Printf("* Url:	        %v\n", post.Url)
	fmt.Printf("* Description:	%v\n", post.Description.String)
	fmt.Printf("* Published At: %v\n", post.PublishedAt.Time)
}