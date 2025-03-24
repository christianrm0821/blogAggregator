package main

import (
	"fmt"
	"workspace/github.com/christianrm0821/blogAggregator/internal/config"
)

func main() {
	var ans *config.Config
	ans, err := config.Read()
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
