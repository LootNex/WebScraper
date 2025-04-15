package scraping

import (
	"math/rand"
	"strings"
	"time"
)

// GetCurrentPrice имитирует парсинг цены с Wildberries или Ozon.
func GetCurrentPrice(url string) float64 {
	if !isValidURL(url) {
	return 0
	}

	// TODO: Реализовать настоящий парсинг с colly или goquery.
	// Здесь заглушка: возвращаем случайную цену.
	rand.Seed(time.Now().UnixNano())
	return 1000 + rand.Float64()*1000
}

func isValidURL(url string) bool {
	return strings.Contains(url, "wildberries") ||
	strings.Contains(url, "ozon") ||
	strings.Contains(url, "wb.ru")
}