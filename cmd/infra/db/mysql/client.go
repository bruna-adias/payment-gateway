package mysql

import (
	"database/sql"
	"fmt"
	"payment-gateway/cmd/infra"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLClient(cfg *infra.Configuration) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)

	var db *sql.DB
	var err error
	for i := 0; i < 15; i++ {
		db, err = sql.Open("mysql", dsn)

		if err == nil {
			if err := db.Ping(); err == nil {
				break
			}

			db.Close()
		}

		fmt.Printf("failed to open database: %v\n", err)

		fmt.Printf("waiting for 5 seconds before retrying...\n")
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}
