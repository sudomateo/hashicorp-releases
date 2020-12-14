package hcrelease

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const (
	// ReleasesURL is the the public HashiCorp releases URL.
	ReleasesURL = "https://releases.hashicorp.com"
)

// Products maps a HashiCorp product to a release.
type Products map[string]*Release

// Release represents a release of a HashiCorp product.
type Release struct {
	Name     string              `json:"name"`
	Versions map[string]*Version `json:"versions"`
}

// Version represents a version of a Release.
type Version struct {
	Name             string   `json:"name"`
	Version          string   `json:"version"`
	Shasums          string   `json:"shasums"`
	ShasumsSignature string   `json:"shasums_signature"`
	Builds           []*Build `json:"builds"`
}

// Build represents a build of a Version.
type Build struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
}

// GetProducts gathers details for all HashiCorp products.
func GetProducts(url string) (Products, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	p := make(Products)
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// GetRelease gathers the release details for a given product.
func (p Products) GetRelease(product string) (*Release, error) {
	r, ok := p[product]
	if !ok {
		return nil, errors.New("product not found")
	}
	return r, nil
}

// GetVersion gathers the version details for a given Release.
func (r *Release) GetVersion(version string) (*Version, error) {
	v, ok := r.Versions[version]
	if !ok {
		return nil, errors.New("version not found")
	}
	return v, nil
}

// GetBuild gathers the build details for a given Version.
func (v *Version) GetBuild(os, arch string) (*Build, error) {
	var b *Build
	for _, build := range v.Builds {
		if build.Arch == arch && build.OS == os {
			b = build
		}
	}
	if b == nil {
		return nil, errors.New("build not found")
	}
	return b, nil
}

// Download downloads a specific build of a HashiCorp product, writing the
// downloaded file to w.
func (b *Build) Download(w io.Writer) error {
	resp, err := http.Get(b.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
