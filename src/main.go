package main

import "github.com/jandedobbeleer/aliae/src/cli"

var (
	Version = "development"
)

func main() {
	cli.Execute(Version)
}
