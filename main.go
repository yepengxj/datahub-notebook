package main

import (
	"github.com/asiainfoLDP/datahub-client/client"
	"github.com/asiainfoLDP/datahub-client/daemon"
	"github.com/asiainfoLDP/datahub-client/datahub"
)

func main() {

	datahub.DatahubInit()

	if datahub.RunDaemon {
		daemon.RunDaemon()
	} else {
		client.RunClient()
	}
}
