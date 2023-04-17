package creditscreening

import "example.layeredarch/domain/customer"

type CreditScreeing interface {
	Perform(customer.Customer) error
}
