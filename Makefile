APPNAME ?= "hashicorp-releases"
GOBIN ?= $(HOME)/go/bin
REVISION ?= $(shell git rev-parse --short HEAD)

$(GOBIN)/$(APPNAME): *.go
	go install -ldflags "-X main.AppRevision=$(REVISION)" ./...

.PHONY: clean
clean:
	rm $(GOBIN)/$(APPNAME)
