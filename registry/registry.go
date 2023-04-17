package registry

import (
	"example.layeredarch/registry/defs"
	"github.com/sarulabs/di"
)

func New() (di.Container, error) {
	return build(
		defs.Server(),
		defs.RequestRouter(),
		defs.Database(),
		defs.Transaction(),
		defs.DomainCustomerBuilder(),
		defs.RdbCustomerReporitory(),
		defs.HttpCreditScreening(),
		defs.HttpCreditScreeningTranslator(),
		defs.RequestHandlerCustomerPost(),
		defs.RequestHandlerCustomerPostTranslator(),
		defs.UsecaseCustomerManagement(),
		defs.CreditScreeningHttpClient(),
	)
}
