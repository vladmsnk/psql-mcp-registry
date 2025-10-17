package api

import (
	"time"
)

// RegisterInstanceRequest represents the request body for registering a new instance
type RegisterInstanceRequest struct {
	Name            string `json:"name" binding:"required"`
	DatabaseName    string `json:"database_name" binding:"required"`
	Description     string `json:"description"`
	CreatorUsername string `json:"creator_username"`
}

// RegisterInstanceResponse represents the response after successful registration
type RegisterInstanceResponse struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	DatabaseName    string    `json:"database_name"`
	Description     string    `json:"description"`
	CreatorUsername string    `json:"creator_username"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ErrorResponse represents the standard error response format
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

