# hashicorp-releases

The `hashicorp-releases` program is a Go CLI tool that allows you to install HashiCorp tools directly from
[releases.hashicorp.com](https://releases.hashicorp.com). Once installed, a desired version of the tool can be symlinked
for use.

## Installation

`hashicorp-releases` currently supports Linux and macOS.

Linux and macOS installation:

```
git clone git@github.com:sudomateo/hashicorp-releases.git
cd hashicorp-releases
make
```

## Usage

Full CLI output:

```
Usage: hashicorp-releases [--version] [--help] <command> [<args>]

Available commands are:
    download    Download a specific version of a product.
    install     Install a specific version of a product.
    list        List the available products.
    use         Use a specific version of a product.
```

List available products.

```
hashicorp-releases list
```

Download a specific version of a product. This will download the products ZIP file to the current working directory.

```
hashicorp-releases download terraform 0.14.3
```

Install a specific version of a product. This will download and extract a product to
`${HOME}/.local/share/hashicorp-releases`.

```
hashicorp-releases install terraform 0.14.3
```

Use a specific version of a product. This installs a product and creates a symlink named `${HOME}/.local/bin/${PRODUCT}`
whose target is the desired product and version.

```
hashicorp-releases use terraform 0.14.3
```
