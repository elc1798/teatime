package main

import (
	"fmt"
	"os"

	fs "github.com/elc1798/teatime/fs"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		return
	}

	cmd := args[0]
	switch cmd {
	case "track":
		err := fs.AddTrackedFile(args[1])
		printErrOrSuccess(err)
	case "back":
		err := fs.WriteBackupFile(args[1])
		printErrOrSuccess(err)
	case "help":
		printHelp()
	default:
		printUsage()
	}
}

func printErrOrSuccess(err error) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Success!")
	}
}

func printHelp() {
    fmt.Println("Commands")
    fmt.Println("--------")
    fmt.Println("track\t[path_to_file]")
    fmt.Println("back\t[tracked_filename]")
    fmt.Println("--------")
}

func printUsage() {
	fmt.Println("Usage: [exe] [cmd] [args]")
	fmt.Println("For list of commands, run cmd 'help'")
}
