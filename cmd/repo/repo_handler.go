package repo

import (
	"github.com/asiainfoLDP/datahub-client/cmd"
)

func Handler() {
	cmd.CmdParser["repo list"] = RepoList
}
