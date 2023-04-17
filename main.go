package main

import (
	"log"
	"os"
	"time"

	"example.layeredarch/logger"
	"example.layeredarch/registry"
	"example.layeredarch/registry/defs"
	"example.layeredarch/waf"
	"github.com/comail/colog"
	"github.com/sarulabs/di"
)

func init() {
	time.Local = time.FixedZone("JST", +9*60*60)
	log.SetPrefix("[LAYERD ARCH EXAMPLE] ")
	colog.Register()
}

func main() {
	reg, err := registry.New()
	if err != nil {
		os.Exit(1)
	}

	logger.Info("Start.")

	if err := runServer(reg); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func runServer(reg di.Container) error {

	server, err := getServer(reg)
	if err != nil {
		return nil
	}

	router, err := getRouter(reg)
	if err != nil {
		return nil
	}

	registry.RegisterRouter(router)

	return server.Run()
}

func getServer(reg di.Container) (waf.Server, error) {
	server, err := func() (waf.Server, error) {
		o, e := reg.SafeGet(defs.Server().Name)
		if e != nil {
			logger.Error(e)
			return nil, e
		}
		return o.(waf.Server), nil
	}()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return server, nil
}

func getRouter(reg di.Container) (waf.Router, error) {
	router, err := func() (waf.Router, error) {
		o, e := reg.SafeGet(defs.RequestRouter().Name)
		if e != nil {
			logger.Error(e)
			return nil, e
		}
		return o.(waf.Router), nil
	}()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return router, nil
}
