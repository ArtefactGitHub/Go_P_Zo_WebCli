package main

import (
	"html/template"
	"log"
	"net/http"
)

type User struct {
	Name     string
	Email    string
	Password string
}

const pubPath = "../public"
const port = "8080"

var user User = User{}

func main() {
	http.HandleFunc("/signup", handleSignUp)
	http.HandleFunc("/signin", handleSignIn)
	http.HandleFunc("/user", handleUser)
	http.HandleFunc("/", handleHome)

	log.Printf("running on http://localhost:%s", port)
	log.Print(http.ListenAndServe("localhost:"+port, nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(pubPath + "/home/index.html")
	t.Execute(w, nil)
}

func handleSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, _ := template.ParseFiles(pubPath + "/auth/signup.html")
		t.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		r.ParseForm()

		username := r.Form.Get("username")
		email := r.Form.Get("email")
		password := r.Form.Get("password")
		log.Printf("%s, %s, %s", username, email, password)

		user = User{Name: username, Email: email, Password: password}

		http.Redirect(w, r, "/user", http.StatusMovedPermanently)
	}
}

func handleSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, _ := template.ParseFiles(pubPath + "/auth/signin.html")
		t.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		r.ParseForm()

		email := r.Form.Get("email")
		password := r.Form.Get("password")
		log.Printf("%s, %s", email, password)

		if email == user.Email && password == user.Password {
			http.Redirect(w, r, "/user", http.StatusMovedPermanently)
		} else {
			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		}
	}
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(pubPath + "/user/mypage.html")
	t.Execute(w, user)
}
