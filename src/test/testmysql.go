package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func insert(db *sql.DB) {
	stmt, err := db.Prepare(`INSERT INTO AccessLogs (Ip1, date_reg, request_url, status_code, body_size, from_url, agent, request_time) VALUES (?, ?,?, ?, ?,?, ?, ?)`)

	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return
	}
	stmt.Exec(1, 1, 1, 1, 1, 1, 1, 1)

}

func main() {
	db, err := sql.Open("mysql", "root:tuandai1921688190@tcp(192.168.8.190:3306)/Tuandai_Log?charset=utf8")
	if err != nil {
		log.Fatalf("Open database error: %s\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	insert(db)

	rows, err := db.Query("select * from AccessLogs where id = ?", 1)
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()
	var Id int
	var Ip1 string
	for rows.Next() {
		err := rows.Scan(&Id, &Ip1)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(Id, Ip1)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
