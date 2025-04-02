package main

import (
	"context"
	"fmt"
	"html"
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
	ctx = context.Background()
	err = s.db.ResetFeed(ctx)
	if err != nil {
		fmt.Println("could not reset the feeds table")
		os.Exit(1)
		return err
	}
	return nil
}

// prints all the users that are registered and states which is the current user
func handlerUserList(s *state, cmd command) error {
	ctx := context.Background()

	//gets a list of all the names
	names, err := s.db.ListUsers(ctx)
	if err != nil {
		os.Exit(1)
		return err
	}

	//prints all the names and adds (current) to the current user
	currUser := *(s.config.CurrentUserName)
	for _, val := range names {
		if val == currUser {
			val = val + " (current)"
		}
		fmt.Printf("* %v\n", val)
	}
	return nil
}

// gets the feed from the website given the website
func handlerAgg(s *state, cmd command) error {
	//gets the feed from the website/ handles error
	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		os.Exit(1)
		return err
	}

	//Handles any special html characters in the string
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	for _, val := range rssFeed.Channel.Item {
		val.Title = html.UnescapeString(val.Title)
		val.Description = html.UnescapeString(val.Description)
	}

	//prints the rssFeed stuct
	fmt.Printf("%+v\n", rssFeed)
	return nil
}

// Adds to thee feed table
func handlerAddFeed(s *state, cmd command) error {
	//gets the current user so we can later use his ID/ handles any errors
	user, err := s.db.GetUser(context.Background(), *s.config.CurrentUserName)
	if err != nil {
		fmt.Println("error getting user in addfeed")
		os.Exit(1)
		return err
	}

	//Makes a feed struct so we can add it into our feeds table
	feedstr := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), feedstr)
	if err != nil {
		fmt.Println("error with creating a new feed. Duplicate")
		os.Exit(1)
		return err
	}

	fmt.Printf("ID: %v\n", feed.ID)
	fmt.Printf("Created at: %v\n", feed.CreatedAt)
	fmt.Printf("Updated at: %v\n", feed.UpdatedAt)
	fmt.Printf("Name: %v\n", feed.Name)
	fmt.Printf("URL: %v\n", feed.Url)
	fmt.Printf("UserID: %v\n", feed.UserID)
	return nil
}

// Prints out the feed name, url, and user that created the feed
func handlerListFeeds(s *state, cmd command) error {
	ctx := context.Background()

	//gets all of the data from the sql query and handles if there is an issue
	feed, err := s.db.ListFeeds(ctx)
	if err != nil {
		fmt.Println("error with getting the list of names/url/username")
		os.Exit(1)
		return err
	}

	//prints out all of the feed data (name/url/username)
	for _, val := range feed {
		fmt.Printf("%s (%s) by %s\n", val.Name, val.Url, val.Name_2)
	}
	return nil
}
