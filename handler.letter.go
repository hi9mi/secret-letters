package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const MaxTTL = 86400
const MinTTL = 60

func indexPageRoute(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{"maxTTL": MaxTTL, "minTTL": MinTTL})
}

func saveLetterRoute(ctx *gin.Context) {
	letter := &Letter{}

	if err := ctx.ShouldBind(letter); err != nil {
		verrs := err.(validator.ValidationErrors)
		messages := make([]string, len(verrs))
		for i, verr := range verrs {
			messages[i] = fmt.Sprintf(
				"%s is required, but was empty.",
				verr.Field())
		}
		ctx.HTML(http.StatusBadRequest, "index.html", gin.H{"errors": messages, "maxTTL": MaxTTL, "minTTL": MinTTL})
		return
	}

	repo := ctx.Value("repo").(*Repository)
	key := (*ctx.Value("keygen").(*KeyGen)).Get()

	err := (*repo).Set(key, letter.Text, letter.Ttl)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error-page.html", gin.H{"cause": err.Error(), "code": http.StatusInternalServerError})
		return
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/key/%s", key))
}

func getKeyRoute(ctx *gin.Context) {
	key, ok := ctx.Params.Get("key")

	if !ok {
		ctx.HTML(http.StatusBadRequest, "error-page.html", gin.H{"cause": "Not existing key", "code": http.StatusBadRequest})
		return
	}

	repo := (*ctx.Value("repo").(*Repository))
	ok = repo.Check(key)

	if !ok {
		ctx.Redirect(http.StatusMovedPermanently, "/")
		return
	}

	var url string
	if strings.Contains(os.Getenv("LOCAL"), "true") {
		url = fmt.Sprintf("http://%s/%s", ctx.Request.Host, key)
	} else {
		url = fmt.Sprintf("https://%s/%s", ctx.Request.Host, key)
	}

	ctx.HTML(http.StatusOK, "key.html", gin.H{"url": url})
}

func getLetterRoute(ctx *gin.Context) {
	key, ok := ctx.Params.Get("key")
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	repo := (*ctx.Value("repo").(*Repository))
	letter, err := repo.Get(key)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "error-page.html", gin.H{"cause": err.Error(), "code": http.StatusNotFound})
		return
	}
	ctx.HTML(http.StatusOK, "letter.html", gin.H{"letter": letter})
}
