package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

type ViewArgs map[string]interface{}

func main() {
	// csrf認証キーの取得
	authKey := Cfg.CsrfAuthKey

	r := mux.NewRouter()
	r.HandleFunc(SignUpPath, handleSignUp)
	r.HandleFunc(SignInPath, handleSignIn)
	r.HandleFunc(SignOutPath, handleSignOut)
	r.HandleFunc(MyPageZosPath, handleMypageZos)
	r.HandleFunc(MyPageUserPath, handleMypageUser)
	r.PathPrefix(externalFilePath).Handler(http.FileServer(http.Dir(pubPath)))
	r.PathPrefix(internalFilePath).Handler(http.FileServer(http.Dir(pubPath)))
	r.HandleFunc(TopPath, handleHome)

	// セッション管理起動
	ctx, cancel := context.WithCancel(context.Background())
	Sm = NewSessionManager(ctx, cookie_key_session, sessionLifetimeDate)
	defer cancel()

	csrfMiddleware := csrf.Protect([]byte(authKey), csrf.Path("/"))(r)
	log.Printf("running on http://localhost:%s", port)
	http.ListenAndServe("localhost:"+port, csrfMiddleware)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, SignInPath, http.StatusMovedPermanently)
}

// ユーザー登録
func handleSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// セッションが存在しない
		if _, err := Sm.GetSession(w, r); err != nil {
			ExecuteTemplate(w, r, "signup", ViewArgs{csrf.TemplateTag: csrf.TemplateField(r)})
		} else {
			http.Redirect(w, r, MyPageZosPath, http.StatusMovedPermanently)
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

			ExecuteTemplate(w, r, "signup", ViewArgs{"message": message, csrf.TemplateTag: csrf.TemplateField(r)})
			return
		} else {
			for _, v := range defaultCategories {
				_, err := RequestPostUserCategory(res.UserToken, NewRequestCategory(v.Name, v.ColorId))
				if err != nil {
					ExecuteTemplate(w, r, "signup", ViewArgs{"message": err.Error(), csrf.TemplateTag: csrf.TemplateField(r)})
					return
				}
			}
		}

		http.Redirect(w, r, SignInPath, http.StatusMovedPermanently)
	}
}

// ログイン
func handleSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// セッションが存在しない
		if _, err := Sm.GetSession(w, r); err != nil {
			ExecuteTemplate(w, r, "signin", ViewArgs{csrf.TemplateTag: csrf.TemplateField(r)})
		} else {
			http.Redirect(w, r, MyPageZosPath, http.StatusMovedPermanently)
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

			ExecuteTemplate(w, r, "signin", ViewArgs{"message": message, csrf.TemplateTag: csrf.TemplateField(r)})
			return
		}

		sessionData := &SessionData{UserToken: *res.UserToken}
		log.Printf("sessionData: %v", sessionData)
		Sm.StartSession(w, r, sessionData)

		http.Redirect(w, r, MyPageZosPath, http.StatusMovedPermanently)
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

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// ユーザー情報
func handleMypageUser(w http.ResponseWriter, r *http.Request) {
	s := NewService()

	// セッションが存在しない
	session, err := Sm.GetSession(w, r)
	if err != nil {
		http.Redirect(w, r, SignInPath, http.StatusMovedPermanently)
		return
	}

	if r.Method == http.MethodGet {
		res, err := s.GetMypageUser(&session.UserToken)
		log.Printf("res: %v, err: %v", res, err)

		message := ""
		if err != nil || res.StatusCode != http.StatusOK {
			message = "ユーザー情報が取得できません。再度ログインをお願いします。"
			if err != nil {
				message += err.Error()
			} else if res.Error != nil {
				message += res.Error.Message
			}

			ExecuteTemplate(w, r, "signin", ViewArgs{"message": message, csrf.TemplateTag: csrf.TemplateField(r)})
			return
		}

		ExecuteTemplate(w, r, "mypage_user", ViewArgs{
			"message":        message,
			"model":          res,
			csrf.TemplateTag: csrf.TemplateField(r),
		})
	} else {
		err = s.PostMypageUser(&session.UserToken, r.Form)
		if err != nil {
			ErrorPage(w, r, err.Error())
			return
		}

		http.Redirect(w, r, MyPageUserPath, http.StatusMovedPermanently)
	}
}

// マイページ
func handleMypageZos(w http.ResponseWriter, r *http.Request) {
	s := NewService()

	// セッションが存在しない
	session, err := Sm.GetSession(w, r)
	if err != nil {
		http.Redirect(w, r, SignInPath, http.StatusMovedPermanently)
		return
	}

	if r.Method == http.MethodGet {
		res, err := s.GetMypageZos(&session.UserToken)
		log.Printf("[Get]handleMypageZos res: %v, err: %v", res, err)

		message := ""
		if err != nil || res.StatusCode != http.StatusOK {
			message = "ユーザー情報が取得できません。再度ログインをお願いします。"
			if err != nil {
				message += err.Error()
			} else if res.Error != nil {
				message += res.Error.Message
			}

			ExecuteTemplate(w, r, "signup", ViewArgs{"message": message, csrf.TemplateTag: csrf.TemplateField(r)})
			return
		}

		ExecuteTemplateWithFunc(w, r, "mypage_zos", ViewArgs{
			"message":        message,
			"model":          res,
			csrf.TemplateTag: csrf.TemplateField(r),
		}, template.FuncMap{"TimeToSimple": gPresenter.TimeToSimple})
	} else {
		_, err = s.PostNewZo(&session.UserToken, r.Form)
		if err != nil {
			ErrorPage(w, r, err.Error())
			return
		}

		http.Redirect(w, r, MyPageZosPath, http.StatusMovedPermanently)
	}
}

// TODO
func ErrorPage(w http.ResponseWriter, r *http.Request, message string) {
	ExecuteTemplate(w, r, "home", ViewArgs{"message": message})
}

func IsLogin(w http.ResponseWriter, r *http.Request) bool {
	result := false
	if _, err := Sm.GetSession(w, r); err == nil {
		result = true
	}
	return result
}

func ExecuteTemplate(w http.ResponseWriter, r *http.Request, viewName string, args ViewArgs) {
	t, _ := template.ParseFiles(layoutFilePath, viewFilePath+viewName+".html")
	t.ExecuteTemplate(w, layoutName, args.Add("isLogin", IsLogin(w, r)))
}

func ExecuteTemplateWithFunc(w http.ResponseWriter, r *http.Request, viewName string, args ViewArgs, funcMap template.FuncMap) {
	t, _ := template.New(viewName+".html").
		Funcs(template.FuncMap{"TimeToSimple": gPresenter.TimeToSimple}).
		ParseFiles(layoutFilePath, viewFilePath+viewName+".html")
	t.ExecuteTemplate(w, layoutName, args.Add("isLogin", IsLogin(w, r)))
}

func (a *ViewArgs) Add(key string, value interface{}) *ViewArgs {
	(*a)[key] = value
	return a
}
