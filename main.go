package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/cli"
)

var (
	// AppVersion is the version of this application.
	AppVersion string = "0.0.1"

	// AppRevision is the source control revision for this application.
	AppRevision string = ""
)

func main() {
	cliVersion := fmt.Sprintf("%s", AppVersion)
	if AppRevision != "" {
		cliVersion = fmt.Sprintf("%s-%s", AppVersion, AppRevision)
	}

	cli := &cli.CLI{
		Args: os.Args[1:],
		Commands: map[string]cli.CommandFactory{
			"":         defaultCommandFactory,
			"list":     listCommandFactory,
			"download": downloadCommandFactory,
		},
		Name:         "hashicorp-releases",
		Version:      cliVersion,
		Autocomplete: true,
		HelpFunc:     cli.BasicHelpFunc("hashicorp-releases"),
		HelpWriter:   os.Stdout,
		ErrorWriter:  os.Stderr,
	}

	exitStatus, err := cli.Run()
	if err != nil {
		log.Println(err)
	}
	os.Exit(exitStatus)
}

func defaultCommandFactory() (cli.Command, error) {
	var d defaultCommand
	return &d, nil
}

type defaultCommand struct{}

func (d *defaultCommand) Help() string {
	return "Help called for default subcommand"
}

func (d *defaultCommand) Run(args []string) int {
	log.Println("Run called for default subcommand")
	log.Println(args)
	return 0
}

func (d *defaultCommand) Synopsis() string {
	return "Synopsis called for default subcommand"
}

func listCommandFactory() (cli.Command, error) {
	var l listCommand
	return &l, nil
}

type listCommand struct{}

func (l *listCommand) Help() string {
	return "Help called for list subcommand"
}

func (l *listCommand) Run(args []string) int {
	log.Println("Run called for list subcommand")
	return 0
}

func (l *listCommand) Synopsis() string {
	return "Synopsis called for list subcommand"
}
