package repo

import (
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub-client/cmd"
	"github.com/asiainfoLDP/datahub-client/utils/mflag"
)

type formatRepoList struct {
	repoName string `json:"repoName"`
	itemName string `json:"itemname"`
}

func RepoList(cmd string, args []string) (j []byte, err error) {
	f := mflag.NewFlagSet(cmd, mflag.ContinueOnError)

	d := formatRepoList{}
	f.StringVar(&d.repoName, []string{"-repo"}, "", "Repository name")
	f.StringVar(&d.itemName, []string{"-item"}, "", "Item name")

	if err = f.Parse(args); err != nil {
		//fmt.Printf(Prompt)
		return nil, err
	}

	if j, err = json.Marshal(d); err != nil {
		return nil, err
	}

	if len(f.Args()) > 0 {
		fmt.Println(f.Args()[0], "not supported.")
		fmt.Printf("See '%s --help' for detail.\n", cmd)
		//fmt.Printf(Prompt)
		return nil, fmt.Errorf("invalid argument")
	}
	return j, err
}

func init() {
	cmd.CmdParser["repo list"] = RepoList
}
