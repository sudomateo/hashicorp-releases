package main

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/pengsrc/go-shared/log"
	"github.com/sudomateo/hashicorp-releases/pkg/hcrelease"
)

func listCommandFactory() (cli.Command, error) {
	var l listCommand
	return &l, nil
}

type listCommand struct{}

func (l *listCommand) Help() string {
	return "Help called for list subcommand"
}

func (l *listCommand) Run(args []string) int {
	url, err := url.Parse(hcrelease.ReleasesURL)
	if err != nil {
		log.Errorf(context.Background(), "failed to parse URL %s: %v", hcrelease.ReleasesURL, err)
		return 1
	}
	url.Path = "index.json"
	products, err := hcrelease.GetProducts(url.String())
	if err != nil {
		log.Errorf(context.Background(), "failed to retreive products from %s: %v", hcrelease.ReleasesURL, err)
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

func (l *listCommand) Synopsis() string {
	return "List the available products"
}
