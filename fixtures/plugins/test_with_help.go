/**
	* 1. Setup the server so cf can call it under main.
				e.g. `cf my-plugin` creates the callable server. now we can call the Run command
	* 2. Implement Run that is the actual code of the plugin!
	* 3. Return an error
**/

package main

import (
	"fmt"

	"github.com/cloudfoundry/cli/plugin"
)

type TestWithHelp struct {
	stringForTestWithHelp string
}

func (c *TestWithHelp) Run(args []string, reply *bool) error {
	if args[0] == "help" {
		theHelpCmd()
	}
	return nil
}

func (c *TestWithHelp) GetCommands() []plugin.Command {
	return []plugin.Command{
		{
			Name:     "help",
			HelpText: "help text for test_with_help",
		},
	}
}

func theHelpCmd() {
	fmt.Println("You called help in test_with_help")
}

func main() {
	plugin.Start(new(TestWithHelp))
}
