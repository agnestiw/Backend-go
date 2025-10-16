package repository

import (
	"latihan2/app/model"
	"latihan2/database"
)

func GetUserByUsername(username string) (*model.User, string, error) {
	var user model.User
	var passwordHash string

	row := database.DB.QueryRow(
		`SELECT id, username, email, password_hash, role, created_at
		 FROM users
		 WHERE username = $1 OR email = $1`,
		username,
	)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &passwordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, "", err
	}

	return &user, passwordHash, nil
}