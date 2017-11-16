package main

import (
	"github.com/elc1798/teatime/crumpet"
	"github.com/elc1798/teatime/fs"
	"github.com/elc1798/teatime/p2p"
	"github.com/urfave/cli"
	"os"
	//tt "github.com/elc1798/teatime"
	//dmp "github.com/sergi/go-diff/diffmatchpatch"
	"fmt"
	"path"
	"sort"
	"time"
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
		{Name: "init",
			Action: func(c *cli.Context) error {
				wd, _ := os.Getwd()
				host := c.Args().Get(0)
				repo, _ := fs.InitRepo(host, wd)

				serverSession := p2p.NewTTNetSession(repo)
				port := 2345
				serverSession.StartListener(port, true)
				if c.Args().Get(1) != "" {
					serverSession.TryTeaTimeConn(fmt.Sprintf("%s:%d", host, port), time.Millisecond*250)
				}
				return nil
			},
		},
		{Name: "start",
			Action: func(c *cli.Context) error {
				crumpet.Start()
				return nil
			},
		},
		{Name: "track",
			Action: func(c *cli.Context) error {
				wd, _ := os.Getwd()
				repoName := c.Args().Get(0)
				repo, _ := fs.LoadRepo(repoName)
				fileName := c.Args().Get(1)
				repo.AddFile(path.Join(wd, fileName))
				return nil
			},
		},
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
