package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	mux := http.NewServeMux()

	// маршруты
	mux.Handle("/user", apiMiddleware(http.HandlerFunc(userHandler)))

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// userHandler обрабатывает и GET, и POST
func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			http.Error(w, `{"error": "invalid id"}`, http.StatusBadRequest)
			return
		}

		resp := map[string]int{"user_id": id}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if r.Method == http.MethodPost {
		var body struct {
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
			http.Error(w, `{"error": "invalid name"}`, http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		resp := map[string]string{"created": body.Name}
		json.NewEncoder(w).Encode(resp)
		return
	}

	http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
}

// Middleware
func apiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL.Path)

		key := r.Header.Get("X-API-Key")
		if key != "secret123" {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
