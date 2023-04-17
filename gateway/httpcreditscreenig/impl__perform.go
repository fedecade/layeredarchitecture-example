package httpcreditscreenig

import (
	"fmt"
	"io"
	"net/http"

	"example.layeredarch/domain/customer"
	"example.layeredarch/domain/errors/unqualifiedcustomer"
	"example.layeredarch/logger"
)

func (i *impl) Perform(customer customer.Customer) error {
	data := i.translator.ToRequestData(customer)

	res, err := i.client.PostJson(data)
	if err != nil {
		logger.Error(err)
		return err
	}

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		res := fmt.Sprintf(
			"StatusCode: %d, Body: %s",
			res.StatusCode,
			func() string {
				if body != nil {
					return string(body)
				} else {
					return ""
				}
			}(),
		)
		return unqualifiedcustomer.New(customer, res)
	}

	return nil
}
