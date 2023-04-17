package rdbcustomerrepository

import (
	"example.layeredarch/domain/customer"
	"github.com/jmoiron/sqlx"
)

type impl struct {
	tx *sqlx.Tx
}

func New(tx *sqlx.Tx) customer.Repository {
	return &impl{tx: tx}
}
