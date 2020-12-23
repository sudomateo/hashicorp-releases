package main

import (
	"fmt"
	"net/url"
	"os"
	"runtime"

	"github.com/mitchellh/cli"
	"github.com/sudomateo/hashicorp-releases/pkg/hcrelease"
)

// downloadCommandFactory is a factory that produces the download command.
func downloadCommandFactory() (cli.Command, error) {
	var l downloadCommand
	return &l, nil
}

// downloadCommand is a blank struct that satisfies the cli.Command interface.
type downloadCommand struct{}

// Help prints help text for the download command.
func (d *downloadCommand) Help() string {
	help := `Usage: hashicorp-releases download <product> <version>`
	return help
}

// Run runs the download command.
func (d *downloadCommand) Run(args []string) int {
	if len(args) < 2 {
		fmt.Println("The download command expects exactly two arguments.")
		fmt.Printf("%s\n", d.Help())
		return 1
	}
	product := args[0]
	version := args[1]

	productURL, err := url.Parse(hcrelease.ReleasesURL)
	if err != nil {
		return 1
	}
	productURL.Path = "index.json"

	products, err := hcrelease.GetProducts(productURL.String())
	if err != nil {
		fmt.Printf("failed to retrieve product details: %v", err)
		return 1
	}

	release, err := products.GetRelease(product)
	if err != nil {
		fmt.Printf("failed to retrieve release details: %v", err)
		return 1
	}

	ver, err := release.GetVersion(version)
	if err != nil {
		fmt.Printf("failed to retrieve version details: %v", err)
		return 1
	}

	build, err := ver.GetBuild(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		fmt.Printf("failed to retrieve build details: %v", err)
		return 1
	}

	outFile, err := os.Create(build.Filename)
	if err != nil {
		return 1
	}
	defer outFile.Close()

	err = build.Download(outFile)
	if err != nil {
		fmt.Printf("failed to download build: %v", err)
		return 1
	}

	return 0
}

// Synopsis prints a one-liner about the download command.
func (d *downloadCommand) Synopsis() string {
	return "Download a specific version of a product."
}
