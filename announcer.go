package main

import (
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	password := os.Getenv("POSTGRESPASSWORD")
	if password == "" {
		fmt.Println("Failed to get $POSTGRESPASSWORD. Have you sourced authdetails.sh?")
	}

	db, err := sql.Open("postgres", "user=postgres password='"+password+"' dbname=tolljardates sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM dates WHERE due = current_date;")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var due time.Time
		var when time.Time
		var to string
		var from string
		var what string
		err = rows.Scan( &due, &when, &to, &from, &what)
		fmt.Println(from)
	}
}
