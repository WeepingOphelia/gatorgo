package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/WeepingOphelia/gatorgo/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, u database.User) error {
	ctx := context.Background()
	if len(cmd.Args) < 1 {
		return errors.New("Feed must have nomenclature")
	}
	if len(cmd.Args) < 2 {
		return errors.New("Where is the feed?")
	}

	nf, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    u.ID})
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    u.ID,
		FeedID:    nf.ID})
	if err != nil {
		return err
	}

	fmt.Printf("Added %s\n", nf.Name)
	fmt.Printf("ID: %v\n", nf.ID)
	fmt.Printf("URL: %s\n", nf.Url)
	fmt.Printf("Created: %v\n", nf.CreatedAt)
	fmt.Printf("Updated: %v\n", nf.UpdatedAt)
	fmt.Printf("Added by: %s\n", u.Name)
	return nil

}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Print(feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()
	fs, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}
	for _, f := range fs {
		fmt.Printf("%s\n", f.Name)
		fmt.Printf("%s\n", f.Url)
		u, err := s.db.GetUserById(ctx, f.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Added by: %s\n\n", u.Name)
	}

	return nil
}

func handlerFollow(s *state, cmd command, u database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("Follow what exactly?")
	}

	ctx := context.Background()

	f, err := s.db.GetFeedByUrl(ctx, cmd.Args[0])
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    u.ID,
		FeedID:    f.ID})
	if err != nil {
		return err
	}

	fmt.Printf("Added feed follow:\nUser: %v\nFeed: %v\n", u.Name, f.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, u database.User) error {
	fs, err := s.db.GetFeedFollowsForUser(context.Background(), u.Name)
	if err != nil {
		return err
	}

	fmt.Printf("%s follows:\n", u.Name)
	for _, f := range fs {
		fmt.Printf("%s\n", f.Name)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, u database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("Unfollow what?")
	}
	ctx := context.Background()

	f, err := s.db.GetFeedByUrl(ctx, cmd.Args[0])
	if err != nil {
		return err
	}

	return s.db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		FeedID: f.ID,
		UserID: u.ID,
	})
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		u, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, u)
	}
}
