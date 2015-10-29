package cmd

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	//"github.com/asiainfoLDP/datahub-client/ds"
	"github.com/asiainfoLDP/datahub-client/utils"
	"github.com/asiainfoLDP/datahub-client/utils/readline"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
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
	User      = UserInfo{}
	Matches   = make([]string, 0, len(Cmds))
	UnixSock  = "/var/run/datahub.sock"
	Prompt    = "datahub> "
	Running   = true
	Logged    = false
	CmdParser = make(map[string]func(string, []string) ([]byte, error))
)

var Cmds = []Commands{
	{
		CmdName:   "login",
		Handler:   cmdHandlerLogin,
		CmdHelper: "Login to datahub server",
	},
	{
		CmdName:   "logout",
		Handler:   cmdHandlerLogout,
		CmdHelper: "Logout from datahub server",
	},
	{
		CmdName: "dp",
		path:    "/datapool",
		subCmd: []cmdMethod{
			{"create", "post"},
			{"list", "get"},
			{"update", "put"},
			{"delete", "delete"},
		},
		Handler:   cmdHandler,
		CmdHelper: "Datapool management",
	},
	{
		CmdName: "job",
		path:    "/job",
		subCmd: []cmdMethod{
			{"list", "get"},
			{"rm", "delete"},
		},
		Handler:   cmdHandler,
		CmdHelper: "Job management",
	},
	{
		CmdName: "ep",
		path:    "/entrypoint",
		subCmd: []cmdMethod{
			{"add", "post"},
			{"update", "put"},
			{"list", "get"},
			{"delete", "delete"},
		},
		Handler:   cmdHandler,
		CmdHelper: "EntryPoint management",
	},
	{
		CmdName: "subscription",
		path:    "/subscriptions",
		subCmd: []cmdMethod{
			{"queryall", "get"},
			{"query", "get"},
			{"pull", "post"},
			{"pullsingle", "post"},
			{"stream", "post"},
		},
		Handler:   cmdHandler, //cmdHandlerSubscription,
		CmdHelper: "Subscription management",
		needLogin: true,
	},
	{
		CmdName: "repo",
		path:    "/Repository",
		subCmd: []cmdMethod{
			{"list", "get"},
			{"put", "post"},
		},
		Handler:   cmdHandler,
		CmdHelper: "Repository management",
		needLogin: true,
	},
	{
		CmdName:   "quit",
		Handler:   cmdHandlerQuit,
		CmdHelper: "Goodbye",
	},
	{
		CmdName:   "!command",
		Handler:   cmdHandlerCMD,
		CmdHelper: "Execute shell command",
	},
}

func cmdHandlerCMD(args []string, v Commands) {
	fmt.Println("FIXME!")

}

func cmdHandler(args []string, c Commands) {
	if len(args) < 2 {
		fmt.Println("too few arguments.")
		//fmt.Printf("%s %s ...\n", args[0], strings.Join(c.subCmd, "|"))
		fmt.Printf("%s ", args[0])
		for _, v := range c.subCmd {
			fmt.Printf("%s|", v.cmd)
		}
		fmt.Println("...")
		//strings.Join(c.subCmd, "|"))
		return
	}
	if c.needLogin && !Logged {
		Login(false)
	}
	for _, v := range c.subCmd {
		if strings.EqualFold(v.cmd, args[1]) {
			//fmt.Println("TODO:", c.CmdName, v.cmd)
			go JobRequest(args, c.path, v.method)
			return
		}
	}
	fmt.Println(args[1], "operation not supported")
}
func cmdHandlerQuit(args []string, c Commands) {
	Running = false
	fmt.Println("bye.")
}

func cmdHandlerLogin(args []string, c Commands) {
	Login(true)
}

func cmdHandlerLogout(args []string, c Commands) {
	Logout()
}

func Logout() {
	Logged = false
}

func Login(interactive bool) {
	if Logged {
		if interactive {
			fmt.Println("you are already logged in as", User.userName)
		}
		return
	}

	fmt.Printf("login: ")
	reader := bufio.NewReader(os.Stdin)
	//loginName, _ := reader.ReadString('\n')
	loginName, _ := reader.ReadBytes('\n')

	loginName = bytes.TrimRight(loginName, "\r\n")
	fmt.Printf("password: ")
	pass := utils.GetPasswd(true) // Silent, for *'s use gopass.GetPasswdMasked()
	//fmt.Printf("[%s]:[%s]\n", string(loginName), string(pass))

	User.userName = string(loginName)
	//User.password = string(pass)
	User.password = fmt.Sprintf("%x", md5.Sum(pass))
	/*
		User.b64 = base64.StdEncoding.EncodeToString([]byte(User.password + ":" + User.password))
		fmt.Printf("%s\n%s:%s\n", User.b64, User.userName, User.password)
	*/
	Logged = true
}
func commandLineArgsParser(args []string, path, method string) {
	command := strings.Join(args[:2], " ")

	var jsonData []byte
	var err error
	//fmt.Println("command is:", command)
	if f, ok := CmdParser[command]; ok {
		if jsonData, err = f(command, args[2:]); err != nil {
			//fmt.Println(err.Error())
			return
		}
	} else {
		fmt.Println("command not implentment yet.")
		return
	}

	commToDaemon(command, method, path, jsonData)

}

func commToDaemon(cmd, method, path string, jsonData []byte) {
	//fmt.Println(method, path, string(jsonData))

	req, err := http.NewRequest(strings.ToUpper(method), path, bytes.NewBuffer(jsonData))
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
	formatResp(cmd, body)
	//fmt.Printf(Prompt)
}

func formatResp(cmd string, body []byte) {
	//fmt.Println(string(body))
	switch cmd {
	case "repo list":
		var value = RepoJson{}
		type Tag struct {
			Repository     string `json:"-"`
			dataitem       string `json:"-"`
			dataitemusages string `json:"-"`
			Tag            []struct {
				a, b, c string
			} `json:"tags"`
		}
		var value2 = Tag{}
		v2Flag := false

		fmt.Println(string(body))
		if err := json.Unmarshal(body, &value); err != nil {
			fmt.Println("111", err)
			if err := json.Unmarshal(body, &value2); err != nil {

				fmt.Println("222", err)
			} else {
				v2Flag = true
			}

		}

		fmt.Println(value2)
		if !v2Flag {
			n, _ := fmt.Printf("\n%-8s%-24s%-16s\n", "ITEMID", "ITEMNAME", "LASTUPDATE")
			printDash(n)
			for _, v := range value.Datas {
				fmt.Printf("%-8d%-24s%-s\n", v.Item.Dataitem_id, v.Item.Dataitem_name,
					v.Usage.Refresh_date)
				//fmt.Printf("%#v", v)
			}
			printDash(n)
		} else {
			fmt.Printf("hello world")
		}
	default:
		fmt.Println(string(body))
	}
	flush()
}

func JobRequest(args []string, path, method string) {
	commandLineArgsParser(args, path, method)

}

func AttemptedCompletion(text string, start, end int) []string {
	if start == 0 { // this is the command to match
		return readline.CompletionMatches(text, CompletionEntry)
	} else {
		return nil
	}
}

//
// this return a match in the "words" array
//
func CompletionEntry(prefix string, index int) string {
	if index == 0 {
		Matches = Matches[:0]

		for _, w := range Cmds {
			if strings.HasPrefix(w.CmdName, prefix) {
				Matches = append(Matches, w.CmdName)
			}
		}
	}

	if index < len(Matches) {
		return Matches[index]
	} else {
		return ""
	}
}

func flush() {
	fmt.Printf(Prompt)
}
func printDash(n int) {
	for i := 0; i < n-2; i++ {
		fmt.Printf("-")
	}
	fmt.Println()
}
func CmdParserInit() {

	//cmdParser["dp create"] = dp.DpCreate
}

/*
func AddBasicAuth(req *http.Request) *http.Request {
	req.Header.Set("Authorization", "Basic "+User.b64)
	return req
}
*/
