package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:mysql@tcp(127.0.0.1:33306)/test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println("DB Error!!!")
	} else {
		fmt.Println("DB OK.")
	}

	addUser := 0
	addPost := 0
	addComment := 0
	//ユーザーの追加
	for i := 0; i < 10; i++ {
		countUser := strconv.Itoa(i)
		_, err := db.Exec(`
		INSERT INTO users(name, password) VALUES(?, ?)`, "user"+countUser, "pass"+countUser)

		if err != nil {
			panic(err.Error())
		}
		addUser++
	}

	//postの追加
	for i := 0; i < 100; i++ {
		countPost := strconv.Itoa(i)
		_, err := db.Exec(`
		INSERT INTO posts(user_id, img_name) VALUES(?, ?)`, i, "img_"+countPost)

		if err != nil {
			panic(err.Error())
		}
		addPost++
	}

	for i := 0; i < 50; i++ {
		countComment := strconv.Itoa(i)
		_, err := db.Exec(`
		INSERT INTO comments(post_id, user_id, comment_body) VALUES(?, ?, ?)`, i, i, "This is comment"+countComment)

		if err != nil {
			panic(err.Error())
		}
		addComment++
	}
	fmt.Printf("added user: %d\n added post: %d\n added comments: %d\n", addUser, addPost, addComment)

}
