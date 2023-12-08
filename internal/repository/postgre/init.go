package postgre

import (
	"fmt"
	"time"

	"github.com/tokopedia/sqlt"
)

var (
	DBChat   *sqlt.DB
	DBMaster string = "user=secure_random_username password=secure_random_password dbname=chat host=postgres port=5432 sslmode=disable"
)

func InitPostgreSqltDB(connMaster string, connSlave string) error {
	var (
		dbURL string
		err   error
	)

	if connSlave == "" {
		dbURL = connMaster
	} else {
		dbURL = fmt.Sprintf("%s;%s", connMaster, connSlave)
	}

	DBChat, err = sqlt.Open("postgres", dbURL)
	if err != nil {
		return err
	}

	if err := DBChat.Ping(); err != nil {
		return err
	}
	DBChat.SetConnMaxLifetime(9 * time.Second)

	return nil
}
