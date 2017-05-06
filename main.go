package main

import (
	"ender/efp"
	"fmt"
	"os"
	"strings"
)

type efpCommand struct {
	name    string
	execute func([]string)
}

func executeEFP(files []string) {

}

func getCommands() []efpCommand {
	return []efpCommand{
		efpCommand{"help", cmdHelp},
		efpCommand{"prototype", cmdPrototype},
	}
}

func cmdPrototype() {
	fmt.Printf("PROTOTYPE\n")
}

func cmdHelp() {
	fmt.Printf("I AM HELPING\n")
}

func processCommand(prototype *element) {
	// get the args
	args := strings.split(" ")
	for _, cmd := range getCommands() {
		if cmd.name == strings.lower(args[0]) {
			processCommand(p)
			return
		}
	}
	executeEFP(args)
	processCommand(p)
}

func main() {
	args := os.Args
	if len(args) < 1 {
		fmt.Printf("Please specify a .efp prototype file\n")
		return
	}
	p := efp.Prototype(args[1])
	processCommand(p)
}
