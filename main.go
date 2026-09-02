package main

import (
	"github.com/al3bdzo/Gator/internal/config"
	"log"
	"os"
)

type state struct {
	cfg *config.Config
}

func main() {
	conf, err := config.ReadJson()
	if err != nil {
		log.Fatal(err)
	}

	programState := &state{
		cfg: &conf,
	}
	programCommands := commands{
		cmds: make(map[string]func(*state, command)error),
	}
	programCommands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = programCommands.run(programState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}