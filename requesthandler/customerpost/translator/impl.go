package translator

import "example.layeredarch/domain/customer"

type impl struct {
	customerBuilder customer.Builder
}

func New(customerBuilder customer.Builder) Translator {
	return &impl{customerBuilder: customerBuilder}
}
