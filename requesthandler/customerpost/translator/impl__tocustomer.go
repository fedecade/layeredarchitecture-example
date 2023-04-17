package translator

import (
	"encoding/json"
	"net/http"

	"example.layeredarch/domain/customer"
)

func (i *impl) ToCustomer(
	r *http.Request,
) (
	customer.Customer,
	error,
) {
	var rc struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&rc); err != nil {
		return nil, err
	}

	return i.customerBuilder.New(rc.Name, rc.Email)
}
