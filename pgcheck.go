package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres@localhost/postgres?sslmode=disable")

	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}

func main() {
	http.HandleFunc("/replica", handler)
	http.HandleFunc("/master", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT pg_is_in_recovery()")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for rows.Next() {
		var result bool
		err = rows.Scan(&result)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		switch r.URL.Path {
		case "/master":
			if result {
				http.Error(w, "I'm not the master!", 403)
			} else {
				fmt.Fprint(w, "I'm the master!")
			}
		case "/replica":
			if result {
				fmt.Fprint(w, "I'm a replica!")
			} else {
				http.Error(w, "I'm not a replica!", 403)
			}
		}
	}
}
