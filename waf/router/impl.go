package router

import (
	"net/http"

	"example.layeredarch/waf"
	"github.com/gorilla/mux"
	"github.com/sarulabs/di"
)

type impl struct {
	router        *mux.Router
	rootContainer di.Container
}

func (i *impl) Handler() http.Handler {
	return i.router
}

func New(
	rootContainer di.Container,
) waf.Router {
	return &impl{
		router:        mux.NewRouter(),
		rootContainer: rootContainer,
	}
}
