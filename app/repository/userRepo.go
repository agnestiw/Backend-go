package repository

import (
	"database/sql"
	"fmt"
	"latihan2/app/model"
	"latihan2/database"
	"log"
)

func GetUsersRepo(search, sortBy, order, role string, limit, offset int) ([]model.User, error) {
	condition := "deleted_at IS NULL" // default untuk user biasa
	if role == "admin" {
		condition = "1=1" // admin bisa lihat semua
	}

	query := fmt.Sprintf(`
		SELECT id, username, email, role, created_at, deleted_at
		FROM users
		WHERE (%s) AND (username ILIKE $1 OR email ILIKE $1)
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, condition, sortBy, order)

	rows, err := database.DB.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		log.Println("Query error:", err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.CreatedAt, &u.DeletedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}


func CountUsersRepo(search string) (int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM users WHERE username ILIKE $1 OR email ILIKE $1`
	err := database.DB.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}

func SoftDeleteUserRepo(id int) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := database.DB.Exec(query, id)
	if err != nil {
		log.Println("Soft delete error:", err)
		return err
	}
	return nil
}

func GetUserByID(id int, role string) (*model.User, error) {
    condition := "deleted_at IS NULL" // default user biasa
    if role == "admin" {
        condition = "1=1" // admin bisa lihat semua
    }

    query := fmt.Sprintf(`
        SELECT id, username, email, role, created_at, deleted_at
        FROM users
        WHERE id = $1 AND %s
    `, condition)

    var u model.User
    row := database.DB.QueryRow(query, id)
    err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.CreatedAt, &u.DeletedAt)
    if err != nil {
        return nil, err
    }
    return &u, nil
}
