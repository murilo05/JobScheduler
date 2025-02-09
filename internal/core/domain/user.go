package domain

import (
	"time"
)

type UserRole string

const (
	Admin   UserRole = "admin"
	Cashier UserRole = "cashier"
)

type User struct {
	ID        uint64
	Name      string
	Email     string
	Password  string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}
