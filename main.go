package main

import (
	"fmt"
	"workspace/github.com/christianrm0821/blogAggregator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	var ans *config.Config
	ans, err := config.Read()
	fmt.Printf("DBurl: %v\n", *((*ans).DbURL))
	fmt.Printf("username: %v\n", *((*ans).CurrentUserName))

	if err != nil {
		return
	}
	user := "Christian"
	err = ans.SetUser(&user)
	if err != nil {
		return
	}
	ans, _ = config.Read()
	fmt.Printf("DBurl: %v\n", *((*ans).DbURL))
	fmt.Printf("username: %v\n", *((*ans).CurrentUserName))

}
