package entities

import (
	"time"
)

type User struct {
	UserID uint64 `json:"user_id"`

	HashedPassword string `json:"password"`
	Salt           string `json:"salt"`
	Hasher         string `json:"hasher"`

	IsBlacklisted bool      `json:"is_blacklisted"`
	Role         string      `json:"roles"`
	CreatedAt     time.Time `json:"created_at"`
	LastLogin     time.Time `json:"last_login"`
	Email         string    `json:"email"`
}

type GetUsersFilters struct {
	OffSet string
	Limit  string
	App    string
}
type UserUpdate struct {
	UserID        string    `json:"user_id"`
	ReferenceId   string    `json:"reference_id"`
	UserType      int64     `json:"type_id"`
	UserTypeName  string    `json:"-"`
	Name          string    `json:"name"`
	Phone         string    `json:"phone"`
	Password      string    `json:"password"`
	DefaultAppID  string    `json:"app_id"`
	DeviceID      string    `json:"device_id"`
	IsBlacklisted bool      `json:"is_blacklisted"`
	Roles         []string  `json:"roles"`
	CreatedAt     time.Time `json:"created_at"`
	LastLogin     time.Time `json:"last_login"`
	Email         string    `json:"email"`
}

type FetchUserResponse struct {
	UserId       int64         `json:"userId"`
	ReferenceId  string        `json:"referenceId"`
	UserTypeName string        `json:"userType"`
	Name         string        `json:"name"`
	Phone        string        `json:"phone"`
	RoleValues   []interface{} `json:"roles"`
	Email        string        `json:"email"`
}
type UserRole struct {
	UserID uint32
	RoleID uint32
}
