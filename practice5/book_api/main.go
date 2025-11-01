package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"


)

type Book struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Genre  string  `json:"genre"`
	Price  int    `json:"price"`
}
func main() {
    connStr := "host=localhost port=5432 user=postgres password=1234 dbname=bookdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() //defer откладывает выполнение указанной функции до тех пор, пока текущая функция не закончится (не вернётся).

	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        genre := r.URL.Query().Get("genre")
        sort := r.URL.Query().Get("sort")
        limitStr := r.URL.Query().Get("limit")
        offsetStr := r.URL.Query().Get("offset")

        query := "SELECT id, title, author, genre, price FROM books"
        var args []interface{}
        var conditions []string

        if genre != "" {
            conditions = append(conditions, "genre = $1")
            args = append(args, genre)
        }

        if len(conditions) > 0 {
            query += " WHERE " + conditions[0]
        }

        if sort == "price_asc" {
            query += " ORDER BY price ASC"
        } else if sort == "price_desc" {
            query += " ORDER BY price DESC"
        }

        limit, _ := strconv.Atoi(limitStr)
        offset, _ := strconv.Atoi(offsetStr)

        if limit > 0 {
            args = append(args, limit)
            query += fmt.Sprintf(" LIMIT $%d", len(args))
        }
        if offset > 0 {
            args = append(args, offset)
            query += fmt.Sprintf(" OFFSET $%d", len(args))
        }

        rows, err := db.Query(query, args...)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var books []Book
        for rows.Next() {
            var b Book
            if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Genre, &b.Price); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            books = append(books, b)
        }

        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("X-Query-Time", time.Since(start).String())
        json.NewEncoder(w).Encode(books)
    })

    fmt.Println("🚀 Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))

}
