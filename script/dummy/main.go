package main

import (
	"database/sql"
	"fmt"
	"os"
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
		fmt.Println("DB Error! --> ", err.Error())
		os.Exit(1)
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
			fmt.Println("Error! User didn't add.", err.Error())
		}
		addUser++
	}

	//postの追加
	for i := 0; i < 100; i++ {
		countPost := strconv.Itoa(i)
		_, err := db.Exec(`
		INSERT INTO posts(user_id, img_name) VALUES(?, ?)`, i, "img_"+countPost)

		if err != nil {
			fmt.Println("Error! Post didn't add.", err.Error())
		}
		addPost++
	}

	//commentの追加
	for i := 0; i < 50; i++ {
		countComment := strconv.Itoa(i)
		_, err := db.Exec(`
		INSERT INTO comments(post_id, user_id, comment_body) VALUES(?, ?, ?)`, i, i, "This is comment"+countComment)

		if err != nil {
			fmt.Println("Error! Comments didn't add.", err.Error())
		}
		addComment++
	}

	fmt.Printf("Done...\n added user: %d\n added post: %d\n added comments: %d\n", addUser, addPost, addComment)

}
