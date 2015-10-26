package dp

import (
	"github.com/asiainfoLDP/datahub-client/cmd"
)

func Handler() {
	cmd.CmdParser["dp create"] = DpCreate
}
