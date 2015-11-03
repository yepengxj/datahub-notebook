package client

import (
	"bytes"
	"fmt"
	"github.com/asiainfoLDP/datahub-client/cmd"
	"github.com/asiainfoLDP/datahub-client/utils/readline"
	"os/exec"
	"strings"
)

func RunClient() {

	fmt.Println("Welcome to datahub(v0.1.0)")
	fmt.Println("Last login: Wed Oct 16 23:38:50 2015")

	for cmd.Running == true {
		result := readline.ReadLine(&cmd.Prompt)
		if result == nil { // exit loop
			break
		}

		//lineSlice := strings.Split(*result, " ")
		lineSlice := strings.Fields(*result)
		if len(lineSlice) < 1 {
			continue
		}

		commandFound := false
		for _, v := range cmd.Cmds {
			if strings.EqualFold(v.CmdName, lineSlice[0]) {
				commandFound = true
				v.Handler(lineSlice, v)
			}

		}

		if !commandFound && len(lineSlice[0]) > 0 {
			if lineSlice[0][0] == '!' {
				cmd := exec.Command(strings.Trim(lineSlice[0], "!"), lineSlice[1:]...)
				var out bytes.Buffer
				cmd.Stdout = &out
				cmd.Run()
				fmt.Printf("%v", out.String())

			} else if len(lineSlice[0]) == 1 && lineSlice[0][0] == '?' {
				ShowUsage()
			} else {
				fmt.Println("no such command, type '?' to show help.")
			}
		}

		readline.AddHistory(*result) //allow user to recall this line
	}
}

func ShowUsage() {
	for _, v := range cmd.Cmds {
		fmt.Printf("%-16s  %s\n", v.CmdName, v.CmdHelper)

	}
}
