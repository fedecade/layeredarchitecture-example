package defs

import (
	"example.layeredarch/domain/customer"
	"github.com/sarulabs/di"
)

func DomainCustomerBuilder() di.Def {
	return di.Def{
		Name:  "domain/customer/builder",
		Scope: di.App,
		Build: func(cnt di.Container) (any, error) {
			return customer.NewBuilder(), nil
		},
	}
}
