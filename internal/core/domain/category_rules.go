package domain

type CategoryRules struct {
	Expenses map[Categories][]string
	Incomes  map[Categories][]string
}

func GetDefaultRules() CategoryRules {
	return CategoryRules{
		Expenses: map[Categories][]string{
			CatServicios:   {"acueducto", "enel", "vanti", "claro", "movistar", "tigo", "epm", "netflix", "spotify", "telecomuni"},
			CatTransporte:  {"transporte masivo", "uber", "cabify", "didi", "taxis", "eds", "terpel", "texaco", "peaje", "sitp"},
			CatHogar:       {"homecenter", "easy", "ferreteria", "administracion", "arriendo", "sodimac"},
			CatRestaurante: {"exito", "carulla", "jumbo", "d1", "ara", "mcdonalds", "rappi", "restaurante", "kfc", "starbucks", "bm 126 sas"},
			CatRopa:        {"zara", "hm", "adidas", "nike", "falabella", "tennis", "bershka", "erato"},
			CatSalud:       {"cruz verde", "farmatodo", "colsubsidio", "eps", "sanitas", "smartfit", "bodytech"},
			CatOtrosGastos: {"impto gobierno", "4x1000", "cuota manejo", "comision", "bold"},
		},
		Incomes: map[Categories][]string{
			CatSalario:       {"nomina", "pago quincena", "salario"},
			CatOtrosIngresos: {"abono intereses", "consignacion", "reintegro", "transferencia desde nequi"},
		},
	}
}
