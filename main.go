package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"workspace/github.com/christianrm0821/blogAggregator/internal/config"
	"workspace/github.com/christianrm0821/blogAggregator/internal/database"

	_ "github.com/lib/pq"
)

// State struct, takes in a type from database and config
type state struct {
	db     *database.Queries
	config *config.Config
}

// map which maps all of the command names to the commands they do
type commands struct {
	cmdMap map[string]func(*state, command) error
}

// Registers a new command into the map which maps command names to commands
func (c *commands) registerCommand(name string, f func(*state, command) error) {
	(*c).cmdMap[name] = f
}

// Runs the command that is passed by going through the map
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
	//gets arguments from command line
	myArgs := os.Args

	//returns a database pointer or an error
	//sql.Open takes in the database driver's name and where it is going to read from (a string)
	db, err := sql.Open("postgres", "postgres://postgres:CT_rm0821@localhost:5432/gator?sslmode=disable")
	if err != nil {
		fmt.Println("failed to open the ")
		return
	}

	//Returns a pointer to a query struct from the database pointer we got above
	dbQueries := database.New(db)

	//checks if we have actually typed a command
	if len(myArgs) < 2 {
		fmt.Println("you need more arguments")
	}

	//makes all the words typed in lower-case
	for i := range myArgs {
		myArgs[i] = strings.ToLower(myArgs[i])
	}

	//gets the first second element in args which is the command name
	actualCommand := myArgs[1]

	//checks if it is login and if it is then checks if it has a username after login
	if actualCommand == "login" && len(myArgs) < 3 {
		fmt.Println("need a username")
		return
	}

	//checks if it is register and if it is then checks if it has a username
	if actualCommand == "register" && len(myArgs) < 3 {
		fmt.Println("need a user to add ")
		return
	}

	//changes the current slice to make the first argument the username
	myArgs = myArgs[2:]

	//fills the ans config with data and handles if there is an error
	var ans *config.Config
	ans, err = config.Read()
	if err != nil {
		fmt.Println("error with the first read")
		os.Exit(1)
	}

	//Fills in a state struct and uses the new ans config as well
	//as the dbQueries we obtained earlier from database.New()
	myState := state{
		db:     dbQueries,
		config: ans,
	}

	//Make a map and maps command names to commands.
	//registers the commands "login", "register", "reset", "users"
	myCommands := commands{
		cmdMap: make(map[string]func(*state, command) error),
	}
	myCommands.registerCommand("login", handlerLogin)
	myCommands.registerCommand("register", handlerRegister)
	myCommands.registerCommand("reset", handlerReset)
	myCommands.registerCommand("users", handlerUserList)
	myCommands.registerCommand("agg", handlerAgg)

	//makes a command struct and assigns it the arguments as well as
	// the command name
	cmd := command{
		name: &actualCommand,
		args: myArgs,
	}

	//runs the command with the currect state and current command which
	//was given in the terminal. If there was an error then it handles it
	err = myCommands.run(&myState, cmd)
	if err != nil {
		fmt.Println("error with running myCommand")
		os.Exit(1)
	}

	//Prints out the current state of the ans config to check if the
	//user is correct after the command has been ran.
	//ans, _ = config.Read()
	/*
		fmt.Printf("DBurl: %v\n", *((*ans).DbURL))
		fmt.Printf("username: %v\n", *((*ans).CurrentUserName))
	*/
}
