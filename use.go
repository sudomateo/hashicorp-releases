package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/cli"
	"github.com/sudomateo/hashicorp-releases/pkg/hcrelease"
)

// fileExists checks whether or not a file exists.
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// useCommandFactory is the factory that produces the use command.
func useCommandFactory() (cli.Command, error) {
	var i useCommand
	return &i, nil
}

// useCommand is a blank struct that satisfies the cli.Command interface.
type useCommand struct{}

// Help prints help text for the use command.
func (u *useCommand) Help() string {
	help := `Usage: hashicorp-releases use <product> <version>`
	return help
}

// Run runs the use command.
func (u *useCommand) Run(args []string) int {
	if len(args) < 2 {
		fmt.Println("The use command expects exactly two arguments.")
		fmt.Printf("%s\n", u.Help())
		return 1
	}
	product := args[0]
	version := args[1]

	user, err := user.Current()
	if err != nil {
		log.Printf("failed to retrieve home directory: %v", err)
		return 1
	}
	homeDir := user.HomeDir

	dataDir := filepath.Join(homeDir, ".local/share/hashicorp-releases")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Printf("failed to create data directory: %v", err)
		return 1
	}

	productURL, err := url.Parse(hcrelease.ReleasesURL)
	if err != nil {
		return 1
	}
	productURL.Path = "index.json"

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

	tmpfile, err := ioutil.TempFile(dataDir, "terraform")
	if err != nil {
		log.Printf("failed to create temporary directory: %v", err)
		return 1
	}
	defer os.Remove(tmpfile.Name())

	err = build.Download(tmpfile)
	if err != nil {
		log.Printf("failed to download build: %v", err)
		return 1
	}

	ext := filepath.Ext(build.Filename)
	if ext != ".zip" {
		log.Printf("invalid file extenstion %s: %v", ext, err)
		return 1
	}

	zipReader, err := zip.OpenReader(tmpfile.Name())
	if err != nil {
		log.Printf("failed to open zip reader: %v", err)
		return 1
	}
	defer zipReader.Close()

	fileName := fmt.Sprintf("%s_%s", build.Name, build.Version)
	outPath := filepath.Join(dataDir, fileName)
	for _, f := range zipReader.File {
		if f.Name == build.Name {
			rc, err := f.Open()
			if err != nil {
				log.Printf("failed to open file: %v", err)
				return 1
			}
			defer rc.Close()

			outFile, err := os.Create(outPath)
			if err != nil {
				log.Printf("failed to create file: %v", err)
				return 1
			}
			if err := outFile.Chmod(0755); err != nil {
				log.Printf("failed to chmod file: %v", err)
				return 1
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, rc)
			if err != nil {
				log.Printf("failed to copy file: %v", err)
				return 1
			}
			break
		}
	}

	symlinkDir := filepath.Join(homeDir, ".local/bin")
	if err := os.MkdirAll(symlinkDir, 0755); err != nil {
		log.Printf("failed to create symlink directory: %v", err)
		return 1
	}

	symlinkPath := filepath.Join(symlinkDir, build.Name)
	exists, err := fileExists(symlinkPath)
	if err != nil {
		log.Printf("failed to check for file existence: %v", err)
		return 1
	}
	if exists {
		if err := os.Remove(symlinkPath); err != nil {
			log.Printf("failed to remove file: %v", err)
			return 1
		}
	}

	if err := os.Symlink(outPath, symlinkPath); err != nil {
		log.Printf("failed to create symlink: %v", err)
		return 1
	}

	return 0
}

// Synopsis prints a one-liner about the use command.
func (u *useCommand) Synopsis() string {
	return "Use a specific version of a product."
}
