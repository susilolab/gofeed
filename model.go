package main

import "time"

type User struct {
	ID                 int64     `json:"id"`
	Username           string    `json:"username"`
	AuthKey            string    `json:"auth_key"`
	PasswordHash       string    `json:"password_hash"`
	PasswordResetToken string    `json:"password_reset_token"`
	Email              string    `json:"email"`
	Status             int       `json:"status"`
	Github             string    `json:"github"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type UserLogin struct {
	Username string
	Password string
}

type News struct {
	ID        int64     `sql:"type:int";json:"id"`
	UserId    int64     `sql:"type:int;not null";json:"user_id"`
	Title     string    `sql:"type:varchar(50)";json:"title"`
	Text      string    `sql:"type:text";json:"text"`
	Link      string    `sql:"type:varchar(255)";json:"link"`
	Status    int       `sql:"type:tinyint(1)";json:"status"`
	CreatedAt time.Time `sql:"type:datetime";json:"created_at"`
}

type Auth struct {
	ID       int64  `json:"id"`
	UserId   int64  `json:"user_id"`
	Source   string `json:"source"`
	SourceId string `json:"source_id"`
}

type Comment struct {
	ID        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	NewsId    int64     `json:"news_id"`
	Text      string    `json:"text"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DbConfig struct {
	Hostname string `json:"hostname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     string `json:"port"`
	DbName   string `json:"db_name"`
}
