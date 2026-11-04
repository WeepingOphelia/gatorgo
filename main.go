package main

import (
	"fmt"
	"os"

	"database/sql"

	"github.com/WeepingOphelia/gatorgo/internal/config"
	"github.com/WeepingOphelia/gatorgo/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("ILLITERACY ERROR:", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.Url)
	if err != nil {
		fmt.Println("DATABASE ERROR:", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	var s state
	s.cfg = &cfg
	s.db = dbQueries

	var cmds commands
	cmds.Commands = make(map[string]func(*state, command) error)
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	var args = os.Args
	if len(args) < 2 {
		fmt.Println("BOREDOM ERROR: Soooo... did you want to do something? Orrr....")
		os.Exit(1)
	}
	cmd := command{
		Name: args[1],
		Args: args[2:],
	}
	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Println("COMMAND ERROR:", err)
		os.Exit(1)
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Println("ILLITERACY ERROR:", err)
		os.Exit(1)
	}

	// fmt.Printf("%+v", cfg)
}
