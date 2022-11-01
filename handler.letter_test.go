package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIndexPageRoute(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	repo := getMemoryRepository()
	keyGen := getTestKeyGen()
	setupRouter(r, &repo, &keyGen)
	req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "<h1 class=\"header-title\">Secret letters 📨</h1>")
}

func TestSaveLetter(t *testing.T) {
	text := "test message"
	ttl := "60"
	postData := strings.NewReader(fmt.Sprintf("text=%s&ttl=%s", text, ttl))

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	repo := getMemoryRepository()
	keyGen := getTestKeyGen()
	setupRouter(r, &repo, &keyGen)

	req, _ := http.NewRequestWithContext(ctx, "POST", "/letter/new", postData)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	assert.Equal(t, 302, w.Code)

	key := keyGen.Get()
	saved_message, _ := repo.Get(key)
	assert.Equal(t, text, saved_message)
}

func TestReadLetter(t *testing.T) {
	testLetter := "test letter"
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	repo := getMemoryRepository()
	keyGen := getTestKeyGen()
	key := keyGen.Get()
	repo.Set(key, testLetter, 60)
	setupRouter(r, &repo, &keyGen)
	resultChannel := make(chan int, 2)

	go func(c chan int, router *gin.Engine) {
		request, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("/%s", key), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)

		resultChannel <- w.Code
	}(resultChannel, r)

	go func(c chan int, router *gin.Engine) {
		request, _ := http.NewRequest("GET", fmt.Sprintf("/%s", key), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)

		resultChannel <- w.Code
	}(resultChannel, r)

	firstCode := <-resultChannel
	secondCode := <-resultChannel
	assert.Equal(t, 200, firstCode)
	assert.Equal(t, 404, secondCode)
}
