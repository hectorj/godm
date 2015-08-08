# gpm
A go1.5+ package manager.

[![Build Status](https://travis-ci.org/hectorj/gpm.svg?branch=master)](https://travis-ci.org/hectorj/gpm) [![GoDoc](https://godoc.org/github.com/hectorj/gpm?status.svg)](https://godoc.org/github.com/hectorj/gpm/) [![Coverage Status](https://coveralls.io/repos/hectorj/gpm/badge.svg?branch=master)](https://coveralls.io/r/hectorj/gpm?branch=master)

## Status

Still in very early development.

## Installation

```bash
go get github.com/hectorj/gpm
```

## Usage

```bash
# Go to your package directory, wherever that is
cd $GOPATH/src/myPackage
# run gpm on the go files from which you want to vendor imported packages
gpm *.go
# if everything went well, you have new git submodules you can commit
git commit -m "Vendoring as git submodules done by gpm"
```