package services

import (
	"strings"

	"finanzas-mvp/internal/core/domain"
)

type Categorizer struct {
	rules         domain.CategoryRules
	categoryIDMap map[string]int
}

func NewCategorizer(rules domain.CategoryRules, categories []domain.Category) *Categorizer {

	catMap := make(map[string]int)
	for _, cat := range categories {
		catMap[cat.Name] = cat.ID
	}

	return &Categorizer{
		rules:         rules,
		categoryIDMap: catMap,
	}
}

func (c *Categorizer) Categorize(mov *domain.Movement) {
	desc := normalizeText(mov.Description)

	var keywordsMap map[domain.Categories][]string
	var defaultSlug domain.Categories

	if mov.Type == domain.TypeExpense {
		keywordsMap = c.rules.Expenses
		defaultSlug = domain.CatOtrosGastos
	} else {
		keywordsMap = c.rules.Incomes
		defaultSlug = domain.CatOtrosIngresos
	}

	targetCategoryName := c.findMatchingCategory(desc, keywordsMap, string(defaultSlug))

	mov.CategoryName = targetCategoryName

	if id, exists := c.categoryIDMap[targetCategoryName]; exists {
		mov.CategoryID = id
	}
}

func (c *Categorizer) findMatchingCategory(desc string, rules map[domain.Categories][]string, defaultCat string) string {
	for slug, keywords := range rules {
		for _, kw := range keywords {
			if strings.Contains(desc, kw) {
				return string(slug)
			}
		}
	}
	return defaultCat
}

func normalizeText(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))

	replacements := strings.NewReplacer(
		"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u",
		"ä", "a", "ë", "e", "ï", "i", "ö", "o", "ü", "u",
	)

	return replacements.Replace(text)
}
