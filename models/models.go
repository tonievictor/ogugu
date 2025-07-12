package models

import "time"

type Subscription struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	RSS       RssFeed   `json:"rss"`
}

type CreateSubscription struct {
	RssID string `json:"rss_id"`
}

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
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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

type RSSMeta struct {
	Channel struct {
		LastBuildDate string `xml:"lastBuildDate" validate:"required,datetime"`
		Title         string `xml:"title" validate:"required"`
		Description   string `xml:"description" validate:"required"`
	} `xml:"channel"`
}
