package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	_ "golang.org/x/crypto/bcrypt"
)

var (
	db *sql.DB
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

	if err := tmp.ExecuteTemplate(w, "top.html.tpl", posts); err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

//signupページ
func signup(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles("template/signup.html.tpl"))

	message := "Please Input"

	fmt.Println("method:", r.Method)
	//レスポンスの解析
	r.ParseForm()
	userName := fmt.Sprint(r.Form["username"])
	password := fmt.Sprint(r.Form["password"])

	fmt.Println("username:", userName)
	fmt.Println("password:", password)

	if r.Method == http.MethodPost {
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
		_, _err := db.Exec(`
		INSERT INTO users(name, password) VALUES(?, ?)`, userName, hash)

		if _err != nil {
			fmt.Println("Error! User didn't add.", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		message = "SignUp Successed!"
	}

	// getのrender処理
	fmt.Println("method:", r.Method)

	if err := tmp.ExecuteTemplate(w, "signup.html.tpl", message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
