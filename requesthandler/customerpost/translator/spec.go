package translator

import (
	"net/http"

	"example.layeredarch/domain/customer"
)

type Translator interface {
	ToCustomer(r *http.Request) (customer.Customer, error)
}
