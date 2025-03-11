# cqlshvm

A simple CLI tool to download specific versions of the cqlsh utility. The tool
present a list of available versions and then download the one the user chooses.

## Requirements:

Golang 1.23+

## Usage:

Build: `go build`

List versions: `cqlshvm list [-lt <version>] [-gt <version>]` or `go run ./main.go list [-lt <version>] [-gt <version>]`

Download version: `cqlshvm download <version> > cqlsh.tar.gz` or `go run ./main.go download <version> > cqlsh.tar.gz`

Help: `cqlshvm` or `cqlshvm help` or `go run ./main.go help` or `go run ./main.go`