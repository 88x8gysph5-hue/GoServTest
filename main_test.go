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
		count int
	}{
		{count: 0},
		{count: 1},
		{count: 2},
		{count: 100},
	}

	cityRequest := []string{"moscow", "tula"}

	for _, city := range cityRequest {
		for _, v := range requests {

			maxLen := len(cafeList[city])
			expected := min(v.count, maxLen)

			response := httptest.NewRecorder()
			req := httptest.NewRequest(
				"GET",
				fmt.Sprintf("/cafe?count=%d&city=%s", v.count, city),
				nil,
			)

			handler.ServeHTTP(response, req)

			body := response.Body.String()

			var resultCount int
			if body == "" {
				resultCount = 0
			} else {
				resultCount = len(strings.Split(body, ","))
			}

			assert.Equal(t, expected, resultCount)
		}
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

		body := response.Body.String()

		var cafes []string
		if body != "" {
			cafes = strings.Split(body, ",")
		}
		assert.Equal(t, v.wantCount, len(cafes))
		for _, cafe := range cafes {
			assert.True(
				t,
				strings.Contains(
					strings.ToLower(cafe),
					strings.ToLower(v.search),
				),
			)
		}
	}
}
