package client

import (
	"github.com/asiainfoLDP/datahub-client/cmd/dp"
	"github.com/asiainfoLDP/datahub-client/cmd/repo"
	"github.com/asiainfoLDP/datahub-client/cmd/subscription"
)

func init() {
	dp.Handler()
	repo.Handler()
	subscription.Handler()
}
