# gpm
A go1.5+ package manager.

[![Build Status](https://travis-ci.org/hectorj/gpm.svg?branch=master)](https://travis-ci.org/hectorj/gpm) [![GoDoc](https://godoc.org/github.com/hectorj/gpm?status.svg)](https://godoc.org/github.com/hectorj/gpm/) [![Coverage Status](https://coveralls.io/repos/hectorj/gpm/badge.svg?branch=master)](https://coveralls.io/r/hectorj/gpm?branch=master)

## Status

Still in very early development.

## Installation

```
go get github.com/hectorj/gpm
```

## Usage

```
cd $GOPATH/src/myPackage
gpm *.go
git commit -m "Vendoring as git submodules done by gpm"
```