package model

import "time"

type User struct {
	Nickname    string
	Email       string
	Password    string
	DateOfBirth time.Time
}
