package main

import (
	"flag"
	"fmt"
	"os"
)

func printUsage() {
	fmt.Fprintf(
		flag.CommandLine.Output(),
		"Usage: %s [-s | <prompt>]\n" +
		"Askpass helper for headless sessions, enter password in another terminal.\n" +
		"  -s  Start a server to provide password from.\n" +
		"  Without -s, ask an existing server for a password with <prompt>.\n",
		os.Args[0])
}

func main() {
	var serverMode bool
	flag.BoolVar(&serverMode, "s", false, "")
	flag.CommandLine.Usage = printUsage
	flag.Parse()

	if serverMode {
		switch len(flag.Args()) {
		case 0:
			RunServer()
		default:
			printUsage()
			os.Exit(1)
		}
		RunServer()
	} else {
		switch len(flag.Args()) {
		case 0:
			RunClient("")
		case 1:
			RunClient(flag.Args()[0])
		default:
			printUsage()
			os.Exit(1)
		}
	}
}
