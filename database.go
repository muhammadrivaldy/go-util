package goutil

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type logLevel int

const (
	LoggerSilent = logLevel(logger.Silent)
	LoggerError  = logLevel(logger.Error)
	LoggerWarn   = logLevel(logger.Warn)
	LoggerInfo   = logLevel(logger.Info)
)

// NewMySQL is a function for set up db mysql connection
func NewMySQL(user, password, url, schema string, parameters []string) (*sql.DB, error) {

	parameter := strings.Join(parameters, "&")
	conn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", user, password, url, schema, parameter)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}

	return db, err

}

// NewPostgreSQL is a function for set up db postgresql connection
func NewPostgreSQL(user, password, host, schema string, sslMode null.String, port int) (*sql.DB, error) {

	if !sslMode.Valid {
		sslMode.SetValid("disable")
	}

	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, schema, sslMode.String)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return db, err

}

// NewGorm is a override sql from native sql to gorm.
// driverName is a package of database your using, currently we have only mysql package.
// level is a info log of query, value is 1 - 4 (Silent, Error, Warn & Info)
func NewGorm(db *sql.DB, driverName string, level logLevel) (*gorm.DB, error) {

	gormConfig := &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.LogLevel(level),
			Colorful:      true,
		})}

	switch driverName {
	case "mysql":
		return gorm.Open(mysql.New(mysql.Config{Conn: db}), gormConfig)
	case "postgres":
		return gorm.Open(postgres.New(postgres.Config{Conn: db}), gormConfig)
	}

	return nil, errors.New("driver name is not registered")

}

// NewSqlx is a override sql from native sql to sqlx
func NewSqlx(db *sql.DB, driverName string) *sqlx.DB {

	sqlxCon := sqlx.NewDb(db, driverName)
	sqlxCon.SetConnMaxLifetime(time.Minute * 3)
	sqlxCon.SetMaxOpenConns(10)
	sqlxCon.SetMaxIdleConns(10)

	return sqlxCon
}
