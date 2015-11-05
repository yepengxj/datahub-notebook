package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub/utils/mflag"
	"io/ioutil"
)

type FormatDpCreate struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Conn string `json:"conn"`
}

func DpCreate(needLogin bool, args []string) (err error) {
	f := mflag.NewFlagSet("dp create", mflag.ContinueOnError)
	d := FormatDpCreate{}
	f.StringVar(&d.Type, []string{"-type", "t"}, "file", "datapool type")
	f.StringVar(&d.Conn, []string{"-conn"}, "", "datapool connection info")

	if len(args) > 0 && args[0][0] != '-' {
		d.Name = args[0]
		args = args[1:]
	}

	if len(args) == 0 {

		fmt.Printf("missing argument.\nSee '%s --help'.\n", f.Name())
		return
	}

	if err = f.Parse(args); err != nil {
		return
	}

	if len(f.Args()) > 0 {
		fmt.Printf("invalid argument.\nSee '%s --help'.\n", f.Name())
		return

	}

	jsonData, err := json.Marshal(d)
	if err != nil {
		return err
	}

	if needLogin && !Logged {
		login(false)
	}

	resp, err := commToDaemon("post", "/datapool", jsonData)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	return nil
}
