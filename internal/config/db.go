package config

import "database/sql"

func InitDb(cfg *DbConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetConnMaxLifetime(cfg.MaxLifetimeConn)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
