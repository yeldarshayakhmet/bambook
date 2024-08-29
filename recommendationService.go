package main

import (
	"bambook-backend/core"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type RecommendationHandler struct{}

func (h *RecommendationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "user"):
		h.recommendForUserHandler(w, r)
		return
	case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "popular"):
		h.recommendPopularHandler(w, r)
		return
	}
}

func getUserRecommendations(user string, recCount string) ([]core.Book, error) {
	params := url.Values{}
	params.Add("n_rec", recCount)

	baseUrl := "http://localhost:8000/embed/recommend/user/"
	urlPath := fmt.Sprintf("%s%s?%s", baseUrl, user, params.Encode())

	recommendations, err := getRecommendations(urlPath)
	if err != nil {
		userId, _ := strconv.Atoi(user)

		type UserBooks struct {
			UserId int64
			BookId int64
		}
		var userBooks []UserBooks

		_ = pgxscan.Select(context.Background(), db, &userBooks, "select * from books_users where user_id = $1", userId)

		var bookIds []string
		for _, userBook := range userBooks {
			bookIds = append(bookIds, strconv.FormatInt(userBook.BookId, 10))
		}
		return getColdRecommendations(bookIds, recCount)
	}

	return recommendations, nil
}

func getPopularRecommendations(recCount string) ([]core.Book, error) {
	params := url.Values{}
	params.Add("n_rec", recCount)
	baseUrl := "http://localhost:8000/embed/recommend/popular"
	urlPath := fmt.Sprintf("%s?%s", baseUrl, params.Encode())
	return getRecommendations(urlPath)
}

func getColdRecommendations(sequence []string, recCount string) ([]core.Book, error) {
	params := url.Values{}
	for i := 0; i < len(sequence); i++ {
		params.Add("seq", sequence[i])
	}
	params.Add("n_rec", recCount)
	baseUrl := "http://localhost:8000/embed/recommend/cold"
	urlPath := fmt.Sprintf("%s?%s", baseUrl, params.Encode())
	return getRecommendations(urlPath)
}

func getRecommendations(url string) ([]core.Book, error) {
	response, err := http.Get(url)
	if response.Status != "200 OK" {
		return nil, errors.New("user does not exist")
	}

	responseBytes, _ := io.ReadAll(response.Body)
	_ = response.Body.Close()
	var recommendations []int
	_ = json.Unmarshal(responseBytes, &recommendations)

	var result []core.Book
	err = pgxscan.Select(context.Background(), db, &result, "select * from books where id = any($1)", recommendations)
	return result, err
}

func (h *RecommendationHandler) recommendForUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user key from header
	user := r.Header.Get("User")

	// Extract rec count from query parameter
	recCount := r.URL.Query().Get("rec")

	recommendations, _ := getUserRecommendations(user, recCount)
	result, _ := json.Marshal(recommendations)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (h *RecommendationHandler) recommendPopularHandler(w http.ResponseWriter, r *http.Request) {
	recCount := r.URL.Query().Get("rec")

	recommendations, _ := getPopularRecommendations(recCount)
	result, _ := json.Marshal(recommendations)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}