package handlers

type MovementForm struct {
	Date        string  `json:"date" form:"movements[].date"` // form tag opcional si usas binding clásico
	Description string  `json:"description" form:"movements[].description"`
	CategoryID  int     `json:"category_id" form:"movements[].category_id"`
	Amount      float64 `json:"amount" form:"movements[].amount"`
	Type        string  `json:"type" form:"movements[].type"`
}

type SaveMovementsForm struct {
	PeriodMonth string         `form:"period_month"`
	FileName    string         `form:"file_name"`
	Movements   []MovementForm `form:"movements"`
}
