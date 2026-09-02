package main 

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Usage: %s <name>", cmd.name)
	}
	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set the current user: %w", err)
	}
	fmt.Println("Login Successful!")
	return nil
}
