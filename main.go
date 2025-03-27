package main

import (
	"fmt"
	"os"
	"strings"
	"workspace/github.com/christianrm0821/blogAggregator/internal/config"

	_ "github.com/lib/pq"
)

type state struct {
	config *config.Config
}

type commands struct {
	cmdMap map[string]func(*state, command) error
}

func (c *commands) registerCommand(name string, f func(*state, command) error) {
	(*c).cmdMap[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	val, ok := c.cmdMap[*cmd.name]
	if ok {
		err := val(s, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	myArgs := os.Args
	if len(myArgs) < 2 {
		fmt.Println("you need more arguments")
	}

	for i, _ := range myArgs {
		myArgs[i] = strings.ToLower(myArgs[i])
	}

	actualCommand := myArgs[1]
	if actualCommand == "login" && len(myArgs) < 3 {
		fmt.Println("need a username")
		//return
	}
	myArgs = myArgs[2:]

	var ans *config.Config
	ans, err := config.Read()
	if err != nil {
		fmt.Println("error with the first read")
		os.Exit(1)
	}

	var myState state
	myState.config = ans

	myCommands := commands{
		cmdMap: make(map[string]func(*state, command) error),
	}
	myCommands.registerCommand("login", handlerLogin)

	login := "login"
	cmd := command{
		name: &login,
		args: myArgs,
	}

	err = myCommands.run(&myState, cmd)
	if err != nil {
		fmt.Println("error with running myCommand")
		os.Exit(1)
	}
	ans, _ = config.Read()
	fmt.Printf("DBurl: %v\n", *((*ans).DbURL))
	fmt.Printf("username: %v\n", *((*ans).CurrentUserName))

}
