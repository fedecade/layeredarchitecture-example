package alreadyexistcustomer

import (
	"fmt"

	"example.layeredarch/domain/customer"
)

type Error struct {
	customer customer.Customer
}

func (i *Error) Error() string {
	return fmt.Sprintf(
		"name: %s, email: %s",
		i.customer.Name(),
		i.customer.Email(),
	)
}

func New(customer customer.Customer) *Error {
	return &Error{customer: customer}
}
