package entities

import (
	"time"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Baskets []*Basket `pg:"rel:has-many"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
