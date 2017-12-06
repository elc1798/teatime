package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/urfave/cli"

	tt "github.com/elc1798/teatime"
	encoder "github.com/elc1798/teatime/encode"
)

func main() {
	app := cli.NewApp()
	app.Name = "teatime"
	app.Usage = "teatime [command] <args>"
	app.Action = func(c *cli.Context) error {
		fmt.Println(app.Usage)
		return nil
	}
	app.Commands = []cli.Command{{
		Name:  "start",
		Usage: "Starts the Teatime Crumpet Daemon. Note: This operation hangs!",
		Action: func(c *cli.Context) error {
			startCrumpetAndHang()
			return nil
		},
	}, {
		Name:  "reset",
		Usage: "Resets Teatime metadata",
		Action: func(c *cli.Context) error {
			tt.ResetTeatime()
			return nil
		},
	}, {
		Name:  "init",
		Usage: "Marks the target directory as repository",
		Action: func(c *cli.Context) error {
			sendCrumpetCommand(append([]string{encoder.COMMAND_INIT_REPO}, os.Args[2:]...))
			return nil
		},
	}, {
		Name:  "connect",
		Usage: "Connects the desired repo to the desired peer",
		Action: func(c *cli.Context) error {
			sendCrumpetCommand(append([]string{encoder.COMMAND_LINK_PEER}, os.Args[2:]...))
			return nil
		},
	}}

	sort.Sort(cli.CommandsByName(app.Commands))
	app.Run(os.Args)
}
