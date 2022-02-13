package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

const (
	pubPath             = "../public"
	staticFilePath      = "/static/"
	viewFilePath        = pubPath + "/views/"
	port                = "8080"
	cookie_key_session  = "go_p_zo_webcli_cookie_key_session"
	sessionLifetimeDate = 1
)

func main() {
	// csrf認証キーの取得
	authKey := Cfg.CsrfAuthKey

	r := mux.NewRouter()
	r.HandleFunc("/signup", handleSignUp)
	r.HandleFunc("/signin", handleSignIn)
	r.HandleFunc("/signout", handleSignOut)
	r.HandleFunc("/mypage", handleMypage)
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(pubPath)))
	r.HandleFunc("/", handleHome)

	// セッション管理起動
	ctx, cancel := context.WithCancel(context.Background())
	Sm = NewSessionManager(ctx, cookie_key_session, sessionLifetimeDate)
	defer cancel()

	csrfMiddleware := csrf.Protect([]byte(authKey), csrf.Path("/"))(r)
	log.Printf("running on http://localhost:%s", port)
	http.ListenAndServe("localhost:"+port, csrfMiddleware)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	isLogin := false
	if _, err := Sm.GetSession(w, r); err == nil {
		isLogin = true
	}

	t, _ := template.ParseFiles(viewFilePath + "home.html")
	t.Execute(w, map[string]interface{}{
		"isLogin":        isLogin,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

// ユーザー登録
func handleSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// セッションが存在しない
		if _, err := Sm.GetSession(w, r); err != nil {
			t, _ := template.ParseFiles(viewFilePath + "signup.html")
			t.Execute(w, map[string]interface{}{
				csrf.TemplateTag: csrf.TemplateField(r),
			})
		} else {
			http.Redirect(w, r, "/mypage", http.StatusMovedPermanently)
			return
		}
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

			t, _ := template.ParseFiles(viewFilePath + "signup.html")
			t.Execute(w, map[string]interface{}{
				"message":        message,
				csrf.TemplateTag: csrf.TemplateField(r),
			})
			return
		}

		http.Redirect(w, r, "/mypage", http.StatusMovedPermanently)
	}
}

// ログイン
func handleSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// セッションが存在しない
		if _, err := Sm.GetSession(w, r); err != nil {
			t, _ := template.ParseFiles(viewFilePath + "signin.html")
			t.Execute(w, map[string]interface{}{
				csrf.TemplateTag: csrf.TemplateField(r),
			})
		} else {
			http.Redirect(w, r, "/mypage", http.StatusMovedPermanently)
			return
		}
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		email := r.Form.Get("email")
		password := r.Form.Get("password")

		res, err := RequestSignin(email, password)
		log.Printf("res: %v, err: %v", res, err)

		if err != nil || res.StatusCode != http.StatusCreated {
			message := "認証情報が正しくありません"
			if err != nil {
				message = err.Error()
			} else if res.Error != nil {
				message = res.Error.Message
			}

			t, _ := template.ParseFiles(viewFilePath + "signin.html")
			t.Execute(w, map[string]interface{}{
				"message":        message,
				csrf.TemplateTag: csrf.TemplateField(r),
			})
			return
		}

		sessionData := &SessionData{UserToken: *res.UserToken}
		log.Printf("sessionData: %v", sessionData)
		Sm.StartSession(w, r, sessionData)

		http.Redirect(w, r, "/mypage", http.StatusMovedPermanently)
	}
}

// ログアウト
func handleSignOut(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// セッションが存在する場合は削除する
		if sessionData, err := Sm.GetSession(w, r); err == nil {
			Sm.EndSession(w, r, sessionData.SessionId)
		}
	}

	t, _ := template.ParseFiles(viewFilePath + "home.html")
	t.Execute(w, nil)
}

// マイページ
func handleMypage(w http.ResponseWriter, r *http.Request) {
	s := NewService()

	// セッションが存在しない
	session, err := Sm.GetSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusMovedPermanently)
		return
	}

	if r.Method == http.MethodGet {
		res, err := s.GetMypage(&session.UserToken)
		log.Printf("res: %v, err: %v", res, err)

		message := ""
		if err != nil || res.StatusCode != http.StatusOK {
			message = "ユーザー情報が取得できません。再度ログインをお願いします。"
			if err != nil {
				message += err.Error()
			} else if res.Error != nil {
				message += res.Error.Message
			}

			t, _ := template.ParseFiles(viewFilePath + "signup.html")
			t.Execute(w, map[string]interface{}{
				"message":        message,
				csrf.TemplateTag: csrf.TemplateField(r),
			})
			return
		}

		funcMap := template.FuncMap{"TimeToSimple": gPresenter.TimeToSimple}
		t := template.New("mypage.html").Funcs(funcMap)
		t, _ = t.ParseFiles(viewFilePath + "mypage.html")
		log.Printf("session %v", session)
		t.Execute(w, map[string]interface{}{
			"message":        message,
			"model":          res,
			csrf.TemplateTag: csrf.TemplateField(r),
		})
	} else {
		_, err = s.PostNewZo(&session.UserToken, r.Form)
		if err != nil {
			ErrorPage(w, r, err.Error())
			return
		}

		http.Redirect(w, r, "/mypage", http.StatusMovedPermanently)
	}
}

// TODO
func ErrorPage(w http.ResponseWriter, r *http.Request, message string) {
	t, _ := template.ParseFiles(viewFilePath + "home.html")
	t.Execute(w, map[string]interface{}{
		"isLogin": true,
		"message": message,
	})
}
