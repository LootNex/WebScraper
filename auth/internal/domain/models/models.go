package models

type User struct{
	ID string
	Login string
	TelegramLogin string
	PassHash []byte
}