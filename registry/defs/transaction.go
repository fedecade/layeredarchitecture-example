package defs

import (
	"github.com/jmoiron/sqlx"
	"github.com/sarulabs/di"
)

func Transaction() di.Def {
	return di.Def{
		Name:  "transaction",
		Scope: di.Request,
		Build: func(cnt di.Container) (interface{}, error) {
			if db, err := cnt.SafeGet(Database().Name); err != nil {
				return nil, err
			} else if tx, err := db.(*sqlx.DB).Beginx(); err != nil {
				return nil, err
			} else {
				return tx, nil
			}
		},
	}
}
