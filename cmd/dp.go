package cmd

import (
    "fmt"
	"encoding/json"
)

type FormatDp struct {
	Name string `json:"dpname"`
	Type string `json:"dptype"`
	//Conn string `json:"dpconn"`
}

type Item struct{
	Repository string `json:"repository"`
	DataItem   string `json:"dataitem"`
	Tag        string `json:"tag"`
	Time       string `json:"time"`
	Publish    string `json:"publish"`
}
type FormatDp_dpname struct {
	Name string `json:"dpname"`
	Type string `json:"dptype"`
	Conn string `json:"dpconn"`
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
		res_body := commToDaemon("/datapool", "GET", jsonData)
		fmt.Println(string(res_body))
	}else{
		for _, v := range args {
		    d.Name = v
		    if v[0] != '-' {
		        str_dp := fmt.Sprintf("/datapool/%s", d.Name)
		        jsonData, err := json.Marshal(d)
	            if err != nil {
		            return err
	            }
	            fmt.Println(jsonData)
		        res_body := commToDaemon(str_dp, "GET", jsonData)
		        fmt.Println(string(res_body))
		    }
		}
		
	}
	return nil
}

