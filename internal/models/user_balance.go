package models

import "time"

type UserBalance struct {
	ID        int64   	`json:"id"`
	UserID    int64   	`json:"user_id"`
	Balance   float64 	`json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}