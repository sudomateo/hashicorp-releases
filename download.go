package main

import (
	"log"
	"net/url"
	"os"
	"path"
	"runtime"

	"github.com/mitchellh/cli"
	"github.com/sudomateo/hashicorp-releases/pkg/hcrelease"
)

func downloadCommandFactory() (cli.Command, error) {
	var l downloadCommand
	return &l, nil
}

type downloadCommand struct{}

func (l *downloadCommand) Help() string {
	return "Help called for download subcommand"
}

func (l *downloadCommand) Run(args []string) int {
	if len(args) < 2 {
		log.Print("must provide at least 2 arguments")
		return 1
	}
	product := args[0]
	version := args[1]
	productURL, err := url.Parse(hcrelease.ReleasesURL)
	if err != nil {
		return 1
	}
	productURL.Path = path.Join("index.json")
	products, err := hcrelease.GetProducts(productURL.String())
	if err != nil {
		log.Printf("failed to retrieve product details: %v", err)
		return 1
	}

	release, err := products.GetRelease(product)
	if err != nil {
		log.Printf("failed to retrieve release details: %v", err)
		return 1
	}

	ver, err := release.GetVersion(version)
	if err != nil {
		log.Printf("failed to retrieve version details: %v", err)
		return 1
	}

	build, err := ver.GetBuild(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		log.Printf("failed to retrieve build details: %v", err)
		return 1
	}

	outFile, err := os.Create(build.Filename)
	if err != nil {
		return 1
	}
	defer outFile.Close()

	err = build.Download(outFile)
	if err != nil {
		log.Printf("failed to download build: %v", err)
		return 1
	}

	return 0
}

func (l *downloadCommand) Synopsis() string {
	return "Synopsis called for download subcommand"
}
