package controllers

import (
	"time"
)

type Basket struct {
	BasketID string `pg:",pk" json:"basket_id"`

	TotalCost    float64 `json:"total_cost"`
	TotalVAT     float64 `json:"total_vat"`
	BasketStatus int     `json:"basket_status"`

	BasketProducts []*BasketProduct `pg:"rel:has-many"`

	UserID string ` json:"user_id"`
	User   *User  `pg:"rel:has-one"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
