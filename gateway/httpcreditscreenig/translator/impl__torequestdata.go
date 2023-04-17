package translator

import "example.layeredarch/domain/customer"

type CreditScreening struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (i *impl) ToRequestData(customer customer.Customer) CreditScreening {
	return CreditScreening{
		customer.Name(),
		customer.Email(),
	}
}
