package entity

import "time"

type Product struct {
	Id          int32
	Name        string
	Description string
	Sku         string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
