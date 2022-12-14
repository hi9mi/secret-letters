package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed static
var staticEmbedFS embed.FS

//go:embed templates
var tmplEmbed embed.FS

type staticFS struct {
	fs fs.FS
}

func (sfs *staticFS) Open(name string) (fs.File, error) {
	return sfs.fs.Open(filepath.Join("static", name))
}

var staticEmbed = &staticFS{staticEmbedFS}

func connectRepository(repo *Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("repo", repo)
	}
}

func connectKeyGen(kg *KeyGen) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("keygen", kg)
	}
}

func setupRouter(r *gin.Engine, repo *Repository, kg *KeyGen) *gin.Engine {
	tmpl := template.Must(template.ParseFS(tmplEmbed, "templates/*.html"))
	r.SetHTMLTemplate(tmpl)
	r.StaticFS("/static", http.FS(staticEmbed))
	r.StaticFile("/robots.txt", "./static/robots.txt")
	r.Use(connectRepository(repo))
	r.Use(connectKeyGen(kg))
	r.GET("/ping", pingRoute)
	r.GET("/", indexPageRoute)
	r.POST("/letter/new", saveLetterRoute)
	r.GET("/key/:key", getKeyRoute)
	r.GET("/letter/:key", getLetterRoute)
	r.NoRoute(func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	})
	return r
}

func main() {
	repo := getRedisRepository()
	kg := getUUIDKeyGen()

	r := setupRouter(gin.Default(), &repo, &kg)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Server started ...")

	var requestUrl string

	if strings.Contains(os.Getenv("LOCAL"), "true") {
		requestUrl = fmt.Sprintf("http://localhost%s/ping", srv.Addr)
	} else {
		requestUrl = fmt.Sprintf("https://%s/ping", "secret-letters.herokuapp.com")
		log.Println("-------------------------------------------")
		log.Println(requestUrl)
		log.Println("-------------------------------------------")
	}

	for interval := range time.Tick(10 * time.Minute) {
		res, _ := http.Get(requestUrl)
		log.Println("-------------------------------------------")
		log.Printf("%v: %v\n", interval, res.StatusCode)
		log.Println("-------------------------------------------")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Println("timeout of 10 seconds.")
	log.Println("Server exiting")
}
