package main

import "context"

func handlerAgg(s *state, cmd command) error {
	fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	return nil
}
