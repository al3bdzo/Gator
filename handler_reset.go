package main

import(
	"fmt"
	"context"
)


func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("DataBase was reset successfully!")
	return nil
}