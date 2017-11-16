package crumpet

import (
	"github.com/elc1798/teatime/fs"
	"github.com/elc1798/teatime/p2p"
)

func Start() {
	x := make(chan bool)
	<-x
	repos, _ := fs.GetAllRepos()
	for i, repo := range repos {
		ttns := p2p.NewTTNetSession(repo)
		ttns.StartListener(12345+i, true)
	}
}
