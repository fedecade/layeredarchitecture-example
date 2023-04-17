package defs

import (
	"example.layeredarch/domain/customer"
	"example.layeredarch/requesthandler/customerpost/translator"
	"github.com/sarulabs/di"
)

func RequestHandlerCustomerPostTranslator() di.Def {
	return di.Def{
		Name:  "requesthandler/customerpost/translator",
		Scope: di.App,
		Build: func(cnt di.Container) (any, error) {
			cb, err := cnt.SafeGet(DomainCustomerBuilder().Name)
			if err != nil {
				return nil, err
			}
			return translator.New(
				cb.(customer.Builder),
			), nil
		},
	}
}
