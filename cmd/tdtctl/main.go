// package main is the entrypoint to the tdtctl CLI app.
//
// tdtctl manages the site.
package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/eculver/tdtv2/pkg/command"
	"github.com/eculver/tdtv2/pkg/version"
	"github.com/mitchellh/cli"
)

func main() {
	// not sure this makes sense but this tool will be used by
	// the site admin for doing things like:
	// - scaffolding a new episode
	// - running the episode publishing workflow
	// - managing/scaling the site?
	// UPDATE 2023-01-13: SoT for media/episodes is Anchor.fm, this binary can
	// serve the site by checking for episodes periodically in the Anchor.fm API

	// examples:
	// (scaffold out files for new episode)
	// tdtctl generate episode/211
	// (sync new episode to site -- this would be run post-merge to deploy changes)
	// tdtctl apply -f ./index/episodes/211.yml

	fmt.Println("I will eventually be the site admin's swiss army knife\n\n")
	fmt.Println("tdtctl create tdt/211")
	fmt.Println("tdtctl apply -f ./index/tdts")
	fmt.Println("tdtctl get feed -o itunes")
	fmt.Println("tdtctl serve -index ./index")
	fmt.Println("tdtctl gen -index ./index -out ./path/to/wwwdir")
	fmt.Println("tdtctl sync -index ./index https://anchor.fm/s/b781db40/podcast/rss")

	if err := inner(); err != nil {
		var cliErr *CLIError
		if errors.As(err, &cliErr) {
			os.Exit(cliErr.ExitCode)
		}
		log.Printf("tdtctl error: %s\n", err)
		os.Exit(1)
	}
}

// CLIError contains information about errors from the CLI
type CLIError struct {
	ExitCode int
}

// Error implements the error interface
func (ce *CLIError) Error() string {
	return fmt.Sprintf("exit code %d", ce.ExitCode)
}

func inner() error {
	c := cli.NewCLI("tdtctl", version.Version)
	c.HelpWriter = os.Stdout
	c.ErrorWriter = os.Stderr
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"create": command.CreateEntryFactory,
		"serve":  command.ServeFactory,
	}
	exitStatus, err := c.Run()
	if err != nil {
		return err
	}
	if exitStatus != 0 {
		return &CLIError{ExitCode: exitStatus}
	}
	return nil
}
