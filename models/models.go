package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserWithAuth struct {
	User      User   `json:"user"`
	AuthToken string `json:"auth_token"`
}

type RssFeed struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Link      string    `json:"link"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Post struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	PubDate     time.Time `json:"pubDate"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreatePost struct {
	Title       string    `xml:"title" validate:"required"`
	Description string    `xml:"description" validate:"required"`
	Link        string    `xml:"link" validate:"required,url"`
	PubDate     time.Time `xml:"pubDate" validate:"required,datetime"`
}

type CreateRssBody struct {
	Name string `json:"name" validate:"required"`
	Link string `json:"link" validate:"required,url"`
}

type CreateUserBody struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Avatar   string `json:"avatar" validate:"omitempty"`
	Password string `json:"password" validate:"required,max=75"`
}

type SigninBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=75"`
}

type Session struct {
	UserID     string
	CreatedAt  time.Time
	ExpiryTime time.Time
}
