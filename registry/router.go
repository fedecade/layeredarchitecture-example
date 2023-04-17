package registry

import (
	"example.layeredarch/registry/defs"
	"example.layeredarch/waf"
	"example.layeredarch/waf/httpmethod"
)

func RegisterRouter(router waf.Router) {
	router.Register(
		waf.HandlerDef{
			Method: httpmethod.Post,
			Path:   "/customer",
			Name:   defs.RequestHandlerCustomerPost().Name,
		},
	)
}
