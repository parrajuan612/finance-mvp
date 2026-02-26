package handlers

import (
	"finanzas-mvp/internal/adapters/parsers"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) Upload(c *gin.Context) {

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "archivo requerido")
		return
	}
	defer file.Close()

	password := c.PostForm("password")

	fmt.Println("Password recibida:", password)

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)

	path := filepath.Join("storage", filename)

	out, err := os.Create(path)
	if err != nil {
		c.String(http.StatusInternalServerError, "no se pudo guardar archivo")
		return
	}

	_, err = io.Copy(out, file)
	if err != nil {
		c.String(http.StatusInternalServerError, "error guardando archivo")
		return
	}

	// CRÍTICO: cerrar antes de usar qpdf
	err = out.Close()
	if err != nil {
		c.String(http.StatusInternalServerError, "error cerrando archivo")
		return
	}

	fmt.Println("Archivo guardado en:", path)

	// LEER PDF
	text, err := parsers.ReadPDFUsingQPDF(path, password)

	movs, err := parsers.ParseBancolombiaText(text)
	if err != nil {
		fmt.Println("Parser error:", err)
	} else {
		fmt.Println("Movimientos detectados:", len(movs))
		for _, m := range movs {
			fmt.Printf("%s | %s | %.2f | %s\n",
				m.Date.Format("2006-01-02"),
				m.Description,
				m.Amount,
				m.Type,
			)
		}
	}
}
