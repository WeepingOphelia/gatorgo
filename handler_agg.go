package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/WeepingOphelia/gatorgo/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	ctx := context.Background()
	if len(cmd.Args) < 1 {
		return errors.New("Feed must have nomenclature")
	}
	if len(cmd.Args) < 2 {
		return errors.New("Where is the feed?")
	}

	u, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return errors.New("Rogue user.")
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

	fmt.Print(nf)
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
