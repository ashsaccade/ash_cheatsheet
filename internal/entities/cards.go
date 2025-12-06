package entities

import "time"

type Card struct {
	Id          int64
	Name        string
	Description string
	Section     string
	UpdatedAt   time.Time
}
