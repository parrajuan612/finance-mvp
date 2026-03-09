package handlers

import (
	"finanzas-mvp/internal/core/ports"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	service ports.StatementService
}

func NewUploadHandler(service ports.StatementService) *UploadHandler {
	return &UploadHandler{
		service: service,
	}
}
func (h *UploadHandler) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
