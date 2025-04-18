package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
)

// gets feed from the given URL. Returns a RSSfeed struct
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	//make the client
	client := &http.Client{}

	//make the request and handle any errors
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		fmt.Println("error getting request")
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
		return nil, err
	}
	//set "User-Agent" header to gator
	//This is to identify the program to the server
	req.Header.Set("User-Agent", "gator")

	//get a response and handle any errors
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error getting response")
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
		return nil, err
	}
	defer res.Body.Close()

	//retrieves the data in the form of byte
	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Print("error getting the data")
		os.Exit(1)
		return nil, err
	}

	//Unmarshalling data into rssfeed stuct
	var rssFeed RSSFeed
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		fmt.Println("error unmarshalling data")
		fmt.Printf("Error: %v", err)
		os.Exit(1)
		return nil, err
	}
	return &rssFeed, nil
}
