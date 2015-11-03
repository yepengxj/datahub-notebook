package subscription

import (
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub-client/ds"
	"github.com/asiainfoLDP/datahub-client/utils/mflag"
)

func SubsInfo(cmd string, args []string) (j []byte, err error) {
	f := mflag.NewFlagSet(cmd, mflag.ContinueOnError)

	d := ds.FormatRepoList{}
	f.StringVar(&d.RepoName, []string{"-repo"}, "", "Repository name")
	f.StringVar(&d.ItemID, []string{"-id"}, "", "Item name")

	if err = f.Parse(args); err != nil {
		//fmt.Printf(Prompt)
		return nil, err
	}

	if len(f.Args()) > 0 && len(d.ItemID) == 0 {
		d.ItemID = f.Args()[0]
	} else if len(f.Args()) > 0 {
		fmt.Println(f.Args(), "not supported.")
		fmt.Printf("See '%s --help' for detail.\n", cmd)
		//fmt.Printf(Prompt)
		return nil, fmt.Errorf("invalid argument")
	}

	if j, err = json.Marshal(d); err != nil {
		return nil, err

	}

	return j, err
}
