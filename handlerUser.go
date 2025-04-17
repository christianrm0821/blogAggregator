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
		fmt.Println("login expects to have a username")
		os.Exit(1)
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
	//err := fmt.Errorf("")
	if len(cmd.args) == 0 {
		fmt.Print("register expects a username after command")
		os.Exit(1)
		return fmt.Errorf("add a username to command")
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

//Fetches the actual feed
/*
	complete_feed, err := s.db.GetFeedByUrl(context.Background(), feed.Url)
	if err != nil {
		fmt.Println("error getting the actual feed using the URL")
		os.Exit(1)
		return err
	}
*/

// Get all the feeds that need to be fetched and print their titles to the console
func scrapeFeeds(s *state, user database.User) error {
	//gets next feed to fetch(not the whole feed just the ID)
	feed, err := s.db.GetNextFeedToFetch(context.Background(), user.ID)
	if err != nil {
		fmt.Println("error getting feed of current user")
		return err
	}

	//Marks the feed that was fetched at current time
	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		fmt.Println("error marking feed as fetched")
		return err
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)

	for _, val := range rssFeed.Channel.Item {
		val.Title = html.UnescapeString(val.Title)
		fmt.Printf("%v\n", val.Title)
	}
	return nil

}

// gets the feed from the next website(1 argument(time (1s,1m,1h)))
func handlerAgg(s *state, cmd command, user database.User) error {
	//checks if argument was input
	if len(cmd.args) < 1 {
		fmt.Println("need to input a time between request(1s,1m,1h)")
		return fmt.Errorf("need to input a time between request(1s,1m,1h)")
	}

	//make a duration type from the argument entered with command handles if input is not valid
	time_between_reqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		fmt.Println("need to enter a valid time")
		return fmt.Errorf("need to enter a valid time")
	}
	fmt.Printf("Collecting feed every %v\n", time_between_reqs)

	//creates a new ticker type with the time given
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s, user)
	}
}

// Adds to the feed table
func handlerAddFeed(s *state, cmd command, user database.User) error {
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
		fmt.Printf("error with creating a new feed. %v\n", err)
		os.Exit(1)
		return err
	}

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), feedFollow)
	if err != nil {
		fmt.Println("error creating feed follow in handlerAddFeed")
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

// follow command takes a url and makes a new follow record for the current user(1 argument: url of feed to follow)
func handlerFollow(s *state, cmd command, user database.User) error {
	//make sure that an argument was added to the command call
	if len(cmd.args) == 0 {
		fmt.Println("need to include a URL to follow")
		os.Exit(1)
		return fmt.Errorf("need to include URL of feed")
	}
	//gets the url into feedurl
	fmt.Printf("%v\n", cmd.args[0])
	feedurl := cmd.args[0]

	//get the feed by using url
	feed, err := s.db.GetFeedByUrl(context.Background(), feedurl)
	if err != nil {
		fmt.Println("error getting feed using the url")
		os.Exit(1)
		return err
	}

	currUser, err := s.db.GetUser(context.Background(), *s.config.CurrentUserName)
	if err != nil {
		fmt.Println("could not get current username")
		os.Exit(1)
		return err
	}

	feedFollows := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		UserID:    currUser.ID,
		FeedID:    feed.ID,
	}

	createFFR, err := s.db.CreateFeedFollow(context.Background(), feedFollows)
	if err != nil {
		//fmt.Println("error creating the feed follows row")
		fmt.Printf("error: %v ", err)
		os.Exit(1)
		return err
	}

	fmt.Printf("Feed Name: %v\n", createFFR.FeedName)
	fmt.Printf("Current user: %v\n", createFFR.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	user, err := s.db.GetUser(context.Background(), *s.config.CurrentUserName)
	if err != nil {
		fmt.Println("error getting the user")
		os.Exit(1)
		return err
	}

	feedsFollowed, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		fmt.Println("error getting the followed feed from current user")
		os.Exit(1)
		return err
	}
	for i, val := range feedsFollowed {
		fmt.Printf("%v. %v\n", i+1, val.FeedName)
		i++
	}
	return nil

}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		fmt.Println("need to input a url")
		os.Exit(1)
		return fmt.Errorf("no url input")
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Println("error getting feed with url")
		os.Exit(1)
		return err
	}
	unfollow := database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.db.UnfollowFeed(context.Background(), unfollow)
	if err != nil {
		fmt.Println("error unfollowing feed")
		os.Exit(1)
		return err
	}
	return nil
}
