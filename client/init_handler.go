package client

import (
	"github.com/asiainfoLDP/datahub-client/cmd/dp"
	"github.com/asiainfoLDP/datahub-client/cmd/repo"
)

func init() {
	dp.Handler()
	repo.Handler()
}
