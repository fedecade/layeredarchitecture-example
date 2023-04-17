package waf

import (
	"net/http"

	"example.layeredarch/waf/httpmethod"
)

type HandlerDef struct {
	Method httpmethod.Method
	Path   string
	Name   string
}

type Router interface {
	Handler() http.Handler
	Register(...HandlerDef)
}
