export GO111MODULE = on
export GOPROXY = https://proxy.golang.org

NAME = imgbb
BINDIR ?= ./bin
PACKDIR ?= ./build/package
GOBUILD := CGO_ENABLED=0 go build --ldflags="-s -w" -v
GOFILES := $(wildcard ./cmd/ipfs-pinner/*.go)
VERSION := $(shell git describe --tags `git rev-list --tags --max-count=1`)
VERSION := $(VERSION:v%=%)
PROJECT := github.com/wabarc/imgbb
PACKAGES := $(shell go list ./...)

PLATFORM_LIST = \
	darwin-amd64 \
	darwin-arm64 \
	linux-amd64

WINDOWS_ARCH_LIST = \
	windows-amd64

.PHONY: all
all: linux-amd64 darwin-amd64 windows-amd64

darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

linux-amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

windows-amd64:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe $(GOFILES)

fmt:
	@echo "-> Running go fmt"
	@go fmt $(PACKAGES)

tar_releases := $(addsuffix .gz, $(PLATFORM_LIST))
zip_releases := $(addsuffix .zip, $(WINDOWS_ARCH_LIST))

$(tar_releases): %.gz : %
	@mkdir -p $(PACKDIR)
	chmod +x $(BINDIR)/$(NAME)-$(basename $@)
	tar -czf $(PACKDIR)/$(NAME)-$(basename $@)-$(VERSION).tar.gz --transform "s/$(notdir $(BINDIR))//g" $(BINDIR)/$(NAME)-$(basename $@)

$(zip_releases): %.zip : %
	@mkdir -p $(PACKDIR)
	zip -m -j $(PACKDIR)/$(NAME)-$(basename $@)-$(VERSION).zip $(BINDIR)/$(NAME)-$(basename $@).exe

all-arch: $(PLATFORM_LIST) $(WINDOWS_ARCH_LIST)

releases: $(tar_releases) $(zip_releases)

clean:
	rm -f $(BINDIR)/*
	rm -f $(PACKDIR)/*

tag:
	git tag v$(VERSION)
