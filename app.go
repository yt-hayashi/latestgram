package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/", top)
	http.ListenAndServe(":8080", nil)
}

func top(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:mysql@tcp(127.0.0.1:33306)/test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println("DB Error! --> ", err.Error())
		os.Exit(1)
	}

	tmp := template.Must(template.ParseFiles("template/top.html.tpl"))

	rows, err := db.Query("SELECT users.id, users.name, posts.id, posts.img_name FROM users as users, posts as posts WHERE(users.id=posts.user_id) ORDER BY posts.created_at limit 50")
	if err != nil {
		fmt.Println(err.Error())
	}

	var contents contents

	for rows.Next() {
		var userID int
		var userName string
		var postID int
		var imgName string
		if err := rows.Scan(&userID, &userName, &postID, &imgName); err != nil {
			fmt.Println(err.Error())
		}
		//fmt.Println(userID, imgName)

		contents = append(contents, makeContent(userName, imgName))
	}

	if err := tmp.ExecuteTemplate(w, "top.html.tpl", contents); err != nil {
		fmt.Println(err.Error())
	}
}

type content struct {
	NameText string
	ImgPath  string
}

type contents []*content

func makeContent(name, path string) (making *content) {
	making = new(content)
	making.NameText = name
	making.ImgPath = path
	return making
}
