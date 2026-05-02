package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count string
		city  string
		want  int
	}{
		{"0", "moscow", 0},
		{"1", "moscow", 1},
		{"2", "moscow", 2},
		{"100", "moscow", 5},

		{"0", "tula", 0},
		{"1", "tula", 1},
		{"100", "tula", 3},

		{"", "moscow", 5},
	}

	for _, re := range requests {
		response := httptest.NewRecorder()
		var str string = fmt.Sprintf("/cafe?city=%s", re.city)
		if re.count != "" {
			str = str + "&count=" + re.count
		}
		req := httptest.NewRequest("GET", str, nil)
		handler.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
		body := response.Body.String()
		var x []string
		if body != "" {
			x = strings.Split(body, ",")
		}
		assert.Len(t, x, re.want)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
	}{
		{search: "фасоль", wantCount: 0},
		{search: "кофе", wantCount: 2},
		{search: "вилка", wantCount: 1},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()

		req := httptest.NewRequest(
			"GET",
			fmt.Sprintf("/cafe?city=moscow&search=%s", v.search),
			nil,
		)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)

		body := response.Body.String()

		var cafes []string
		if body != "" {
			cafes = strings.Split(body, ",")
		}
		assert.Len(t, cafes, v.wantCount)
		for _, cafe := range cafes {
			assert.Contains( // благодаря ей можно пррверять как строки так и мапы так и слайсы
				t,
				strings.ToLower(cafe),     // строка
				strings.ToLower(v.search), //. подстрока
			)
		}
	}
}
