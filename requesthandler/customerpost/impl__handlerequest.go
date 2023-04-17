package customerpost

import (
	"net/http"

	"example.layeredarch/domain/errors/alreadyexistcustomer"
	"example.layeredarch/domain/errors/unqualifiedcustomer"
	"example.layeredarch/logger"
)

func (i *impl) HandleRequest(w http.ResponseWriter, r *http.Request) error {
	customer, err := i.translator.ToCustomer(r)
	if err != nil {
		logger.Error(err)
		i.ResponseError(http.StatusBadRequest, w, err)
		return err
	}

	if err := i.cutomerManagement.Register(customer); err != nil {
		switch e := err.(type) {
		case *unqualifiedcustomer.Error:
			logger.Error(e)
			i.ResponseError(http.StatusBadRequest, w, e)
		case *alreadyexistcustomer.Error:
			logger.Error(e)
			i.ResponseError(http.StatusConflict, w, e)
		default:
			logger.Error(e)
			i.ResponseError(http.StatusInternalServerError, w, e)
		}
		return err
	}

	logger.Info("Customer created. [name: %s, email: %s]", customer.Name(), customer.Email())
	i.ResponseEmpty(http.StatusCreated, w)
	return nil
}
