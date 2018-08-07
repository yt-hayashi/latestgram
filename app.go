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
	db, err := sql.Open("mysql", "root:mysql@tcp(127.0.0.1:33306)/test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println("DB Error! --> ", err.Error())
		os.Exit(1)
	}

	http.HandleFunc("/", top)
	http.ListenAndServe(":8080", nil)
}

func top(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles("template/top.html.tpl"))

	textBody := "This is Test."

	if err := tmp.ExecuteTemplate(w, "top.html.tpl", textBody); err != nil {
		fmt.Println(err.Error())
	}
}
