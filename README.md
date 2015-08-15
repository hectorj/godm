# gpm
A go1.5+ package manager. [![Build Status](https://travis-ci.org/hectorj/gpm.svg?branch=master)](https://travis-ci.org/hectorj/gpm) [![GoDoc](https://godoc.org/github.com/hectorj/gpm?status.svg)](https://godoc.org/github.com/hectorj/gpm/) [![Coverage Status](https://coveralls.io/repos/hectorj/gpm/badge.svg?branch=master)](https://coveralls.io/r/hectorj/gpm?branch=master)

More precisely, a tool to manage your project's dependencies by vendoring them at pinpointed versions.

It relies on the "GO15VENDOREXPERIMENT", so that other people (users and developers) can simply `go get` your project
 without being forced to use `gpm` or any other tool that doesn't come with Go right out of the box.
 
If you wish to see how does a project using `gpm` looks, well you got one right here :)

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
  - [remove](#remove)
- [Bash Autocompletion](#bash-autocompletion)
- [Migrating from another dependenices management tool](#migrating-from-another-dependenices-management-tool)
  - [Godep](#godep)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Status

Still in very early development.

## Dependencies

Go 1.5+ (with `GO15VENDOREXPERIMENT=1`, else even if it runs it's kind of pointless), and Git (available in your $PATH)

## Installation

```bash
# If you haven't already, enable the Go 1.5 vendor experiment (personally that line is in my ~/.bashrc).
export GO15VENDOREXPERIMENT=1
# Then it's a simple go get.
go get github.com/hectorj/gpm/cmd/gpm
```

## Upgrade

```bash
go get -u github.com/hectorj/gpm/cmd/gpm
```

## Usage

Note : does not support Mercurial yet

### help

Auto-generated help is available like this :

```bash
gpm --help
```

(thanks to https://github.com/codegangsta/cli)

### vendor

The `vendor` sub-command takes the current project you're in, extract imports from it, and vendor these imports if possible and necessary.

```bash
# Go to your package directory, wherever that is.
cd $GOPATH/src/myPackage
# Run it.
gpm vendor
# If everything went well, you have new git submodules you can commit.
git commit -m "Vendoring as git submodules done by gpm"
```

### remove

The `remove` sub-command unvendors an import path.
```bash
# Go to your package directory, wherever that is.
cd $GOPATH/src/myPackage
# Run gpm
gpm remove github.com/my/import/path
# If everything went well, you have a git submodule removal you can commit.
git commit -m "Unvendoring done by gpm"
```

## Bash Autocompletion

Copy [gpm_bash_autocomplete.bash](gpm_bash_autocomplete.bash) to `/etc/bash_completion.d/` to get command autocompletion (highly recommended)

## Migrating from another dependenices management tool

### Godep

```bash
godep restore
GOPATH=`godep path`:$GOPATH gpm save
```

Once you have checked the migration went well, you can eventually `rm -rf ./Godeps`
