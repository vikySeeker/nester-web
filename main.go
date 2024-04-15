package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	t "github.com/vikySeeker/nester-web/tasks"
	u "github.com/vikySeeker/nester-web/user"
	"github.com/vikySeeker/nester-web/utils"
)

var (
	key   = []byte("NesterSessionSecret-Key-!23")
	store = sessions.NewCookieStore(key)
)

// var asset_dir = "/home/seeker/projects/nester/nester-web/asset"
var host_path = utils.GetWd()
var view_dir = host_path + "/view/"

func getFileContents(filename string) []byte {
	login_page, err := os.ReadFile(view_dir + filename)
	if err != nil {
		log.Print(err)
		return []byte("Internal Error")
	}
	return login_page
}

/* func renderTemplate(filename string) {
	tmpl := template.Must(template.ParseFiles(filename))
	tmpl.
} */

func isAuthenticated(r *http.Request) bool {
	session, _ := store.Get(r, "SESSION-ID")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return false
	}

	return true
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if isAuthenticated(r) {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	}

	if r.Method == "POST" {
		username, password := r.FormValue("username"), r.FormValue("password")
		log.Printf("username:%s and password:%s\n", username, password)
		status, err := u.LoginUser(username, password)
		if err != nil {
			log.Println(err)
		}

		if !status {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
		}

		log.Print(err)
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	}

	w.Write(getFileContents("login.html"))
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write(getFileContents("signup.html"))
	} else if r.Method == "POST" {
		var user u.User
		user.Uid = -1
		user.Email = r.FormValue("email")
		user.Password = r.FormValue("password")
		user.Username = r.FormValue("username")
		err := u.CreateUser(&user)
		if err != nil {
			http.Redirect(w, r, "/signup", http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if r.FormValue("action") == "delete" {
			taskid, _ := strconv.Atoi(r.FormValue("taskid"))
			log.Println(taskid)
			t.DeleteTask(taskid)
			http.Redirect(w, r, "/dashboard", http.StatusFound)
		} else {
			var task t.Tasks
			task.Taskid = -1
			task.Uid = 69
			task.Taskname = r.FormValue("taskname")
			task.Task_IP = r.FormValue("tip")
			task.Task_Domain = r.FormValue("tdomain")
			if task.Task_Domain == "" {
				task.Task_Domain = "nc"
			}
			if task.Task_IP == "" {
				task.Task_IP = "nc"
			}
			task.Created_at = time.Now().Format("2006-01-02 15:04:05")
			task.Completed_at = "nc"
			log.Println(task)
			err := t.AddTask(&task)
			log.Println(err)
		}

	} else {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./view/dashboard.html"))
	Data := t.GetTaskList()
	tmpl.Execute(w, Data)
	//w.Write(getFileContents("dashboard.html"))
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(getFileContents("profile.html"))
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	//logoutUser()
	http.Redirect(w, r, "/login", http.StatusOK)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}

func listtaskHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./view/tasklist.html"))
	Data := t.GetTaskList()
	tmpl.Execute(w, Data)
	//w.Write(getFileContents("tasklist.html"))
}

func main() {
	pwd, _ := os.Getwd()
	log.Print(pwd)
	t.GetTaskList()
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/task", taskHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/listtask", listtaskHandler)
	fs := http.FileServer(http.Dir("asset/"))
	http.Handle("/asset/", http.StripPrefix("/asset/", fs))
	log.Println("Running server...")
	log.Fatal(http.ListenAndServe(":80", nil))
}
