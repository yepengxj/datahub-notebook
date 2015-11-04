package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub-client/utils/mflag"

	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
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

	commToDaemon("/datapool", jsonData)

	return nil
}

func commToDaemon(path string, jsonData []byte) {
	//fmt.Println(method, path, string(jsonData))

	req, err := http.NewRequest("POST", path, bytes.NewBuffer(jsonData))
	if len(User.userName) > 0 {
		req.SetBasicAuth(User.userName, User.password)
	}
	conn, err := net.Dial("unix", UnixSock)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Datahub Daemon not running?")
		return
	}
	//client := &http.Client{}
	client := httputil.NewClientConn(conn, nil)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//formatResp(cmd, body)
	fmt.Println(string(body))
}
