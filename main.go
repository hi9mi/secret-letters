package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	r.Use(connectRepository(repo))
	r.Use(connectKeyGen(kg))
	r.GET("/ping", pingRoute)
	r.GET("/", indexPageRoute)
	r.POST("/letter/new", saveLetterRoute)
	r.GET("/key/:key", getKeyRoute)
	r.GET("/:key", getLetterRoute)
	r.NoRoute(func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	})
	return r
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	repo := getMemoryRepository()
	kg := getTestKeyGen()

	r := setupRouter(gin.Default(), &repo, &kg)

	r.Run(":8080")
}