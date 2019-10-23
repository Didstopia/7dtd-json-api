SHELL := /bin/bash

export GO111MODULE=on
export GOPROXY=

export PATH := $(GOPATH)/bin:$(PATH)

BINARY_VERSION?=0.0.1
BINARY_OUTPUT?=7dtd-json-api
EXTRA_FLAGS?=-mod=vendor

define timed_function
	@d=$$(date +%s); \
	$(shell echo $1); \
	echo "=> Ran $1 in $$(($$(date +%s)-d)) seconds"
endef

.PHONY: all install uninstall build test clean deps upgrade tidy

all: deps build

install:
	$(call timed_function,'go install -v $(EXTRA_FLAGS) -ldflags "-X main.Version=$(BINARY_VERSION)"')

uninstall:
	$(call timed_function,'rm -f $(GOPATH)/bin/$(BINARY_OUTPUT)')

build:
	$(call timed_function,'go build -v $(EXTRA_FLAGS) -ldflags "-X main.Version=$(BINARY_VERSION)" -o $(BINARY_OUTPUT)')

test:
	$(call timed_function,'go test -v $(EXTRA_FLAGS) -race -coverprofile=coverage.txt -covermode=atomic ./...')

clean:
	$(call timed_function,'go clean')
	$(call timed_function,'rm -f $(BINARY_OUTPUT)')

deps:
	$(call timed_function,'go build -v $(EXTRA_FLAGS) ./...')

upgrade:
	$(call timed_function,'go get -u ./...')
	$(call timed_function,'go mod vendor')
	$(call timed_function,'go mod tidy')

tidy:
	$(call timed_function,'go mod tidy')
