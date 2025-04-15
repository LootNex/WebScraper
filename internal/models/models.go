package models

import "time"

type Item struct {
	ID           int
	UserID       int64
	Link         string
	CurrentPrice float64
	LastChecked  time.Time
	CreatedAt    time.Time
}

type PriceHistory struct {
	ID        int
	ItemID    int
	Price     float64
	CheckedAt time.Time
	}