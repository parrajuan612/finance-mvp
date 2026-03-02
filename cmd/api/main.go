package main

import (
	"finanzas-mvp/internal/adapters/handlers"
	"finanzas-mvp/internal/core/domain"
	"finanzas-mvp/internal/core/services"
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
	rules := domain.GetDefaultRules()
	dBCategories := []domain.Category{
		{ID: 1, Name: string(domain.CatRestaurante), IsActive: true},
		{ID: 2, Name: string(domain.CatRopa), IsActive: true},
		{ID: 3, Name: string(domain.CatSalud), IsActive: true},
		{ID: 4, Name: string(domain.CatServicios), IsActive: true},
		{ID: 5, Name: string(domain.CatTransporte), IsActive: true},
		{ID: 6, Name: string(domain.CatViajes), IsActive: true},
		{ID: 7, Name: string(domain.CatEducacion), IsActive: true},
		{ID: 8, Name: string(domain.CatEntretenimiento), IsActive: true},
		{ID: 9, Name: string(domain.CatHogar), IsActive: true},
		{ID: 10, Name: string(domain.CatSalario), IsActive: true},
		{ID: 11, Name: string(domain.CatOtrosIngresos), IsActive: true},
		{ID: 12, Name: string(domain.CatOtrosGastos), IsActive: true},
	}
	categorizer := services.NewCategorizer(rules, dBCategories)
	statementService := services.NewStatementService(categorizer)
	uploadHandler := handlers.NewUploadHandler(statementService)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/upload", uploadHandler.Upload)

	log.Println("Server running on :9080")

	r.Run(":9080")
}
