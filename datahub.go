package main

import (
	"fmt"
	"github.com/asiainfoLDP/datahub-client/client"
	"github.com/asiainfoLDP/datahub-client/daemon"
	flag "github.com/asiainfoLDP/datahub-client/utils/mflag"
	"os"
)

var (
	runDaemon bool
)

func init() {
	flagParse()
}

func flagParse() {
	flDaemon := flag.Bool([]string{"D", "-daemon"}, false, "Enable daemon mode")
	flVersion := flag.Bool([]string{"V", "-version"}, false, "Show version")

	flag.Usage = client.ShowUsage
	//flag.PrintDefaults()
	flag.Parse()
	//fmt.Printf("run daemon: %v, version: %v\n", *flDaemon, *flVersion)

	if *flVersion {
		fmt.Println("datahub v0.1.0")
		os.Exit(0)
	}

	if *flDaemon {
		runDaemon = true
	}
}

func main() {

	if runDaemon {
		daemon.RunDaemon()
	} else {
		client.RunClient()
	}
}
