package main

import (
	"github.com/urfave/cli"
	"os"
	//tt "github.com/elc1798/teatime"
	//dmp "github.com/sergi/go-diff/diffmatchpatch"
	"fmt"
	"sort"
)

func main() {
	app := cli.NewApp()
	app.Name = "teatime"
	app.Usage = "make an explosive entrance"
	app.Action = func(c *cli.Context) error {
		fmt.Println("boom! I say!")
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:    "lorem",
			Aliases: []string{"l"},
			Usage:   "lorem ipsum",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}
	sort.Sort(cli.CommandsByName(app.Commands))
	app.Run(os.Args)
}
