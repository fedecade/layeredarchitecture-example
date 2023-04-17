package defs

import (
	"example.layeredarch/domain/creditscreening"
	"example.layeredarch/domain/customer"
	"example.layeredarch/logger"
	"example.layeredarch/usecase/customermanagement"
	"github.com/sarulabs/di"
)

func UsecaseCustomerManagement() di.Def {
	return di.Def{
		Name:  "usecase/customermanagement",
		Scope: di.Request,
		Build: func(cnt di.Container) (any, error) {
			repo, err := cnt.SafeGet(RdbCustomerReporitory().Name)
			if err != nil {
				logger.Error(err)
				return nil, err
			}
			gtwy, err := cnt.SafeGet(HttpCreditScreening().Name)
			if err != nil {
				logger.Error(err)
				return nil, err
			}
			return customermanagement.New(
				repo.(customer.Repository),
				gtwy.(creditscreening.CreditScreeing),
			), nil
		},
	}
}
