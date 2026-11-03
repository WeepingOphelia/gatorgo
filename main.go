package main

import (
	"fmt"
	"os"

	"github.com/WeepingOphelia/gatorgo/internal/config"
)

type state struct {
	Config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("ILLITERACY ERROR:", err)
		return
	}
	var s state
	s.Config = &cfg

	var cmds commands
	cmds.Commands = make(map[string]func(*state, command) error)
	cmds.register("login", handlerLogin)

	var args = os.Args
	if len(args) < 2 {
		fmt.Println("Where you trying to do something?")
		return
	}
	cmd := command{
		Name: args[1],
		Args: args[2:],
	}
	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Println("COMMAND ERROR:", err)
		return
	}

	fmt.Printf("%+v", cfg)
}
