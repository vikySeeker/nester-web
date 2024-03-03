package main

import (
	"log"
	"net/http"
	"os"

	"github.com/vikySeeker/nester-web/db"
)

// var asset_dir = "/home/seeker/projects/nester/nester-web/asset"
var view_dir = "/home/seeker/projects/nester/nester-web/view/"

func getFileContents(filename string) []byte {
	login_page, err := os.ReadFile(view_dir + filename)
	if err != nil {
		log.Print(err)
		return []byte("Internal Error")
	}
	return login_page
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write(getFileContents("login.html"))
	} else if r.Method == "POST" {
		username, password := r.FormValue("username"), r.FormValue("password")
		log.Printf("username:%s and password:%s\n", username, password)
		status, err := db.LoginUser(username, password)
		if err != nil {
			log.Println(err)
		}

		if !status {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
		}

		log.Print(err)
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write(getFileContents("signup.html"))
	} else if r.Method == "POST" {
		var user db.User
		user.Uid = -1
		user.Email = r.FormValue("email")
		user.Password = r.FormValue("password")
		user.Username = r.FormValue("username")
		err := db.CreateUser(&user)
		if err != nil {
			http.Redirect(w, r, "/signup", http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(getFileContents("dashboard.html"))
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(getFileContents("profile.html"))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}

func main() {
	pwd, _ := os.Getwd()
	log.Print(pwd)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	fs := http.FileServer(http.Dir("asset/"))
	http.Handle("/asset/", http.StripPrefix("/asset/", fs))
	log.Fatal(http.ListenAndServe(":80", nil))
}
