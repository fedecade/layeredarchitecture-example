package customermanagement

import (
	"example.layeredarch/domain/creditscreening"
	"example.layeredarch/domain/customer"
)

type impl struct {
	customerRepository customer.Repository
	creditScreening    creditscreening.CreditScreeing
}

func New(
	customerRepository customer.Repository,
	creditScreening creditscreening.CreditScreeing,
) CustomerManagement {
	return &impl{
		customerRepository: customerRepository,
		creditScreening:    creditScreening,
	}
}
