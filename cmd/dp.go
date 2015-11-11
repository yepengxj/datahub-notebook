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
			err = errors.New(string(body))
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
					fmt.Println(string(body))
					err = errors.New(string(body))
				}
			}
		}
	}
	return nil
}

func dpResp(bDetail bool, byRespBody []byte) {
	if bDetail == true {
		//support: dp name1 name2 name3
		strcDp := FormatDp_dpname{}
		err := json.Unmarshal(byRespBody, &strcDp)
		if err != nil {
			fmt.Println("Get /datapools/:dpname  dpResp json.Unmarshal error!")
		}
		n, _ := fmt.Printf("datapool:%-16s\t%-16s\t%-16s\n", strcDp.Name, strcDp.Type, strcDp.Conn)
		printDash(n + 12)
		for _, item := range strcDp.Items {
			if item.Publish == "Y" {
				fmt.Printf("%s/%s/%s\t%-16s\t%4s\n", item.Repository, item.DataItem, item.Tag, item.Time, "pub")
			} else {
				fmt.Printf("%s/%s/%s\t%-16s\t%4s\n", item.Repository, item.DataItem, item.Tag, item.Time, "pull")
			}

		}
	} else {
		strcDps := []FormatDp{}
		err := json.Unmarshal(byRespBody, &strcDps)
		if err != nil {
			fmt.Println("Get /datapools  dpResp json.Unmarshal error!")
		}
		n, _ := fmt.Printf("%-16s\t%-8s\n", "datapool", "type")
		printDash(n + 12)
		for _, tag := range strcDps {
			fmt.Printf("%-16s\t%-8s\n", tag.Name, tag.Type)
		}
	}
}
