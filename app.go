package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

var (
	db    *sql.DB
	store = sessions.NewCookieStore([]byte("something-very-secret"))
)

func main() {
	_db, err := sql.Open("mysql", "root:mysql@tcp(127.0.0.1:33306)/test")
	if err != nil {
		fmt.Println("DB Error! --> ", err.Error())
		os.Exit(1)
	}
	db = _db

	defer db.Close()

	http.HandleFunc("/", top)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.ListenAndServe(":8080", nil)
}

type post struct {
	NameText string
	ImgPath  string
}

type contents []*post

//topページ
func top(w http.ResponseWriter, r *http.Request) {

	tmp := template.Must(template.ParseFiles("template/top.html.tpl"))

	rows, err := db.Query("SELECT name, img_name FROM posts INNER JOIN users ON posts.user_id=users.id ORDER BY posts.created_at limit 50")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var posts contents

	for rows.Next() {
		var userName string
		var imgName string
		if err := rows.Scan(&userName, &imgName); err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		posts = append(posts, &post{userName, imgName})
	}

	//session 読み出し
	session, err := store.Get(r, "user-session")
	fmt.Println(session.Values["userName"])

	if err := tmp.ExecuteTemplate(w, "top.html.tpl", posts); err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

//signupページ
func signup(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles("template/signup.html.tpl"))

	message := ""

	if r.Method == http.MethodPost {
		//レスポンスの解析
		r.ParseForm()
		userName := fmt.Sprint(r.Form.Get("username"))
		password := fmt.Sprint(r.Form.Get("password"))
		if (userName == "") || (password == "") {
			message = "Input Form!"
			w.WriteHeader(http.StatusNotAcceptable)
			if err := tmp.ExecuteTemplate(w, "signup.html.tpl", message); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		// signup時の処理
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
			fmt.Println("Hash Error!")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//DBに追加
		if _, err := db.Exec(`
		INSERT INTO users(name, password) VALUES(?, ?)`, userName, hash); err != nil {
			fmt.Println("Error! User didn't add.", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := tmp.ExecuteTemplate(w, "signup.html.tpl", message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

//loginページ
func login(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles("template/login.html.tpl"))
	message := ""

	if r.Method == http.MethodPost {
		//レスポンスの解析
		r.ParseForm()
		userName := fmt.Sprint(r.Form.Get("username"))
		password := fmt.Sprint(r.Form.Get("password"))
		if (userName == "") || (password == "") {
			message = "Input Form!"
			w.WriteHeader(http.StatusNotAcceptable)
			if err := tmp.ExecuteTemplate(w, "login.html.tpl", message); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		// login時の処理

		//DBから読み出し
		rows, err := db.Query("SELECT name, password FROM users WHERE name=?", userName)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_userName := ""
		_password := ""
		for rows.Next() {
			if err := rows.Scan(&_userName, &_password); err != nil {
				fmt.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		//password確認
		fmt.Println(_userName, password, _password)
		if err := bcrypt.CompareHashAndPassword([]byte(_password), []byte(password)); err != nil {
			fmt.Println(err.Error())
			message = "Something is wrong."
			w.WriteHeader(http.StatusNotAcceptable)
			if err := tmp.ExecuteTemplate(w, "login.html.tpl", message); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		fmt.Println("Success!")
		//session書き込み
		session, err := store.Get(r, "user-session")
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		session.Values["userName"] = userName
		// Save it before we write to the response/return from the handler.
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := tmp.ExecuteTemplate(w, "login.html.tpl", message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
