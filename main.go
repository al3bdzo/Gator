package main

import (
	"github.com/al3bdzo/Gator/internal/config"
	"github.com/al3bdzo/Gator/internal/database"

	"log"
	"os"

	"database/sql"
	_ "github.com/lib/pq"
)

type state struct {
	db *database.Queries
	cfg *config.Config
}

func main() {
	conf, err := config.ReadJson()
	if err != nil {
		log.Fatal(err)
	}

	dbURL := conf.DbURL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	programState := &state{
		db: dbQueries,
		cfg: &conf,
	}
	programCommands := commands{
		cmds: make(map[string]func(*state, command)error),
	}
	programCommands.register("login", handlerLogin)
	programCommands.register("register", handlerRegister)
	programCommands.register("reset", handlerReset)
	programCommands.register("users", handlerUsers)

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