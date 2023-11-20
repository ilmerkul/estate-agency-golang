package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"gilab.com/estate-agency-api/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

func NewClient(storageCfg *config.StorageConfig) (*sql.DB, error) {

	DB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%s)/%s", storageCfg.ConfigDSN.User, storageCfg.ConfigDSN.Password, storageCfg.ConfigDSN.Protocol, storageCfg.ConfigDSN.Host, storageCfg.ConfigDSN.Port, storageCfg.ConfigDSN.NameDB))

	if err != nil {
		return nil, err
	}

	DB.SetConnMaxLifetime(time.Duration(storageCfg.ConnMaxLifetime))
	DB.SetMaxIdleConns(storageCfg.MaxIdleConns)
	DB.SetMaxOpenConns(storageCfg.MaxOpenConns)

	if err := DB.Ping(); err != nil {
		return nil, err
	}

	/*
		stmt, err := DB.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
		`)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		_, err = stmt.Exec()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	*/

	return DB, nil
}
