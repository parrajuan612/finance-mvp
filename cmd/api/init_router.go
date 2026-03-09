package api

import (
	"finanzas-mvp/internal/adapters/handlers"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, uploadHandler *handlers.UploadHandler) {

	r.GET("/", uploadHandler.Home)

	r.POST("/upload", uploadHandler.Upload)

	r.POST("/save-movements", uploadHandler.SaveMovements)
}
