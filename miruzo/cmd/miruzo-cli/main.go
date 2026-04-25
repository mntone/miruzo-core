package main

import "os"

var version = "0.0.0+dev"

func main() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
