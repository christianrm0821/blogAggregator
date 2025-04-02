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
		fmt.Print("error getting request\n")
		os.Exit(1)
		return nil, err
	}

	//get a response and handle any errors
	res, err := client.Do(req)
	if err != nil {
		fmt.Print("error getting response\n")
		os.Exit(1)
		return nil, err
	}

	//set "User-Agent" header to gator
	//This is to identify the program to the server
	res.Header.Set("User-Agent", "gator")

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
		os.Exit(1)
		return nil, err
	}
	return &rssFeed, nil
}
