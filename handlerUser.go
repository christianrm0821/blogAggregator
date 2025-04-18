package main

import (
	"context"
	"fmt"
	"html"
	"os"
	"strconv"
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

// allows us to interpret the time
func time_parse(time_input string) time.Time {
	var layouts []string
	layouts = append(layouts, "2006-01-02 15:04:05")
	layouts = append(layouts, "02 Jan 06 15:04 MST")
	layouts = append(layouts, "2006-01-02T15:04:05Z")
	layouts = append(layouts, "2006-01-02 15:04:05 MST")
	layouts = append(layouts, time.RFC1123)
	layouts = append(layouts, time.RFC822)
	layouts = append(layouts, "02 Jan 2006 15:04:05 -0700")
	var parsed_time time.Time
	var err error

	for _, val := range layouts {
		parsed_time, err = time.Parse(val, time_input)
		if err == nil {
			return parsed_time
		}
	}
	return time.Time{}
}

// Get all the feeds that need to be fetched and print their titles to the console
func scrapeFeeds(s *state, user database.User) error {
	//gets next feedID to fetch(not the whole feed just the ID)
	feed, err := s.db.GetNextFeedToFetch(context.Background(), user.ID)
	if err != nil {
		fmt.Println("error getting feed of current user")
		fmt.Printf("Error: %v\n", err)
		return err
	}

	//Marks the feed that was fetched at current time
	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		fmt.Println("error marking feed as fetched")
		return err
	}

	feedurl, err := s.db.GetFeedURLFromFeedID(context.Background(), feed.FeedID)
	if err != nil {
		fmt.Println("error getting the feed url from ID")
		fmt.Printf("Error: %v\n", err)
	}

	//fetches the feed and places it in rssFeed
	rssFeed, err := fetchFeed(context.Background(), feedurl)
	if err != nil {
		return err
	}
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)

	//Adds the all the posts in the rssfeed into the post table
	for _, val := range rssFeed.Channel.Item {
		post := database.CreatPostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       html.UnescapeString(val.Title),
			Url:         html.UnescapeString(val.Link),
			Description: html.UnescapeString(val.Description),
			PublishedAt: time_parse(val.PubDate),
			FeedID:      feed.FeedID,
		}
		_, err := s.db.CreatPost(context.Background(), post)
		if err != nil {
			fmt.Println("error creating the post")
			fmt.Printf("Error: %v\n", err)
			return err
		}
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
		fmt.Printf("Error: %v", err)
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
		fmt.Printf("Error: %v", err)
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

// Allows users to browse through their posts(takes in 1 optional argument(number of posts to show))
func handlerBrowse(s *state, cmd command, user database.User) error {
	num_posts := 2
	var err error
	if len(cmd.args) > 0 {
		num_posts, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			num_posts = 2
		}
	}
	if num_posts > 10 {
		num_posts = 10
		fmt.Println("The limit it 10")
	}

	posts, err := s.db.GetPostForUser(context.Background(), user.ID)
	if err != nil {
		fmt.Println("make sure to run agg command before running browse")
		fmt.Printf("Error: %v\n", err)
		return err
	}
	for i := 0; i < num_posts; i++ {
		feed_name, err := s.db.GetFeedNameFromID(context.Background(), posts[i].FeedID)
		if err != nil {
			fmt.Println("could not get the feed name")
		}
		fmt.Printf("Feed Name: %v\n", feed_name)
		fmt.Println()
		fmt.Printf("Title: %v", posts[i].Title)
		fmt.Println()
		fmt.Printf("Published: %v", posts[i].PublishedAt)
		fmt.Println()
		fmt.Printf("URL: %v", posts[i].Url)
		fmt.Println()
		fmt.Println("-------------------------------------------------")
		fmt.Println()
	}
	return nil

}
