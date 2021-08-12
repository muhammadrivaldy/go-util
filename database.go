package util

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewMySQL is a function for set up db connection
func NewMySQL(user, password, url, schema string) (*sql.DB, error) {
	conn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, url, schema)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}

	return db, err
}

// NewGorm is a override sql from native sql to gorm
// driverName is a package of database your using
// level is a info log of query, value is 1 - 4 and default is 1
func NewGorm(db *sql.DB, driverName string, level int) (*gorm.DB, error) {
	var loglevel = logger.Silent
	if level == 1 {
		loglevel = logger.Silent
	} else if level == 2 {
		loglevel = logger.Error
	} else if level == 3 {
		loglevel = logger.Warn
	} else if level == 4 {
		loglevel = logger.Info
	}

	switch driverName {
	case "mysql":
		return gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{
			Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      loglevel,
				Colorful:      true,
			})})
	default:
		return nil, errors.New("driver name is not registered")
	}
}

// NewSqlx is a override sql from native sql to sqlx
func NewSqlx(db *sql.DB, driverName string) *sqlx.DB {
	sqlxCon := sqlx.NewDb(db, driverName)
	sqlxCon.SetConnMaxLifetime(time.Minute * 3)
	sqlxCon.SetMaxOpenConns(10)
	sqlxCon.SetMaxIdleConns(10)

	return sqlxCon
}
