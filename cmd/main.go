package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
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
	// csrf認証キーの取得
	authKey := os.Getenv("Go_P_Zo_WebCli_CsrfAuthKey")

	r := mux.NewRouter()
	r.HandleFunc("/signup", handleSignUp)
	r.HandleFunc("/signin", handleSignIn)
	r.HandleFunc("/user", handleUser)
	r.HandleFunc("/", handleHome)

	csrfMiddleware := csrf.Protect([]byte(authKey), csrf.Path("/"))(r)
	http.ListenAndServe("localhost:"+port, csrfMiddleware)
	log.Printf("running on http://localhost:%s", port)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(pubPath + "/home/index.html")
	t.Execute(w, nil)
}

func handleSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, _ := template.ParseFiles(pubPath + "/auth/signup.html")
		t.Execute(w, map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(r),
		})
	} else if r.Method == http.MethodPost {
		r.ParseForm()

		username := r.Form.Get("username")
		email := r.Form.Get("email")
		password := r.Form.Get("password")

		user = User{Name: username, Email: email, Password: password}

		http.Redirect(w, r, "/user", http.StatusMovedPermanently)
	}
}

func handleSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, _ := template.ParseFiles(pubPath + "/auth/signin.html")
		t.Execute(w, map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(r),
		})
	} else if r.Method == http.MethodPost {
		r.ParseForm()

		email := r.Form.Get("email")
		password := r.Form.Get("password")

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
