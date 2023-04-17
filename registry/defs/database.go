package defs

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"github.com/sarulabs/di"
)

func Database() di.Def {
	const driver = "mysql"
	const dsntpl = "%s:%s@tcp(%s:%s)/%s?parseTime=true&time_zone=%%27Asia%%2FTokyo%%27&loc=Local&multiStatements=true"
	return di.Def{
		Name:  "database",
		Scope: di.App,
		Build: func(cnt di.Container) (interface{}, error) {
			host := os.Getenv("DB_HOST")
			port := os.Getenv("DB_PORT")
			user := os.Getenv("DB_USER")
			pass := os.Getenv("DB_PASS")
			name := os.Getenv("DB_NAME")
			dsn := fmt.Sprintf(dsntpl, user, pass, host, port, name)
			db, err := sqlx.Connect(driver, dsn)
			if err != nil {
				return nil, err
			}
			db.SetMaxIdleConns(0)
			return db, nil
		},
	}
}
