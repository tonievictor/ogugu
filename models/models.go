package models

import "time"

type User struct {
	ID        string    `json:"id" validate:"required"`
	Username  string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required:email"`
	Avatar    string    `json:"avatar" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}

type RssFeed struct {
	ID        string    `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Link      string    `json:"link" validate:"required,url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateRssBody struct {
	Name string `json:"name" validate:"required"`
	Link string `json:"link" validate:"required,url"`
}

type CreateUserBody struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required:email"`
	Avatar   string `json:"avatar" validate:"omitempty,string"`
	Password string `json:"password" validate:"required,string"`
}
