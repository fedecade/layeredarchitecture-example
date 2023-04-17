package unqualifiedcustomer

import (
	"fmt"

	"example.layeredarch/domain/customer"
)

type Error struct {
	customer customer.Customer
	msg      string
}

func (i *Error) Error() string {
	return fmt.Sprintf(
		"name: %s, email: %s, reason: %s",
		i.customer.Name(),
		i.customer.Email(),
		i.msg,
	)
}

func New(customer customer.Customer, msg string) *Error {
	return &Error{customer: customer, msg: msg}
}
