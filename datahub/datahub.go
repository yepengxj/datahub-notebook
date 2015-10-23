package datahub

import (
	"fmt"
	"github.com/asiainfoLDP/datahub-client/cmd"
	flag "github.com/asiainfoLDP/datahub-client/utils/mflag"
	"github.com/asiainfoLDP/datahub-client/utils/readline"
	"os"
)

var (
	RunDaemon bool
)

func flagParse() {
	flDaemon := flag.Bool([]string{"D", "-daemon"}, false, "Enable daemon mode")
	flPort := flag.String([]string{"P", "-port"}, "18000", "Binding port")
	flVersion := flag.Bool([]string{"V", "-version"}, false, "Show version")
	flHost := flag.String([]string{"H", "-host"}, "hub.dataos.io:8080", "Server host address")

	flag.Parse()

	fmt.Printf("run daemon: %v, listening port: %v, version: %v, host: %s\n",
		*flDaemon, *flPort, *flVersion, *flHost)

	if len(flag.Args()) > 0 {
		fmt.Println(flag.Args()[0], "option not recongized.")
		cmd.Running = false
	}

	if *flVersion {
		fmt.Println("datahub v0.1.0")
		os.Exit(0)
	}

	if *flDaemon {
		RunDaemon = true
	}
}

func DatahubInit() {
	cmd.Running = true
	readline.SetCompletionEntryFunction(cmd.CompletionEntry)
	readline.SetAttemptedCompletionFunction(nil)
	flagParse()
	cmd.CmdParserInit()
}
