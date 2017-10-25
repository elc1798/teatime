package main

import (
	"fmt"
	"fs"
	"os"
)

func main() {
	args := os.Args[1:]

	err := fs.AddTrackedFile(args[0])
	if err != nil {
		fmt.Println(err)
	} else {
        fmt.Println("Link created!")
    }
}
