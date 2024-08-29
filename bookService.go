package main

import (
	"bambook-backend/core"
	"context"
	"encoding/json"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"net/http"
)

type BooksHandler struct{}

func (h *BooksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.searchBooks(w, r)
}

func (h *BooksHandler) searchBooks(w http.ResponseWriter, r *http.Request) {
	searchToken := r.URL.Query().Get("search")
	searchToken = fmt.Sprintf("%s%%", searchToken)
	fmt.Println(searchToken)
	var books []core.Book
	_ = pgxscan.Select(context.Background(), db, &books, "select * from books where title ilike $1", searchToken)
	result, _ := json.Marshal(books)
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}