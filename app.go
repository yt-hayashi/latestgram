package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
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
	_db, err := sql.Open("mysql", "root:mysql@tcp(127.0.0.1:3306)/test")
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
	http.HandleFunc("/commentdel", commentdel)
	http.HandleFunc("/logout", logout)
	http.Handle("/img/", http.FileServer(http.Dir("./")))
	http.Handle("/css/", http.FileServer(http.Dir("./")))
	http.ListenAndServe(":8080", nil)
}

type post struct {
	PostID   int
	NameText string
	ImgPath  string
	Comments []string
}

type contents []*post

//topページ
func top(w http.ResponseWriter, r *http.Request) {

	tmp := template.Must(template.ParseFiles("template/top.html.tpl"))

	rows, err := db.Query("SELECT posts.id, name, img_name FROM posts INNER JOIN users ON posts.user_id=users.id ORDER BY posts.id DESC limit 50")
	if err != nil {
		log.Println("querry request", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var posts contents

	for rows.Next() {
		var (
			postID   int
			userName string
			imgName  string
			comments []string
		)
		if err := rows.Scan(&postID, &userName, &imgName); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//コメント読み込み
		commentsRows, err := db.Query("SELECT users.name ,comment_body FROM comments LEFT JOIN posts ON comments.post_id=posts.id LEFT JOIN users ON comments.user_id=users.id WHERE comments.post_id =?", postID)
		if err != nil {
			log.Println("read comment", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//このループが何故か複数回される？？
		for commentsRows.Next() {
			var commenterName string
			var comment string
			if err := commentsRows.Scan(&commenterName, &comment); err != nil {
				fmt.Println("comment db Err!!!")
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			//fmt.Println(comment)
			commentStr := commenterName + ": " + comment
			comments = append(comments, commentStr)
		}

		posts = append(posts, &post{postID, userName, imgName, comments})
	}

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
		//ユーザー名のディレクトリを作成
		dirPath := "./img/" + userName
		if err := os.Mkdir(dirPath, 0777); err != nil {
			fmt.Println(err.Error())
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
	userID, ok := session.Values["userID"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	userName := fmt.Sprint(session.Values["userName"])

	//画像の受け取り
	if r.Method == http.MethodPost {
		imgFile, header, err := r.FormFile("upload")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if header.Size == 0 {
			if err := tmp.ExecuteTemplate(w, "upload.html.tpl", userName); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		defer imgFile.Close()

		//ファイルの作成
		imgPath := fmt.Sprintf("./img/%s/%d.jpg", userName, time.Now().Unix())

		f, err := os.Create(imgPath)
		defer f.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		io.Copy(f, imgFile)

		//DBに書き込み
		if _, err := db.Exec(`
		INSERT INTO posts(user_id, img_name) VALUES(?, ?)`, userID, imgPath); err != nil {
			fmt.Println("Error! Post didn't add.", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
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
	userID, ok := session.Values["userID"]
	//ログインしてない場合はリダイレクト
	if ok == false {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	//URLのパラメーターの解析
	param, ok := r.URL.Query()["id"]
	if !ok || param[0] == "" {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	postID := param[0]

	commnenNum := 0
	if err := db.QueryRow(`SELECT COUNT(*) FROM comments INNER JOIN posts ON comments.post_id=posts.id WHERE comments.post_id =?`, postID).Scan(&commnenNum); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if commnenNum == 5 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

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
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

//logoutの処理
func logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	session.Values["userID"] = ""
	session.Values["userName"] = ""
	// Save it before we write to the response/return from the handler.
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func commentdel(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, ok := session.Values["userID"]
	//ログインしてない場合はリダイレクト
	if ok == false {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	//URLのパラメーターの解析
	param, ok := r.URL.Query()["id"]
	if !ok || param[0] == "" {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	commentID := param[0]

	//if r.Method == http.MethodPost {
	r.ParseForm()

	fmt.Print(commentID)
	//DBに書き込み
	if _, err := db.Exec(`DELETE FROM comments WHERE id =?`, commentID); err != nil {
		fmt.Println("Error! Comment didn't delete.", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
	//}
}
