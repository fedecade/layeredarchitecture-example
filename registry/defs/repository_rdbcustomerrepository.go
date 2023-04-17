package defs

import (
	"example.layeredarch/logger"
	"example.layeredarch/repository/rdbcustomerrepository"
	"github.com/jmoiron/sqlx"
	"github.com/sarulabs/di"
)

func RdbCustomerReporitory() di.Def {
	return di.Def{
		Name:  "repository/rdbcustomerrepository",
		Scope: di.Request,
		Build: func(cnt di.Container) (any, error) {
			tx, err := cnt.SafeGet(Transaction().Name)
			if err != nil {
				logger.Error(err)
				return nil, err
			}
			return rdbcustomerrepository.New(
				tx.(*sqlx.Tx),
			), nil
		},
	}
}
