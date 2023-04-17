package defs

import (
	"os"

	"example.layeredarch/waf"
	"example.layeredarch/waf/server"
	"github.com/sarulabs/di"
)

func Server() di.Def {
	return di.Def{
		Name:  "server",
		Scope: di.App,
		Build: func(cnt di.Container) (any, error) {
			port := os.Getenv("LISTEN_PORT")
			router, err := cnt.SafeGet(RequestRouter().Name)
			if err != nil {
				return nil, err
			}
			return server.New(
				port,
				router.(waf.Router),
			), nil
		},
	}
}
