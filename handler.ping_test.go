package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	repo := getMemoryRepository()
	keyGen := getTestKeyGen()
	setupRouter(r, &repo, &keyGen)
	req, _ := http.NewRequestWithContext(ctx, "GET", "/ping", nil)
	r.ServeHTTP(w, req)

	result := w.Result()
	defer result.Body.Close()
	data, _ := io.ReadAll(result.Body)
	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"message":"pong"}`, string(data))
}
