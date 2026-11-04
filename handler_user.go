package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/WeepingOphelia/gatorgo/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("Username required for login")
	}
	if _, err := s.db.GetUser(context.Background(), cmd.Args[0]); err != nil {
		return errors.New("Non-existant user: " + cmd.Args[0])
	}
	s.cfg.SetUser(cmd.Args[0])
	fmt.Println("User set to", cmd.Args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("Username required for registration")
	}
	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err == nil {
		return errors.New("User already exists")
	}
	s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.Args[0],
		},
	)
	s.cfg.SetUser(cmd.Args[0])
	fmt.Println("Registered:", cmd.Args[0])
	fmt.Print(s.db.GetUser(context.Background(), cmd.Args[0]))
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("User reset")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	usr, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, u := range usr {
		fmt.Printf("* %s", u.Name)
		if u.Name == s.cfg.CurrentUserName {
			fmt.Printf(" (current)\n")
		} else {
			fmt.Print("\n")
		}
	}
	return nil
}
