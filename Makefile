VERSION := 0.15.0

all: build

.PHONY: build
build:
	go build

.PHONY: format
format:
	find . -name '*.go'