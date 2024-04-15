package db

import (
	"database/sql"
	"log"

	"github.com/vikySeeker/nester-web/utils"

	_ "github.com/mattn/go-sqlite3"
)

var host_path = utils.GetWd()

var db_file_path = host_path + "/db/nester.db"
var dbconn *sql.DB

func GetConn() (*sql.DB, error) {
	log.Println(host_path)
	if dbconn != nil {
		return dbconn, nil
	}

	dbconn, err := sql.Open("sqlite3", db_file_path)

	if err != nil {
		return nil, err
	}

	return dbconn, nil
}
