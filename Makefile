#! /usr/bin/make
#
# Makefile for solebtc
#
# Targets:
# - "depend" retrieves the Go packages needed to run the linter and tests
# - "metalint" runs the linter and checks the code format using goimports
# - "test" runs the tests
#
# Meta targets:
# - "all" is the default target, it runs all the targets in the order above.
#
DEPEND=\
			 golang.org/x/tools/cmd/cover \
			 github.com/smartystreets/goconvey/convey \
			 github.com/smartystreets/assertions \
			 github.com/alecthomas/gometalinter

all: depend metalint test

depend:
	@go get -v $(DEPEND)

metalint:
	gometalinter \
		--disable=gotype \
		--disable=errcheck \
		--disable=deadcode \
		--enable=goimports \
		--deadline=60s \
		./...

test:
	go test -cover ./...
