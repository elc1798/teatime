package main

import (
    "fs"
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        return
    }

    fileObj, err := fs.GetFileObjFromFile(os.Args[1])
    if err != nil {
        fmt.Println(err)
        return
    }
    err = fs.WriteFileObjToPath(fileObj, os.Args[1] + ".out")
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("Complete!  Diff input file with .out file to ensure their contents are the same.")
}
