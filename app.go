package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"

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

	fmt.Println("method:", r.Method)
	if r.Method == http.MethodPost {
		// signup時の処理
		return
	}

	// getのrender処理
	fmt.Println("method:", r.Method)
	//レスポンスの解析
	r.ParseForm()
	fmt.Println(r.Form)
	for i, v := range r.Form {
		fmt.Println("index:", i)
		fmt.Println("value:", v)
	}

	message := "Please Input"
	if err := tmp.ExecuteTemplate(w, "signup.html.tpl", message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
