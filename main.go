package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"os"
)

var db *pgxpool.Pool

/*
func (h *RecommendationHandler) recommendForNewUserHandler(w http.ResponseWriter, r *http.Request) {
	sequence := r.URL.Query()["seq"]
	recCount := r.URL.Query().Get("rec")

	recommendations, _ := getColdRecommendations(sequence, recCount)
	result, _ := json.Marshal(recommendations)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}*/

func main() {
	var err error
	db, err = pgxpool.New(context.Background(), "postgresql://postgres@localhost:5432/bambook")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	mux := http.NewServeMux()
	mux.Handle("/recommend/", &RecommendationHandler{})
	mux.Handle("/books/", &BooksHandler{})
	mux.Handle("/users/", &UserHandler{})
	http.ListenAndServe(":5050", mux)
}