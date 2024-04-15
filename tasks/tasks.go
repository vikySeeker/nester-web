package tasks

import (
	"errors"
	"log"

	"github.com/vikySeeker/nester-web/db"
)

type Tasks struct {
	Taskid       int
	Uid          int
	Taskname     string
	Created_at   string
	Completed_at string
	Task_IP      string
	Task_Domain  string
}

type TaskList struct {
	Err_msg       bool
	List          []Tasks
	Taskcount     int
	Taskcompleted int
	Bugsfound     int
}

func isTaskEmpty(task *Tasks) bool {
	if task.Taskname == "" || task.Task_IP == "" || task.Task_Domain == "" {
		return true
	}
	return false
}

func GetTaskList() TaskList {
	dbconn, err := db.GetConn()
	if err != nil {
		return TaskList{
			Err_msg: true,
			List:    nil,
		}
	}
	log.Println("querying rows...")
	rows, err := dbconn.Query("select * from tasks;")
	if err != nil {
		return TaskList{
			Err_msg: true,
			List:    nil,
		}
	}

	log.Println("Scanning values...")
	var tasklist TaskList
	tasklist.Err_msg = false
	tasklist.Taskcount = 0
	tasklist.Taskcompleted = 0
	tasklist.Bugsfound = 0
	for rows.Next() {
		log.Println("getting values...")
		var t Tasks
		var tid int
		var uid int
		var tname string
		var cat string
		var cot string
		var tip string
		var td string
		_ = rows.Scan(&tid, &uid, &tname, &cat, &cot, &tip, &td)
		t.Taskid = tid
		t.Uid = uid
		t.Taskname = tname
		t.Created_at = cat
		t.Completed_at = cot
		t.Task_IP = tip
		t.Task_Domain = td
		tasklist.List = append(tasklist.List, t)
	}
	log.Println("Printing task lists...")
	for _, t := range tasklist.List {
		log.Printf("taskid=%d, uid=%d, taskname=%s, cat=%s, coat=%s, taskip=%s,taskdomain=%s\n", t.Taskid, t.Uid, t.Taskname, t.Created_at, t.Completed_at, t.Task_IP, t.Task_Domain)
	}
	count := len(tasklist.List)
	tasklist.Taskcount = count
	if count == 0 {
		tasklist.Err_msg = true
	}

	log.Println(tasklist)
	return tasklist
}

func AddTask(task *Tasks) error {
	if isTaskEmpty(task) {
		return errors.New("empty task details")
	}
	dbconn, err := db.GetConn()
	if err != nil {
		return err
	}

	stmt, err := dbconn.Prepare("insert into tasks(uid, taskname, created_at, completed_at, task_ip, task_domain) values(?, ?, ?, ?, ?, ?);")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(task.Uid, task.Taskname, task.Created_at, task.Completed_at, task.Task_IP, task.Task_Domain)

	return err
}

func DeleteTask(taskid int) error {
	dbconn, err := db.GetConn()
	if err != nil {
		return err
	}

	log.Println("before exec...")
	rowsaffected, err := dbconn.Exec("DELETE FROM tasks WHERE taskid = ?", taskid)
	log.Println("after exec...")

	log.Println(rowsaffected.RowsAffected())
	if err != nil {
		return err
	}

	return err
}
