package main

import (
	"database/sql"
	"log"
	// "time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	from, err := sql.Open("mysql", "root:test_test@tcp(127.0.0.1:3306)/wordpress")
	if err != nil {
		log.Fatalln("open", err)
	}

	defer from.Close()

	rows, err := from.Query("select post_title, post_content, post_date from wp_posts where post_status='publish' and post_type='post'")
	if err != nil {
		log.Fatalln("query", err)
	}

	defer rows.Close()

	to, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalln("open to", err)
	}

	for rows.Next() {
		var title, content, date string

		if err := rows.Scan(&title, &content, &date); err != nil {
			log.Fatalln("scan", err)
		}

		// log.Println(title, date, len(content))

		_, err := to.Exec("insert into post (title, slug, create_date, update_date, category_id, tag_id, author_id, status, content) values (?, ?, ?, ?, 1, 3, 1, 'Publish', ?)", title, title, date, date, content)
		if err != nil {
			log.Println("exec", err, title)
		}
	}
}
