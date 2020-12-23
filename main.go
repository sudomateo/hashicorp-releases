package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/cli"
)

var (
	// AppVersion is the version of this application.
	AppVersion string = "0.2.0"

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
			"list":     listCommandFactory,
			"download": downloadCommandFactory,
			"install":  installCommandFactory,
			"use":      useCommandFactory,
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
