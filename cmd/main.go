package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

const pubPath = "../public"
const port = "8080"

var user User = User{}

func main() {
	// csrf認証キーの取得
	authKey := Cfg.CsrfAuthKey

	r := mux.NewRouter()
	r.HandleFunc("/signup", handleSignUp)
	r.HandleFunc("/signin", handleSignIn)
	r.HandleFunc("/user", handleUser)
	r.HandleFunc("/", handleHome)

	csrfMiddleware := csrf.Protect([]byte(authKey), csrf.Path("/"))(r)
	log.Printf("running on http://localhost:%s", port)
	http.ListenAndServe("localhost:"+port, csrfMiddleware)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(pubPath + "/home/index.html")
	t.Execute(w, nil)
}

// ユーザー登録
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

		res, err := RequestSignup(username, "", email, password)
		log.Printf("res: %v, err: %v", res, err)

		if err != nil || res.StatusCode != http.StatusCreated {
			message := ""
			if err != nil {
				message = err.Error()
			} else {
				message = res.Error.Message
			}

			t, _ := template.ParseFiles(pubPath + "/auth/signup.html")
			t.Execute(w, map[string]interface{}{
				"message":        message,
				csrf.TemplateTag: csrf.TemplateField(r),
			})
			return
		}

		user.FamilyName = res.User.FamilyName
		user.Email = res.User.Email
		http.Redirect(w, r, "/user", http.StatusMovedPermanently)
	}
}

// ログイン
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

		res, err := RequestSignin(email, password)
		log.Printf("res: %v, err: %v", res, err)

		if err != nil || res.StatusCode != http.StatusOK {
			message := "認証情報が正しくありません"
			if err != nil {
				message = err.Error()
			} else if res.Error != nil {
				message = res.Error.Message
			}

			t, _ := template.ParseFiles(pubPath + "/auth/signin.html")
			t.Execute(w, map[string]interface{}{
				"message":        message,
				csrf.TemplateTag: csrf.TemplateField(r),
			})
			return
		}

		user.FamilyName = res.FamilyName
		user.Email = res.Email
		http.Redirect(w, r, "/user", http.StatusMovedPermanently)
	}
}

// マイページ
func handleUser(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(pubPath + "/user/mypage.html")
	log.Printf("user %v", user)
	t.Execute(w, user)
}
