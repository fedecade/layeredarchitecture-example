package router

import (
	"net/http"

	"example.layeredarch/logger"
	"example.layeredarch/registry/errors/incorrectregtype"
	"example.layeredarch/requesthandler"
	"example.layeredarch/waf"
	"github.com/jmoiron/sqlx"
	"github.com/sarulabs/di"
)

func (i *impl) Register(defs ...waf.HandlerDef) {
	for _, def := range defs {
		i.router.HandleFunc(def.Path, i.handlerFunc(def.Name)).Methods(def.Method.Name)
	}
}

func (i *impl) handlerFunc(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctn, err := i.rootContainer.SubContainer()
		if err != nil {
			i.responseError(err, w)
			return
		}
		defer func() {
			if err := ctn.Delete(); err != nil {
				logger.Error(err)
			}
		}()

		handler, err := i.getRequestHandler(ctn, name)
		if err != nil {
			i.responseError(err, w)
			return
		}

		tx, err := i.beginTx(ctn)
		if err != nil {
			i.responseError(err, w)
			return
		}

		if err := handler.HandleRequest(w, r); err != nil {
			logger.Error(err)
			if e := tx.Rollback(); e != nil {
				logger.Error(e)
			}
		} else {
			if err := tx.Commit(); err != nil {
				logger.Error(err)
			}
		}
	}
}

func (i *impl) getRequestHandler(ctn di.Container, name string) (requesthandler.RequestHandler, error) {
	obj, err := ctn.SafeGet(name)
	if err != nil {
		return nil, err
	}
	handler, ok := obj.(requesthandler.RequestHandler)
	if !ok {
		return nil, incorrectregtype.New(name, "requesthandler.RequestHandler")
	}

	return handler, nil
}

func (i *impl) beginTx(ctn di.Container) (*sqlx.Tx, error) {
	name := "transaction"
	obj, err := ctn.SafeGet(name)
	if err != nil {
		return nil, err
	}
	tx, ok := obj.(*sqlx.Tx)
	if !ok {
		return nil, incorrectregtype.New(name, "requesthandler.RequestHandler")
	}

	return tx, nil
}

func (i *impl) responseError(err error, w http.ResponseWriter) {
	body := []byte(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Content-Length", string(rune(len(body))))
	w.Write(body)
}
