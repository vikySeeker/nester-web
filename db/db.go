package db

import (
	"database/sql"

	bc "golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Uid      int
	Username string
	Email    string
	Password string
	Passhash string
}

var db_file_path = "/home/seeker/projects/nester/nester-web/db/nester.db"
var dbconn *sql.DB

func getConn() error {
	if dbconn != nil {
		return nil
	}

	db, err := sql.Open("sqlite3", db_file_path)
	dbconn = db

	if err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bc.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func validatePassword(password string, hash string) bool {
	err := bc.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(newuser *User) error {
	err := getConn()
	if err != nil {
		return err
	}

	stmt, err := dbconn.Prepare("insert into user(username, email, password) values(?, ?, ?)")
	if err != nil {
		return err
	}
	hash, err := hashPassword(newuser.Password)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(newuser.Username, newuser.Email, hash)

	if err != nil {
		return nil
	}
	return nil
}

func LoginUser(username string, password string) (bool, error) {

	err := getConn()

	if err != nil {
		return false, err
	}

	stmt, err := dbconn.Prepare("Select * from user where username=?")
	if err != nil {
		return false, err
	}

	var user User

	err = stmt.QueryRow(username).Scan(&user.Uid, &user.Username, &user.Email, &user.Passhash)

	if err != nil {
		return false, err
	}

	if validatePassword(password, user.Passhash) {
		return true, nil
	}

	return false, nil
}
