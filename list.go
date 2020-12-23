package main

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/sudomateo/hashicorp-releases/pkg/hcrelease"
)

// listCommandFactory is the factory that produces the list command.
func listCommandFactory() (cli.Command, error) {
	var l listCommand
	return &l, nil
}

// listCommand is a blank struct that satisfies the cli.Command interface.
type listCommand struct{}

// Help prints the list help text
func (l *listCommand) Help() string {
	help := `Usage: hashicorp-releases list`
	return help
}

// Run runs the list command.
func (l *listCommand) Run(args []string) int {
	url, err := url.Parse(hcrelease.ReleasesURL)
	if err != nil {
		fmt.Printf("failed to parse URL %s: %v", hcrelease.ReleasesURL, err)
		return 1
	}
	url.Path = "index.json"
	products, err := hcrelease.GetProducts(url.String())
	if err != nil {
		fmt.Printf("failed to retreive products from %s: %v", hcrelease.ReleasesURL, err)
		return 1
	}
	allProducts := make([]string, 0, len(products))
	for product := range products {
		allProducts = append(allProducts, product)
	}
	sort.Strings(allProducts)
	fmt.Println(strings.Join(allProducts, "\n"))
	return 0
}

// Synopsis prints a one-liner about the list command.
func (l *listCommand) Synopsis() string {
	return "List the available products."
}
