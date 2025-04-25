package models

type Item struct {
	Name            string  `json:"name"`
	StartPrice      float32 `json:"start_price"`
	CurrentPrice    float32 `json:"current_price"`
	DifferencePrice float32 `json:"difference_price"`
}
