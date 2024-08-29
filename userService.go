package main

import (
	"bambook-backend/core"
	"context"
	"encoding/json"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strings"
)

type UserHandler struct{}

type CreateUserRequest struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	ReadingHistory []int  `json:"readingHistory"`
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "register"):
		h.Register(w, r)
		return
	case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "login"):
		h.Login(w, r)
		return
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request CreateUserRequest
	_ = json.NewDecoder(r.Body).Decode(&request)

	tx, err := db.Begin(context.Background())
	if err != nil {
		panic(err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(context.Background())
		} else {
			err = tx.Commit(context.Background())
		}

		if err != nil {
			panic(err)
		}
	}()

	_, err = db.Exec(
		context.Background(),
		"insert into users (email, password) VALUES ($1, $2)",
		request.Email,
		request.Password)

	rows, _ := db.Query(context.Background(), "select id, email from users where email = $1", request.Email)
	var user core.User
	_ = pgxscan.ScanOne(&user, rows)

	type UserBook struct {
		UserId int64
		BookId int64
	}
	var readingHistory [][]interface{}
	for _, bookId := range request.ReadingHistory {
		readingHistory = append(readingHistory, []interface{}{bookId, user.Id})
	}

	_, err = tx.CopyFrom(context.Background(), pgx.Identifier{"books_users"}, []string{"book_id", "user_id"}, pgx.CopyFromRows(readingHistory))

	result, _ := json.Marshal(user)
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Email string `json:"email"`
	}
	var request LoginRequest
	_ = json.NewDecoder(r.Body).Decode(&request)

	rows, _ := db.Query(context.Background(), "select id, email from users where email = $1", request.Email)
	var user core.User
	_ = pgxscan.ScanOne(&user, rows)

	result, _ := json.Marshal(user)
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}