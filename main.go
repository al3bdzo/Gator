package main

import (
	"github.com/al3bdzo/Gator/internal/config"
	"github.com/al3bdzo/Gator/internal/database"

	"log"
	"os"
	"context"

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
	programCommands.register("agg", handlerAgg)
	programCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	programCommands.register("feeds", handlerGetFeeds)
	programCommands.register("follow", middlewareLoggedIn(handlerFollow))
	programCommands.register("following", middlewareLoggedIn(handlerFollowing))
	programCommands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	programCommands.register("browse", middlewareLoggedIn(handlerBrowse))

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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}