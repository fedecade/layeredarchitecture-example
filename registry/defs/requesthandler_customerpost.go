package defs

import (
	"example.layeredarch/requesthandler/customerpost"
	"example.layeredarch/requesthandler/customerpost/translator"
	"example.layeredarch/usecase/customermanagement"
	"github.com/sarulabs/di"
)

func RequestHandlerCustomerPost() di.Def {
	return di.Def{
		Name:  "requesthandler/customerpost",
		Scope: di.Request,
		Build: func(cnt di.Container) (any, error) {
			uc, err := cnt.SafeGet(UsecaseCustomerManagement().Name)
			if err != nil {
				return nil, err
			}
			tl, err := cnt.SafeGet(RequestHandlerCustomerPostTranslator().Name)
			if err != nil {
				return nil, err
			}
			return customerpost.New(
				uc.(customermanagement.CustomerManagement),
				tl.(translator.Translator),
			), nil
		},
	}
}
