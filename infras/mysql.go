package infras

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

// using pointer to reference value in the memory address so it will not create new variable
type Conn struct {
	Read  *sqlx.DB
	Write *sqlx.DB
}

func ProvideConn() Conn {
	return Conn{
		Read:  CreateReadConnection(),
		Write: CreateWriteConnection(),
	}
}

func CreateReadConnection() *sqlx.DB {
	return CreateNewConnection(
		"root",
		"root",
		"localhost",
		"33061",
		"bootcamp_auth",
		"Asia%2FBangkok", // %2F (URL encoding) is equal to /
	)
}

func CreateWriteConnection() *sqlx.DB {
	return CreateNewConnection(
		"root",
		"root",
		"localhost",
		"33061",
		"bootcamp_auth",
		"Asia%2FBangkok",
	)
}

func CreateNewConnection(username, pass, host, port, dbname, tz string) *sqlx.DB {
	connectionDetail := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=%s&parseTime=True", username, pass, host, port, dbname, tz)
	sqlConn, err := sqlx.Connect("mysql", connectionDetail)
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to connect")
	}

	// Set maximum number of open connections
	sqlConn.SetMaxOpenConns(100)

	// Set maximum number of idle connections
	sqlConn.SetMaxIdleConns(10)

	// Set maximum connection lifetime
	sqlConn.SetConnMaxLifetime(time.Hour)

	// Set maximum connection idle time
	sqlConn.SetConnMaxIdleTime(time.Minute)

	return sqlConn
}
