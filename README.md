sole-server
==========

[![Build Status](https://travis-ci.org/solefaucet/sole-server.svg?branch=master)](https://travis-ci.org/solefaucet/solebtc)
[![Go Report Card](http://goreportcard.com/badge/solefaucet/sole-server)](http://goreportcard.com/report/solefaucet/solebtc)
[![codecov.io](https://codecov.io/github/solefaucet/sole-server/coverage.svg?branch=master)](https://codecov.io/github/solefaucet/solebtc?branch=master)

======

## Requirement

* go1.6
* mysql5.7

## Installation

```bash
# easy enough by go get
$ go get -u github.com/solefaucet/sole-server
```

## DB Migration

#### Requirement

goose is needed for DB migration

```bash
$ go get bitbucket.org/liamstask/goose/cmd/goose
```

#### How to

```bash
# Migrate DB to the most recent version available
$ goose up

# Roll back version by 1
$ goose down

# Create a new migration
$ goose create SomeThingDescriptiveEnoughForYourChangeToDB sql
```

## Development

#### Dependency Management

```bash
# After third party library is introduced or removed
$ GO15VENDOREXPERIMENT="0" godep save -r ./...
```

#### Lint

```bash
$ make metalint
```

#### Test

```bash
$ make test
```

#### Benchmark

```bash
$ make benchmark
```

## Deployment

#### Requirement

fabric is needed for deployment.

```bash
$ pip install fabric
```

#### How to

If you have access to my server, simply run

```bash
$ fab -R production deploy:branch_name=master
```

But I am sure you do not XD
