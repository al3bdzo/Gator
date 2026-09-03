package main 

import (
	"fmt"
	"time"
	"context"
	"github.com/al3bdzo/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Usage: %s <name>", cmd.name)
	}
	name := cmd.args[0]

	user, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("Error getting user: %v", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set the current user: %w", err)
	}

	fmt.Println("Login Successful!")
	fmt.Printf("user data: %v\n", user)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <name>", cmd.name)
	}

	name := cmd.args[0]
	now := time.Now()
	userParams := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name: name,
	}

	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("problem creating user: %v", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}

	fmt.Printf("user: %s was registered successfully!\ndata: %v\n", name, user)
	return nil
}


func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	currentUser := s.cfg.CurrentUserName
	for _, user := range users {
		if user.Name == currentUser{
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}