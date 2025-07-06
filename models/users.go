package models

import (
	"be-tickitz/utils"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID             int     `json:"idUser" db:"id"`
	Email          string  `json:"email"`
	Password       string  `json:"password"`
	FullName       string  `json:"fullName" db:"full_name,omitempty"`
	PhoneNumber    *string `json:"phoneNumber" db:"phone_number,omitempty"`
	ProfilePicture *string `json:"profilePicture" db:"profile_picture"`
	Role           string  `json:"role" db:"role"`
}

type UserLogin struct {
	ID       int    `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Role     string `db:"role"`
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

	if user.Role == "" {
		user.Role = "user"
	}

	hashedPassword, err := utils.HashString(user.Password)
	if err != nil {
		return user, err
	}
	user.Password = hashedPassword

	err = conn.QueryRow(
		context.Background(),
		`
		INSERT INTO users (email, password, full_name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, full_name, role
		`,
		user.Email, user.Password, user.FullName, user.Role,
	).Scan(&user.ID, &user.Email, &user.FullName, &user.Role)

	user.Password = ""

	return user, err
}

func FindOneUserByEmail(email string) (UserLogin, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return UserLogin{}, err
	}
	defer conn.Release()

	rows, err := conn.Query(
		context.Background(),
		`
    SELECT id, email, password, role
    FROM users
    WHERE email = $1
    `,
		email,
	)
	if err != nil {
		return UserLogin{}, err
	}

	user, err := pgx.CollectOneRow[UserLogin](rows, pgx.RowToStructByName)
	return user, err
}

func UpdateUserPassword(userID int, newPassword string) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Release()

	hashedPassword, err := utils.HashString(newPassword)
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(),
		`UPDATE users SET password = $1 WHERE id = $2`,
		hashedPassword, userID,
	)

	return err
}

func GetAllUsers() ([]User, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `
    SELECT id, email, full_name, phone_number, profile_picture, role
    FROM users
    ORDER BY id ASC
  `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.FullName,
			&u.PhoneNumber,
			&u.ProfilePicture,
			&u.Role,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func GetUserByID(userID int) (User, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return User{}, err
	}
	defer conn.Release()

	var u User
	err = conn.QueryRow(context.Background(), `
    SELECT id, email, full_name, phone_number, profile_picture, role
    FROM users
    WHERE id = $1
  `, userID).Scan(
		&u.ID,
		&u.Email,
		&u.FullName,
		&u.PhoneNumber,
		&u.ProfilePicture,
		&u.Role,
	)
	return u, err
}

func DeleteUserByID(userID int) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Release()

	var exists bool
	err = conn.QueryRow(context.Background(),
		`SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`, userID,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user not found")
	}

	_, err = conn.Exec(context.Background(),
		`DELETE FROM users WHERE id = $1`, userID,
	)
	return err
}
