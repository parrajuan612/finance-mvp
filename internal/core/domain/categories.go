package domain

type Categories string

const (
	CatRestaurante     Categories = "Restaurante/Comida/Supermercado"
	CatRopa            Categories = "Ropa/Calzado"
	CatSalud           Categories = "Salud/Deporte"
	CatServicios       Categories = "Servicios/Suscripciones"
	CatTransporte      Categories = "Transporte"
	CatViajes          Categories = "Viajes"
	CatEducacion       Categories = "Educación"
	CatEntretenimiento Categories = "Entretenimiento"
	CatHogar           Categories = "Hogar/Alquiler"
	CatSalario         Categories = "Salario"
	CatOtrosIngresos   Categories = "Otros Ingresos"
	CatOtrosGastos     Categories = "Otros Gastos"
)

var CategoryIDMap = map[Categories]int{
	CatRestaurante:     1,
	CatRopa:            2,
	CatSalud:           3,
	CatServicios:       4,
	CatTransporte:      5,
	CatViajes:          6,
	CatEducacion:       7,
	CatEntretenimiento: 8,
	CatHogar:           9,
	CatSalario:         10,
	CatOtrosIngresos:   11,
	CatOtrosGastos:     12,
}
