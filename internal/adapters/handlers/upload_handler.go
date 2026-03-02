package handlers

import (
	"finanzas-mvp/internal/core/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	service *services.StatementService
}

func NewUploadHandler(service *services.StatementService) *UploadHandler {
	return &UploadHandler{
		service: service,
	}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "archivo requerido")
		return
	}
	defer file.Close()

	password := c.PostForm("password")
	bankID := c.PostForm("bank_id")

	movements, err := h.service.ProcessStatement(file, header.Filename, password, bankID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "results.html", gin.H{
		"Movements": movements,
		"Filename":  header.Filename,
	})
}
