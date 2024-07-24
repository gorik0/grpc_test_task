package model

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
type UserDB struct {
	ID    string
	Email string
}
