package command

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jaffee/commandeer"
	"github.com/mitchellh/cli"
)

// CreateEntryFactory initializes an instance of the create.entry command
// and configures its defaults.
func CreateEntryFactory() (cli.Command, error) {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return &CreateEntry{
		UI:      ui,
		Options: &CreateEntryOptions{},
	}, nil
}

// CreateEntryOptions holds options for the create.entry command.
type CreateEntryOptions struct {
	DryRun bool `help:"Validates input and then prints information about the entry to stdout; does not write any files to disk"`
	Debug  bool `help:"Print debug info."`

	ObjectPath string `flag:"-"` // first positional argument
}

// CreateEntry encapsulates the create.entry command.
type CreateEntry struct {
	UI      cli.Ui
	Options *CreateEntryOptions
}

// Synopsis summarizes the create.entry command functionality.
func (c *CreateEntry) Synopsis() string {
	return "Create a new site entry with the given path (eg. episode/101)"
}

// Help returns usage information about the fix subcommand
func (c *CreateEntry) Help() string {
	var out strings.Builder
	fs := flag.NewFlagSet("create.entry", flag.ContinueOnError)
	fs.SetOutput(&out)
	if err := commandeer.Flags(fs, c.Options); err != nil {
		panic(err)
	}
	fs.PrintDefaults()
	return fmt.Sprintf(`Usage: tdtctl create [options] <path>

Create a new site entry (eg. episode/101, page/about). The entry's should reflect
a hierarchichal structure for the its associated object. For example, episode objects
are naturally arranged by season (optional) and episode number, so its path should
include those values in its path. A path value of 'episode/101' would mean
'episode number 101'.

Options:
%s
Arguments:
  path
	Path of object.
`, out.String())
}

// Run kicks off the fix command
func (c *CreateEntry) Run(args []string) int {
	if err := c.ParseFlags(args); err != nil {
		return c.exitError(err)
	}

	// create index file

	return 0
}

// ParseFlags populates the command's options by parsing command-line flags.
// A descriptive error is returned if the flags cannot be parsed.
func (c *CreateEntry) ParseFlags(args []string) error {
	fs := flag.NewFlagSet("create.entry", flag.ContinueOnError)
	if err := commandeer.Flags(fs, c.Options); err != nil {
		return err
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	p := fs.Arg(0)
	if p == "" {
		return errors.New("path is required")
	}

	fmt.Printf("got entry path %s\n", p)

	return nil
}

// exitError is a helper function that makes Run's error handling a little
// more Go-like.
func (c *CreateEntry) exitError(err error) int {
	c.UI.Error(fmt.Sprintf("tdtctl create error: %s", err.Error()))
	return 1
}

// debug is a helper for printing debug info when the debug flag is set
func (c *CreateEntry) debug(format string, a ...interface{}) {
	if c.Options.Debug {
		c.UI.Info(fmt.Sprintf(format, a...))
	}
}
