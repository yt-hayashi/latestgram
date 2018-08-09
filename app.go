package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

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
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/comment", comment)
	http.ListenAndServe(":8080", nil)
}

type post struct {
	PostID   int
	NameText string
	ImgPath  string
}

type contents []*post

//topページ
func top(w http.ResponseWriter, r *http.Request) {

	tmp := template.Must(template.ParseFiles("template/top.html.tpl"))

	rows, err := db.Query("SELECT posts.id, name, img_name FROM posts INNER JOIN users ON posts.user_id=users.id ORDER BY posts.id  DESC limit 50")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var posts contents

	for rows.Next() {
		var (
			postID   int
			userName string
			imgName  string
		)
		if err := rows.Scan(&postID, &userName, &imgName); err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		posts = append(posts, &post{postID, userName, imgName})
	}

	//session 読み出し
	session, err := store.Get(r, "user-session")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Print(session.Values["userID"])

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

	if r.Method == http.MethodPost {
		//リクエストの解析
		r.ParseForm()
		inputUserName := fmt.Sprint(r.Form.Get("username"))
		inputPassword := fmt.Sprint(r.Form.Get("password"))
		if (inputUserName == "") || (inputPassword == "") {
			message := "Input Form!"
			w.WriteHeader(http.StatusNotAcceptable)
			if err := tmp.ExecuteTemplate(w, "login.html.tpl", message); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		// login時の処理

		//DBから読み出し
		rows, err := db.Query("SELECT id, name, password FROM users WHERE name=?", inputUserName)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		readUserName := ""
		readPassword := ""
		var id int
		for rows.Next() {
			if err := rows.Scan(&id, &readUserName, &readPassword); err != nil {
				fmt.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		//password確認
		if err := bcrypt.CompareHashAndPassword([]byte(readPassword), []byte(inputPassword)); err != nil {
			fmt.Println(err.Error())
			message := "Something is wrong."
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
		session.Values["userID"] = id
		session.Values["userName"] = readUserName
		// Save it before we write to the response/return from the handler.
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	message := ""
	if err := tmp.ExecuteTemplate(w, "login.html.tpl", message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles("template/upload.html.tpl"))
	//session 読み出し
	session, err := store.Get(r, "user-session")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID := session.Values["userID"]
	userName := fmt.Sprint(session.Values["userName"])

	if userID == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	//画像の受け取り
	if r.Method == http.MethodPost {
		imgFile, _, err := r.FormFile("upload")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer imgFile.Close()

		//ユーザー名のディレクトリを作成
		dirPath := "./img/" + userName
		if err := os.Mkdir(dirPath, 0777); err != nil {
			fmt.Println(err.Error())
		}

		//ファイルの作成
		nowTime := fmt.Sprint(time.Now().Unix())
		imgName := nowTime + ".jpg"
		imgPath := "./img/" + userName + "/" + imgName

		f, err := os.Create(imgPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		io.Copy(f, imgFile)

		//DBに書き込み
		if _, err := db.Exec(`
		INSERT INTO posts(user_id, img_name) VALUES(?, ?)`, userID, imgName); err != nil {
			fmt.Println("Error! Post didn't add.", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/top", http.StatusFound)
	}

	if err := tmp.ExecuteTemplate(w, "upload.html.tpl", userName); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func comment(w http.ResponseWriter, r *http.Request) {
	//session 読み出し
	session, err := store.Get(r, "user-session")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID := session.Values["userID"]
	//ログインしてない場合はリダイレクト
	if userID == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	//URLのパラメーターの解析
	param, ok := r.URL.Query()["id"]
	if !ok || len(param[0]) < 1 {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	postID := param[0]

	if r.Method == http.MethodPost {
		r.ParseForm()
		commentBody := fmt.Sprint(r.Form.Get("comment_text"))
		fmt.Print(postID)
		//DBに書き込み
		if _, err := db.Exec(`
		INSERT INTO comments(post_id, user_id, comment_body) VALUES(?, ?, ?)`, postID, userID, commentBody); err != nil {
			fmt.Println("Error! Post didn't add.", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/top", http.StatusFound)
	}
}
