package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login expects to have a username")
	}
	username := cmd.args[0]
	err := (*s).config.SetUser(&username)
	if err != nil {
		return err
	}
	fmt.Printf("User (%v) has been set\n", username)
	return nil
}
