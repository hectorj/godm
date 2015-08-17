# godm
A go1.5+ dependencies manager. [![Build Status](https://travis-ci.org/hectorj/godm.svg?branch=master)](https://travis-ci.org/hectorj/godm) [![GoDoc](https://godoc.org/github.com/hectorj/godm?status.svg)](https://godoc.org/github.com/hectorj/godm/) [![Coverage Status](https://coveralls.io/repos/hectorj/godm/badge.svg?branch=master)](https://coveralls.io/r/hectorj/godm?branch=master)

More precisely, a tool to manage your project's dependencies by vendoring them at pinpointed versions.

It relies on the "GO15VENDOREXPERIMENT", so that other people (users and developers) can simply `go get` your project
 without being forced to use `godm` or any other tool that doesn't come with Go right out of the box.
 
If you wish to see how does a project using `godm` looks, well you got one right here :)

Vendors are [there](vendor) as [Git submodules](.gitmodules)

## Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Status](#status)
- [Dependencies](#dependencies)
- [Installation](#installation)
- [Upgrade](#upgrade)
- [Usage](#usage)
  - [help](#help)
  - [vendor](#vendor)
  - [clean](#clean)
  - [remove](#remove)
- [Bash Autocompletion](#bash-autocompletion)
- [Migrating from another dependenices management tool](#migrating-from-another-dependenices-management-tool)
  - [Godep](#godep)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Status

Still in very early development.

## Dependencies

Go 1.5+ (with `GO15VENDOREXPERIMENT=1`), and Git (available in your $PATH)

## Installation

```bash
# If you haven't already, enable the Go 1.5 vendor experiment (personally that line is in my ~/.bashrc).
export GO15VENDOREXPERIMENT=1
# Then it's a simple go get.
go get github.com/hectorj/godm/cmd/godm
```

## Upgrade

```bash
go get -u github.com/hectorj/godm/cmd/godm
```

## Usage

Note : does not support Mercurial yet

### help

Auto-generated help is available like this :

```bash
godm --help
```

(thanks to https://github.com/codegangsta/cli)

### vendor

The `vendor` sub-command vendors imports that are not vendored yet in the current project. Outputs the import paths of vendors successfully added.

```bash
# Go to your package directory, wherever that is.
cd $GOPATH/src/myPackage
# Run it.
godm vendor
# If everything went well, you have new git submodules you can commit.
git commit -m "Vendoring done by godm"
```

### clean

The `clean` sub-command removes vendors that are not imported in the current project. Outputs the import paths of vendors successfully removed.

```bash
# Go to your package directory, wherever that is.
cd $GOPATH/src/myPackage
# Run it.
godm clean
# If everything went well, you may have some Git submodules removals you can commit.
git commit -m "Vendors cleaning done by godm"
```

### remove

The `remove` sub-command unvendors an import path. Takes a single import path as argument.
```bash
# Go to your package directory, wherever that is.
cd $GOPATH/src/myPackage
# Run godm
godm remove github.com/my/import/path
# If everything went well, you have a Git submodule removal you can commit.
git commit -m "Unvendoring done by godm"
```

## Bash Autocompletion

Copy [godm_bash_autocomplete.bash](godm_bash_autocomplete.bash) to `/etc/bash_completion.d/` to get command autocompletion (highly recommended)

## Migrating from another dependenices management tool

### Godep

```bash
godep restore
GOPATH=`godep path`:$GOPATH godm save
```

Once you have checked the migration went well, you can eventually `rm -rf ./Godeps`
