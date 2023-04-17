package customermanagement

import (
	"example.layeredarch/domain/customer"
)

type CustomerManagement interface {
	Register(customer.Customer) error
}
