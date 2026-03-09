package handlers

import (
	"encoding/json"
	"finanzas-mvp/internal/core/domain"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *UploadHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "archivo requerido")
		return
	}
	defer file.Close()

	password := c.PostForm("password")
	bankID := c.PostForm("bank_id")

	movements, periodMonth, err := h.service.ProcessStatement(file, header.Filename, password, bankID)
	if err != nil {

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "results.html", gin.H{
		"Movements":   movements,
		"Periodmonth": periodMonth,
		"Filename":    header.Filename,
	})
}
func (h *UploadHandler) SaveMovements(c *gin.Context) {
	period := c.PostForm("period_month")
	filename := c.PostForm("file_name")

	// 1) intentar leer JSON generado por JS
	mJson := c.PostForm("movements_json")
	var movements []MovementForm

	if mJson != "" {
		if err := json.Unmarshal([]byte(mJson), &movements); err != nil {
			c.String(http.StatusBadRequest, "movements_json inválido: %v", err)
			return
		}
	} else {

		var fallback SaveMovementsForm
		if err := c.ShouldBind(&fallback); err != nil {
			c.String(http.StatusBadRequest, "no se enviaron movimientos")
			return
		}

	}

	// ahora tienes: period, filename, movements[]

	// Convertir MovementForm -> domain.Movement si vas a pasarlo al service
	var domainMovs []domain.Movement
	for _, m := range movements {
		// Parsear fecha (asumiendo YYYY-MM-DD)
		var parsedDate time.Time
		if m.Date != "" {
			t, err := time.Parse("2006-01-02", m.Date)
			if err != nil {
				parsedDate = time.Now()
			} else {
				parsedDate = t
			}
		} else {
			parsedDate = time.Now()
		}
		uuidUser, err := uuid.Parse("296f368f-f7b4-4388-8934-209e146de03c")
		uuidAccount, err := uuid.Parse("3bf374ea-db66-4f47-ab8a-7d0156c4440f")
		if err != nil {
			log.Fatal(err)
		}
		d := domain.Movement{
			ID:          uuid.Nil,
			UserID:      uuidUser,
			AccountID:   uuidAccount,
			StatementID: nil,
			CategoryID:  m.CategoryID,
			Date:        parsedDate,
			Description: m.Description,
			Amount:      m.Amount,
			Type:        domain.MovementType(m.Type),
			CreatedAt:   time.Now(),
		}
		domainMovs = append(domainMovs, d)
	}
	ctx := c.Request.Context()

	uuidUser, _ := uuid.Parse("296f368f-f7b4-4388-8934-209e146de03c")
	uuidAccount, _ := uuid.Parse("3bf374ea-db66-4f47-ab8a-7d0156c4440f")

	statement := domain.Statement{
		UserID:      uuidUser,
		AccountID:   uuidAccount,
		BankID:      1,
		FileName:    filename,
		PeriodMonth: period,
		UploadDate:  time.Now(),
		Status:      domain.StatusPending,
	}

	err := h.service.SaveStatementWithMovements(
		ctx,
		statement,
		domainMovs,
	)

	if err != nil {

		// detectar error de extracto duplicado
		if strings.Contains(err.Error(), "unique_statement_account_period") {
			c.String(400, "Este extracto ya fue cargado anteriormente")
			return
		}

		log.Println("error guardando movimientos:", err)

		c.String(500, "Error guardando movimientos")
		return
	}

	c.String(200, "Movimientos guardados correctamente")
}
