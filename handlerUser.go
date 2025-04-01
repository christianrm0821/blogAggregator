package main

import (
	"fmt"
)

// Login function which is added to the commands map
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

// Register function which is added to the commands map
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		fmt.Print("register expects a username after command")
	}
	//username := cmd.args[0]

	//s.db.CreateUser()

	return nil
}
