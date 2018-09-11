package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/magicly/monkey-go/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s - %s! This is the Monkey programming language!\n", user.Username, user.Name)
	repl.Start(os.Stdin, os.Stdout)
}
