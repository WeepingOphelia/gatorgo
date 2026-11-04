package main

import "errors"

type command struct {
	Name string
	Args []string
}

type commands struct {
	Commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if f, exists := c.Commands[cmd.Name]; exists {
		return f(s, cmd)
	}
	return errors.New("I'm sorry Dave, I'm afraid I can't do that.")
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.Commands[name] = f
}
