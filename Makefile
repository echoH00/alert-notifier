# Makefile for building alert-notifier

# Go parameters
GOOS=linux
GOARCH=amd64
CGO_ENABLED=0
OUTPUT=alert-notifier
CMD_DIR=./cmd

.PHONY: all build clean

all: build

build:
	@echo "Building $(OUTPUT)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build -o $(OUTPUT) $(CMD_DIR)

clean:
	@echo "Cleaning up..."
	@rm -f $(OUTPUT)
