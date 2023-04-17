package defs

import (
	"example.layeredarch/waf/router"
	"github.com/sarulabs/di"
)

func RequestRouter() di.Def {
	return di.Def{
		Name:  "router",
		Scope: di.App,
		Build: func(cnt di.Container) (any, error) {
			return router.New(cnt), nil
		},
	}
}
