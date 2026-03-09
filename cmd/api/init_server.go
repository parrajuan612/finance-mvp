package api

import (
	"finanzas-mvp/internal/adapters/handlers"
	"finanzas-mvp/internal/adapters/repository/postgres"
	"finanzas-mvp/internal/core/domain"
	"finanzas-mvp/internal/core/services"
	"log"

	"github.com/gin-gonic/gin"
)

func InitServer(r *gin.Engine) {

	r.LoadHTMLGlob("web/templates/*")
	r.Static("/static", "./web/static")

	// 🔹 Inicializar base de datos
	db, err := postgres.InitPostgres()
	if err != nil {
		// Esto detendrá el programa y te dirá si es usuario/password/host incorrecto
		log.Fatalf("❌ ERROR CRÍTICO DB: %v", err)
	}
	log.Println("✅ Conexión a base de datos exitosa")

	// 🔹 Repositories
	statementRepo := postgres.NewStatementRepository(db)

	// 🔹 Reglas de categorización
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
	statementService := services.NewStatementService(
		categorizer,
		statementRepo,
	)

	// 🔹 Handler
	uploadHandler := handlers.NewUploadHandler(statementService)

	// 🔹 Rutas
	InitRoutes(r, uploadHandler)

	log.Println("Server running on :9080")
	err = r.Run(":9080")
	if err != nil {
		log.Fatalf("❌ ERROR AL INICIAR SERVIDOR: %v", err)
	}
}
