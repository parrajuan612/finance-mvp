package domain

type CategoryRules struct {
	Expenses map[Categories][]string
	Incomes  map[Categories][]string
}

func GetDefaultRules() CategoryRules {
	return CategoryRules{
		Expenses: map[Categories][]string{
			CatRestaurante: {
				"exito", "carulla", "jumbo", "d1", "ara", "mcdonalds", "rappi",
				"restaurante", "kfc", "starbucks", "bm 126 sas", "domino s pizza",
				"didi food", "pan latino", "goyurt", "el tambor", "andres paradero",
				"frisby", "crepes", "el corral", "surtifruver",
			},
			CatTransporte: {
				"transporte masivo", "uber", "cabify", "didi", "taxis", "eds",
				"terpel", "texaco", "peaje", "sitp", "transmilenio", "bold*corporacio",
			},
			CatViajes: {
				"latam", "avianca", "viva air", "booking", "airbnb", "despegar",
				"hotel", "lifemiles", "expedia",
			},
			CatServicios: {
				"acueducto", "enel", "vanti", "claro", "movistar", "tigo", "epm",
				"telecomuni", "apple.com", "icloud", "google storage", "etb",
				"emermedica", "seguros",
			},
			CatEntretenimiento: {
				"netflix", "spotify", "playstation", "psn", "disney+", "hbo",
				"cine colombia", "procinal", "hotmart", "twitch", "steam",
			},
			CatRopa: {
				"zara", "hm", "adidas", "nike", "falabella", "tennis", "bershka",
				"erato", "gef", "punto blanco", "chevignon", "patprimo", "decathlon",
			},
			CatHogar: {
				"homecenter", "easy", "ferreteria", "administracion", "arriendo",
				"sodimac", "amazon", "mercado libre", "mercadoli", "alkosto",
				"ktronix", "ikea",
			},
			CatSalud: {
				"cruz verde", "farmatodo", "colsubsidio", "eps", "sanitas",
				"smartfit", "bodytech", "drogueria", "medico", "odontologia",
				"assul medical", "beconfiden",
			},
			CatEducacion: {
				"universidad", "colegio", "udemy", "platzi", "coursera", "pension",
			},
			CatOtrosGastos: {
				"impto gobierno", "4x1000", "cuota manejo", "comision", "bold",
				"ajuste", "intereses mora", "palatino",
			},
		},
		Incomes: map[Categories][]string{
			CatSalario: {
				"nomina", "pago quincena", "salario", "honorarios",
				"pago empresa",
			},
			CatOtrosIngresos: {
				"abono intereses", "consignacion", "reintegro",
				"transferencia desde nequi", "gracias por tu pago",
				"devuelto", "devolucion", "mercado pago", "premiumpe",
			},
		},
	}
}
