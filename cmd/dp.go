package cmd

import (
	"encoding/json"
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
	d := FormatDp{}
	if needLogin && !Logged {
		login(false)
	}

	if len(args) == 0 {
		jsonData, err := json.Marshal(d)
		if err != nil {
			return err
		}
		fmt.Println(jsonData)
		resp, err := commToDaemon("GET", "/datapools", jsonData)
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	} else {
		for _, v := range args {
			d.Name = v
			if v[0] != '-' {
				str_dp := fmt.Sprintf("/datapools/%s", d.Name)
				jsonData, err := json.Marshal(d)
				if err != nil {
					return err
				}
				fmt.Println(jsonData)
				resp, _ := commToDaemon("GET", str_dp, jsonData)
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				fmt.Println(string(body))
			}
		}

	}
	return nil
}
