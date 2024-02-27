package main

import (
	"log"
	"net/http"
	"os"
)

var asset_dir = "/home/seeker/projects/NeSTer/web/asset"
var view_dir = "/home/seeker/projects/NeSTer/web/view/"

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
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write(getFileContents("signup.html"))
	} else if r.Method == "POST" {
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
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	fs := http.FileServer(http.Dir("asset/"))
	http.Handle("/asset/", http.StripPrefix("/asset/", fs))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
