package models

import (
	"be-tickitz/utils"
	"context"
	"strings"
)

type User struct {
	ID             int    `json:"idUser" db:"id"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	FullName       string `json:"fullName" db:"full_name"`
	PhoneNumber    string `json:"phoneNumber" db:"phone_number"`
	ProfilePicture string `json:"profilePicture" db:"profile_picture"`
}

func Register(user User) (User, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return user, err
	}
	defer conn.Release()

	if strings.TrimSpace(user.FullName) == "" {
		user.FullName = utils.ExtractNameFromEmail(user.Email)
	}

	hashedPassword, err := utils.HashString(user.Password)
	if err != nil {
		return user, err
	}
	user.Password = hashedPassword

	err = conn.QueryRow(
		context.Background(),
		`
		INSERT INTO users (email, password, full_name)
		VALUES ($1, $2, $3)
		RETURNING id, email, full_name
		`,
		user.Email, user.Password, user.FullName,
	).Scan(&user.ID, &user.Email, &user.FullName)

	user.Password = ""

	return user, err
}
