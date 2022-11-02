package main

import (
    "context"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
    "time"
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
    srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")

}
