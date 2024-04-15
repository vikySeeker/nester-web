package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbconn, err := sql.Open("sqlite3", "./test.db")

	if err != nil {
		fmt.Println(err)
		return
	}

	rows, err := dbconn.Query("select * from test;")
	if err != nil {
		fmt.Println(err)
		return
	}
	var taskid int
	var taskip string
	var uid int
	var taskname string
	var cat string
	var coat string
	var taskdomain string
	for rows.Next() {
		rows.Scan(&taskid, &uid, &taskname, &cat, &coat, &taskip, &taskdomain)
		fmt.Printf("taskid=%d, uid=%d, taskname=%s, cat=%s, coat=%s, taskip=%s,taskdomain=%s\n", taskid, uid, taskname, cat, coat, taskip, taskdomain)
	}

}
