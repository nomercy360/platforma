package db

import (
	"time"
)

type User struct {
	ID        int64      `db:"id" json:"id"`
	Email     string     `db:"email" json:"email"`
	Password  string     `db:"password" json:"-"` // Don't expose password in JSON responses
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	AvatarURL string     `db:"avatar_url" json:"avatar_url"`
	Role      string     `db:"role" json:"role"`
	Name      *string    `db:"name" json:"name"`
}

func (s Storage) getUserByQuery(query string, args ...interface{}) (*User, error) {
	var user User

	err := s.db.QueryRow(query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.AvatarURL,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s Storage) GetUserByEmail(email string) (*User, error) {
	query := `
		SELECT id, email, name, avatar_url, password, created_at, updated_at, deleted_at
		FROM users WHERE email = ?
	`
	return s.getUserByQuery(query, email)
}

func (s Storage) GetUserByID(id int64) (*User, error) {
	query := `
		SELECT id, email, name, avatar_url, password, created_at, updated_at, deleted_at
		FROM users WHERE id = ?
	`
	return s.getUserByQuery(query, id)
}

func (s Storage) CreateUser(user User) (*User, error) {
	query := `
		INSERT INTO users (email, name, avatar_url, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	res, err := s.db.Exec(query, user.Email, user.Name, user.AvatarURL, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if id, err := res.LastInsertId(); err != nil {
		return nil, err
	} else {
		return s.GetUserByID(id)
	}
}

func (s Storage) ListUsers() ([]User, error) {
	query := `
		SELECT id, email, name, role, avatar_url, created_at, updated_at, deleted_at
		FROM users
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]User, 0)

	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Role,
			&user.AvatarURL,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)

	}

	return users, nil
}
