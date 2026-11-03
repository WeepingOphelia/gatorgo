package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("Username required")
	}
	s.Config.SetUser(cmd.Args[0])
	fmt.Println("User set to", cmd.Args[0])
	return nil
}
