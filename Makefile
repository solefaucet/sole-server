#! /usr/bin/make
#
# Makefile for sole-server
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
			 bitbucket.org/liamstask/goose/cmd/goose \
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
		--linter='dupl:dupl -plumbing -threshold {duplthreshold} ./*.go | grep -v "_test.go"::(?P<path>[^\s][^:]+?\.go):(?P<line>\d+)-\d+:\s*(?P<message>.*)' \
		./...

test:
	# run test
	go test -cover ./...
	# cleanup
	mysql -uroot -e 'drop database if exists sole_test;'

benchmark:
	# benchmark
	go test -v -bench=. -benchmem -timeout 30m ./...
