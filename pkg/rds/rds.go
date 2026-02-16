package rds

import (
	"database/sql"

	"github.com/andreparelho/order-api/pkg/config"
	"github.com/go-sql-driver/mysql"
)

func GetConnection(cfg config.Configuration) (*sql.DB, error) {
	rdsCfg := mysql.Config{
		User:                 cfg.RDS.User,
		Passwd:               cfg.RDS.Password,
		Net:                  "tcp",
		Addr:                 cfg.RDS.Addr,
		DBName:               cfg.RDS.DBName,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", rdsCfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	return db, nil
}
