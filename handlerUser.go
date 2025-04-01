package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"workspace/github.com/christianrm0821/blogAggregator/internal/database"

	"github.com/google/uuid"
)

// Login function which is added to the commands map
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login expects to have a username")
	}

	//get username and make a new user with it. Handle error if username does not exist
	username := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Printf("User not found: %v\n", username)
		os.Exit(1)
		return err
	}

	//set the current user as the user we found
	err = (*s).config.SetUser(&user.Name)
	if err != nil {
		os.Exit(1)
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

	//Get the username, create an empty context, create new uuid
	username := cmd.args[0]
	ctx := context.Background()
	id := uuid.New()

	//Create a new user struct
	newUser := database.CreateUserParams{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	//create a new User and checking if there is an error
	user, err := s.db.CreateUser(ctx, newUser)
	if err != nil {
		os.Exit(1)
		return fmt.Errorf("%v is already in the database", username)
	}

	//Setting the current user to the username you just created
	//Handle if there is an error
	err = s.config.SetUser(&user.Name)
	if err != nil {
		return err
	}

	//Print all of the new users information
	fmt.Printf("user ID: %v\n", user.ID)
	fmt.Printf("user Created At: %v\n", user.CreatedAt)
	fmt.Printf("user updated at: %v\n", user.UpdatedAt)
	fmt.Printf("user name: %v\n", user.Name)

	return nil
}

// Deletes all of the rows(users) in the Users table
func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.db.ResetUsers(ctx)
	if err != nil {
		fmt.Println("could not reset the users table")
		os.Exit(1)
		return err
	}
	return nil
}

func handlerUserList(s *state, cmd command) error {
	ctx := context.Background()
	names, err := s.db.ListUsers(ctx)
	if err != nil {
		os.Exit(1)
		return err
	}
	currUser := *(s.config.CurrentUserName)
	for _, val := range names {
		if val == currUser {
			val = val + " (current)"
		}
		fmt.Printf("* %v\n", val)
	}
	return nil
}
