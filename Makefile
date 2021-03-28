PROJECTNAME=$(shell basename "$(PWD)")

help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: help
all: help

## build: build the parser
build:
	go build -o bin/parser cmd/main.go

## vendor: constructs a directory named vendor in the main module's root directory that contains copies of all packages needed to support builds and tests of packages in the main module
vendor:
	go mod vendor
