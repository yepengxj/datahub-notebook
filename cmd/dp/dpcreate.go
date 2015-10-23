package dp

import (
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub-client/utils/mflag"
)

type formatDpCreate struct {
	Dpname  string `json:"dpname"`
	Dptype  string `json:"dptype"`
	Dpconn  string `json:"dpConn"`
	Dpquota string `json:"dpquota"`
}

func DpCreate(cmd string, args []string) (j []byte, err error) {
	f := mflag.NewFlagSet(cmd, mflag.ContinueOnError)
	d := formatDpCreate{}
	f.StringVar(&d.Dpname, []string{"-dpname"}, "", "datapool name")
	f.StringVar(&d.Dptype, []string{"-dptype", "T"}, "FILE", "datapool type")
	f.StringVar(&d.Dpconn, []string{"-dpconn"}, "", "datapool connection info")
	f.StringVar(&d.Dpquota, []string{"-dpquota"}, "", "datapool quota value")
	if err = f.Parse(args); err != nil {
		//fmt.Printf(Prompt)
		return
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
