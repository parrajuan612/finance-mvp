package main

import (
	"finanzas-mvp/internal/adapters/handlers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	r := gin.Default()

	r.LoadHTMLGlob("web/templates/*")
	r.Static("/static", "./web/static")

	uploadHandler := handlers.NewUploadHandler()

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/upload", uploadHandler.Upload)

	log.Println("Server running on :9080")

	r.Run(":9080")
}
