package customermanagement

import (
	"example.layeredarch/domain/customer"
	"example.layeredarch/logger"
)

func (i *impl) Register(
	customer customer.Customer,
) error {
	if err := i.creditScreening.Perform(customer); err != nil {
		logger.Error(err)
		return err
	}

	if err := i.customerRepository.Register(customer); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
