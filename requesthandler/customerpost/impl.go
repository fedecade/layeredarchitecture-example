package customerpost

import (
	"example.layeredarch/requesthandler"
	"example.layeredarch/requesthandler/customerpost/translator"
	"example.layeredarch/usecase/customermanagement"
	"example.layeredarch/waf/defaulthandler"
)

type impl struct {
	*defaulthandler.DefaultHandler
	cutomerManagement customermanagement.CustomerManagement
	translator        translator.Translator
}

func New(
	customerManagement customermanagement.CustomerManagement,
	translator translator.Translator,
) requesthandler.RequestHandler {
	return &impl{
		defaulthandler.New(),
		customerManagement,
		translator,
	}
}
