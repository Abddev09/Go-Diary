package main

import (
	"time"
)

type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type Post struct {
	ID             int       `json:"id" db:"id"`
	Title          string    `json:"title" db:"title"`
	Description    string    `json:"description" db:"description"`
	ImageURL       *string   `json:"image_url" db:"image_url"`
	ImageFile      *string   `json:"image_file" db:"image_file"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	AuthorID       int       `json:"author_id" db:"author_id"`
	AuthorUsername string    `json:"author_username" db:"author_username"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PostRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

type AuthResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Message  string `json:"message"`
}
