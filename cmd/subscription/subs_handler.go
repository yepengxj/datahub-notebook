package subscription

import (
	"github.com/asiainfoLDP/datahub-client/cmd"
)

func Handler() {
	cmd.CmdParser["subs info"] = SubsInfo
	cmd.CmdParser["subs pull"] = SubsPull

}
