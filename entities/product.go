package entities

import (
	"time"
)

type Product struct {
	ProductID    string  `pg:",pk" json:"product_id"`
	ProductName  string  `pg:"uniqe" json:"product_name"`
	ProductStock int     `json:"product_stock"`
	Price        float64 `json:"price"`
	VatRate      int     `json:"vat_rate"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
