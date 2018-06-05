.DEFAULT_GOAL := help

VERSION := $(shell git describe --tags --abbrev=0)
VERSION_LONG := $(shell git describe --tags)
VAR_VERSION := github.com/tmtk75/aws-switch-role/cmd.Version

LDFLAGS := -ldflags "-X $(VAR_VERSION)=$(VERSION) \
	-X $(VAR_VERSION)Long=$(VERSION_LONG)"

SRCS := $(shell find . -type f -name '*.go')

.PHONY: build
build: aws-switch-role  ## Build

aws-switch-role: $(SRCS)
	go build $(LDFLAGS) -o aws-switch-role

.PHONY: install
install:  ## Install in GOPATH
	go install $(LDFLAGS) .

.PHONY: clean
clean:  ## Clean
	rm -f aws-switch-role


## Release targets
.PHONY: build-release archive
build-release: build/aws-switch-role_linux_amd64
archive: build/aws-switch-role_linux_amd64.gz
release: upload-archives

upload-archives: build/aws-switch-role_linux_amd64.zip
	ghr -u tmtk75 $(VERSION) ./build/*.zip

build/aws-switch-role_linux_amd64.zip: build-release
	(cd build; zip aws-switch-role_linux_amd64.zip aws-switch-role_linux_amd64)

build/aws-switch-role_linux_amd64: generate
	GOARCH=amd64 GOOS=linux go build -o build/aws-switch-role_linux_amd64 ./cmd/aws-switch-role/main.go

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'
