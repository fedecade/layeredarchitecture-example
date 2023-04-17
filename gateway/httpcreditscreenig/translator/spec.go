package translator

import "example.layeredarch/domain/customer"

type Translator interface {
	ToRequestData(customer.Customer) CreditScreening
}
