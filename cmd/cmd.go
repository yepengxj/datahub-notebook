package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/asiainfoLDP/datahub-client/cmd/dp"
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
}

var (
	User     = UserInfo{}
	Matches  = make([]string, 0, len(Cmds))
	UnixSock = "/var/run/datahub.sock"
	Running  = true
	Logged   = false
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
	User.password = string(pass)
	Logged = true
}
func commandLineArgsParser(args []string, path, method string) {
	command := strings.Join(args[:2], " ")

	var jsonData []byte
	var err error
	//fmt.Println("command is:", command)
	if f, ok := cmdParser[command]; ok {
		if jsonData, err = f(command, args[2:]); err != nil {
			//fmt.Println(err.Error())
			return
		}
	} else {
		fmt.Println("command not implentment yet.")
		return
	}

	commToDaemon(method, path, jsonData)

}

func commToDaemon(method, path string, jsonData []byte) {
	//fmt.Println(method, path, string(jsonData))

	req, err := http.NewRequest(strings.ToUpper(method), path, bytes.NewBuffer(jsonData))
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
	fmt.Printf(string(body))
	//fmt.Printf(Prompt)
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

var cmdParser = make(map[string]func(string, []string) ([]byte, error))

func CmdParserInit() {
	cmdParser["dp create"] = dp.DpCreate
}
