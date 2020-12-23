package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/cli"
	"github.com/sudomateo/hashicorp-releases/pkg/hcrelease"
)

// installCommandFactory is the factory that produces the install command.
func installCommandFactory() (cli.Command, error) {
	var i installCommand
	return &i, nil
}

// installCommand is a blank struct that satisfies the cli.Command interface.
type installCommand struct{}

// Help prints help text for the install command.
func (i *installCommand) Help() string {
	help := `Usage: hashicorp-releases install <product> <version>`
	return help
}

// Run runs the install command.
func (i *installCommand) Run(args []string) int {
	if len(args) < 2 {
		fmt.Println("The install command expects exactly two arguments.")
		fmt.Printf("%s\n", i.Help())
		return 1
	}
	product := args[0]
	version := args[1]

	user, err := user.Current()
	if err != nil {
		fmt.Printf("failed to retrieve home directory: %v", err)
		return 1
	}
	homeDir := user.HomeDir

	dataDir := filepath.Join(homeDir, ".local/share/hashicorp-releases")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Printf("failed to create data directory: %v", err)
		return 1
	}

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

	tmpfile, err := ioutil.TempFile(dataDir, "terraform")
	if err != nil {
		fmt.Printf("failed to create temporary directory: %v", err)
		return 1
	}
	defer os.Remove(tmpfile.Name())

	err = build.Download(tmpfile)
	if err != nil {
		fmt.Printf("failed to download build: %v", err)
		return 1
	}

	ext := filepath.Ext(build.Filename)
	if ext != ".zip" {
		fmt.Printf("invalid file extenstion %s: %v", ext, err)
		return 1
	}

	zipReader, err := zip.OpenReader(tmpfile.Name())
	if err != nil {
		fmt.Printf("failed to open zip reader: %v", err)
		return 1
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		if f.Name == build.Name {
			rc, err := f.Open()
			if err != nil {
				fmt.Printf("failed to open file: %v", err)
				return 1
			}
			defer rc.Close()

			fileName := fmt.Sprintf("%s_%s", build.Name, build.Version)
			outPath := filepath.Join(dataDir, fileName)
			outFile, err := os.Create(outPath)
			if err != nil {
				fmt.Printf("failed to create file: %v", err)
				return 1
			}
			if err := outFile.Chmod(0755); err != nil {
				fmt.Printf("failed to chmod file: %v", err)
				return 1
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, rc)
			if err != nil {
				fmt.Printf("failed to copy file: %v", err)
				return 1
			}
			break
		}
	}

	return 0
}

// Synopsis prints a one-liner about the install command.
func (i *installCommand) Synopsis() string {
	return "Install a specific version of a product."
}
