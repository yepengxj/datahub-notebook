package cmd

import (
	"fmt"
	//"github.com/asiainfoLDP/datahub-client/ds"
)

type Commands struct {
	CmdName   string
	subCmd    []cmdMethod
	path      string
	Handler   func([]string, Commands)
	CmdHelper string
	needLogin bool
}
type cmdMethod struct {
	cmd    string
	method string
}

type UserInfo struct {
	userName string
	password string
	b64      string
}
type Data struct {
	Item  DataItem
	Usage DataItemUsage
}
type MsgResp struct {
	Msg string `json:"msg"`
}
type Repository struct {
	Repository_id   int    `json:"repository_id,omitempty"`
	Repository_name string `json:"repository_name,omitempty"`
	User_id         int    `json:"user_id,omitempty"`
	Permit_type     int    `json:"permit_type,omitempty"`
	Arrange_type    int    `json:"arrange_type,omitempty"`
	Comment         string `json:"comment,omitempty"`
	Rank            int    `json:"rank,omitempty"`
	Status          int    `json:"status,omitempty"`
	Dataitems       int    `json:"dataitems,omitempty"`
	Tags            int    `json:"tags,omitempty"`
	Stars           int    `json:"stars,omitempty"`
	Optime          string `json:"optime,omitempty"`
}
type DataItem struct {
	Repository_id   int     `json:"repository_id,omitempty"`
	User_id         int     `json:"user_id,omitempty"`
	Dataitem_id     int     `json:"dataitem_id,omitempty"`
	Dataitem_name   string  `json:"dataitem_name,omitempty"`
	Ico_name        string  `json:"ico_name,omitempty"`
	Permit_type     int     `json:"permit_type,omitempty"`
	Key_words       string  `json:"key_words,omitempty"`
	Supply_style    int     `json:"supply_style,omitempty"`
	Priceunit_type  int     `json:"priceunit_type,omitempty"`
	Price           float32 `json:"price,omitempty"`
	Optime          string  `json:"optime,omitempty"`
	Data_format     int     `json:"data_format,omitempty"`
	Refresh_type    int     `json:"refresh_type,omitempty"`
	Refresh_num     int     `json:"refresh_num,omitempty"`
	Meta_filename   string  `json:"meta_filename,omitempty"`
	Sample_filename string  `json:"sample_filename,omitempty"`
	Comment         string  `json:"comment,omitempty"`
}

type DataItemUsage struct {
	Dataitem_id   int    `json:"-,omitempty"`
	Dataitem_name string `json:"-,omitempty"`
	Views         int    `json:"views"`
	Follows       int    `json:"follows"`
	Downloads     int    `json:"downloads"`
	Stars         int    `json:"stars"`
	Refresh_date  string `json:"refresh_date,omitempty"`
	Usability     int    `json:"usability,omitempty"`
}

type RepoJson struct {
	Datas []Data
	Total int
}

var (
	User          = UserInfo{}
	UnixSock      = "/var/run/datahub.sock"
	DefaultServer = "http://10.1.51.32:8080"
	Logged        = false
)

type Command struct {
	Name      string
	SubCmd    []Command
	Handler   func(login bool, args []string) error
	Desc      string
	NeedLogin bool
}

var Cmd = []Command{
	{
		Name:    "dp",
		Handler: Dp,
		SubCmd: []Command{
			{
				Name:    "create",
				Handler: DpCreate,
			},
			{
				Name:    "rm",
				Handler: DpRm,
			},
		},
		Desc: "list all of datapool.",
	},
	{
		Name:      "subs",
		Handler:   Subs,
		Desc:      "subscription of item.",
		NeedLogin: true,
	},
	{
		Name:      "pull",
		Handler:   Pull,
		Desc:      "pull item from peer.",
		NeedLogin: true,
	},
	{
		Name:    "login",
		Handler: Login,
		Desc:    "login in to dataos.io.",
	},
}

func login(interactive bool) {
	if Logged {
		if interactive {
			fmt.Println("you are already logged in as", User.userName)
		}
		return
	}

}
