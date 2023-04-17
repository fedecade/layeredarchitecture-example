package server

import (
	"fmt"

	"example.layeredarch/waf"
)

type impl struct {
	port   string
	router waf.Router
}

func New(listenPort string, router waf.Router) waf.Server {
	port := fmt.Sprintf(":%s", listenPort)
	return &impl{port, router}
}
