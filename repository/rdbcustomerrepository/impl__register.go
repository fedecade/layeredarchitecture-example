package rdbcustomerrepository

import (
	"example.layeredarch/domain/customer"
	"example.layeredarch/domain/errors/alreadyexistcustomer"
	"example.layeredarch/domain/errors/unexpected"
	"example.layeredarch/logger"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func (i *impl) Register(customer customer.Customer) error {
	sql := `
INSERT INTO customers (
 name
,email
) values (
 :Name
,:Email
)
`
	qry, param, _ := sqlx.Named(
		sql,
		map[string]any{
			"Name":  customer.Name(),
			"Email": customer.Email(),
		},
	)

	if _, err := i.tx.Exec(qry, param...); err != nil {
		switch e := err.(type) {
		case *mysql.MySQLError:
			if e.Number == 1062 {
				logger.Error(e)
				return alreadyexistcustomer.New(customer)
			} else {
				logger.Error(e)
				return unexpected.New(err)
			}
		default:
			logger.Error(e)
			return unexpected.New(err)
		}
	} else {
		return nil
	}
}
