package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func pingRoute(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
