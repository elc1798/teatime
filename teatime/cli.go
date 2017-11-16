package main

import (
	"github.com/urfave/cli"
	"os"
    "github.com/elc1798/teatime/p2p"
    //fs "github.com/elc1798/teatime/fs"
	//tt "github.com/elc1798/teatime"
	//dmp "github.com/sergi/go-diff/diffmatchpatch"
	"fmt"
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
        {   Name:   "init",
			Action: func(c *cli.Context) error {
                wd, _ := os.Getwd()
                host := c.Args().Get(0)
                repo, _ := fs.InitRepo(host, wd)

                serverSession := p2p.NewTTNetSession(repo)
                port := 12345
                serverSession.StartListener(port, true)
                if c.Args().Get(1) != "" {
                    _ := testSession.TryTeaTimeConn(fmt.Sprintf("%s:%d", host, port), time.Millisecond*250)
                }
                return nil
			},
        },
        {   Name:   "start",
			Action: func(c *cli.Context) error {
                repos := fs.Load 
                func (this *Repo) AddFile(relativePath string) error {
                repo.StartPollingRepo()
                return nil
			},
        },
        {   Name:   "track",
			Action: func(c *cli.Context) error {
                //wd, _ := os.Getwd()
                fs.AddTrackedFile(wd + "/" + c.Args().Get(0))
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
