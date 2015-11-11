package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type FormatDp struct {
	Name string `json:"dpname"`
	Type string `json:"dptype"`
	//Conn string `json:"dpconn"`
}

type Item struct {
	Repository string `json:"repository"`
	DataItem   string `json:"dataitem"`
	Tag        string `json:"tag"`
	Time       string `json:"time"`
	Publish    string `json:"publish"`
}
type FormatDp_dpname struct {
	Name  string `json:"dpname"`
	Type  string `json:"dptype"`
	Conn  string `json:"dpconn"`
	Items []Item `json:"items"`
}

func Dp(needLogin bool, args []string) (err error) {
	if needLogin && !Logged {
		login(false)
	}

	if len(args) == 0 {
		resp, _ := commToDaemon("GET", "/datapools", nil)
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			dpResp(false, body)
			//fmt.Println(string(body))
		} else {
			fmt.Println(string(body))
			strmsg := ShowMsgResp(body, true)
			err = errors.New(strmsg)
		}

	} else {
		for _, v := range args {
			if v[0] != '-' {
				str_dp := fmt.Sprintf("/datapools/%s", v)
				resp, _ := commToDaemon("GET", str_dp, nil)
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				if resp.StatusCode == 200 {
					dpResp(true, body)
					//fmt.Println(string(body))
				} else {
					strmsg := ShowMsgResp(body, true)
					err = errors.New(strmsg)
				}
			}
		}
	}
	return err
}

func dpResp(bDetail bool, RespBody []byte) {
	if bDetail == true {
		//support: dp name1 name2 name3
		strcDp := FormatDp_dpname{}
		err := json.Unmarshal(RespBody, &strcDp)
		if err != nil {
			fmt.Println("Get /datapools/:dpname  dpResp json.Unmarshal error!")
			return
		}
		n, _ := fmt.Printf("datapool:%-16s\t%-16s\t%-16s\n", strcDp.Name, strcDp.Type, strcDp.Conn)

		for _, item := range strcDp.Items {
			RepoItemTag := item.Repository + "/" + item.DataItem + "/" + item.Tag
			if item.Publish == "Y" {
				fmt.Printf("%-32s\t%-16s\t%-4s\n", RepoItemTag, item.Time, "pub")
			} else {
				fmt.Printf("%-32s\t%-16s\t%-4s\n", RepoItemTag, item.Time, "pull")
			}
		}
		printDash(n + 12)
	} else {
		strcDps := []FormatDp{}
		err := json.Unmarshal(RespBody, &strcDps)
		if err != nil {
			fmt.Println("Get /datapools  dpResp json.Unmarshal error!")
			return
		}
		n, _ := fmt.Printf("%-16s\t%-8s\n", "datapool", "type")
		printDash(n + 12)
		for _, dp := range strcDps {
			fmt.Printf("%-16s\t%-8s\n", dp.Name, dp.Type)
		}
	}
}

func ShowMsgResp(RespBody []byte, bprint bool) (sMsgResp string) {
	msg := MsgResp{}
	err := json.Unmarshal(RespBody, &msg)
	if err != nil {
		sMsgResp = err.Error() + "ShowMsgResp unmarshal error!"
	} else {
		sMsgResp = msg.Msg
		if bprint {
			fmt.Println(sMsgResp)
		}
	}
	return sMsgResp
}
