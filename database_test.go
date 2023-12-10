package goutil_test

import (
	"testing"

	_ "github.com/lib/pq"
	goutil "github.com/muhammadrivaldy/go-util"
	"gopkg.in/guregu/null.v4"
)

// go test -v -run=TestPostgreSQLConnectionWithGorm
func TestPostgreSQLConnectionWithGorm(t *testing.T) {

	user := "postgres"
	password := "postgres"
	host := "localhost"
	schema := "database"
	port := 5432

	sqlDB, err := goutil.NewPostgreSQL(user, password, host, schema, null.String{}, port)
	if err != nil {
		panic(err)
	}

	_, err = goutil.NewGorm(sqlDB, "postgres", goutil.LoggerInfo)
	if err != nil {
		panic(err)
	}

}
