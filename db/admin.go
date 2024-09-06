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
	Role      string     `db:"role" json:"role"`
	Name      *string    `db:"name" json:"name"`
}

type UserQuery struct {
	Email string
	ID    int64
}

func (s Storage) GetUser(query UserQuery) (*User, error) {
	var user User
	var args interface{}

	q := `
		SELECT id, email, name, password, created_at, updated_at, deleted_at
		FROM users
	`

	if query.Email != "" {
		q += " WHERE email = ?"
		args = query.Email
	} else {
		q += " WHERE id = ?"
		args = query.ID
	}

	err := s.db.QueryRow(q, args).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
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

func (s Storage) CreateUser(email, hashedPassword string, name *string) (*User, error) {
	query := `
		INSERT INTO users (email, name, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	var user User

	if err := s.db.QueryRow(query, email, name, hashedPassword, now, now).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	user.Email = email
	return &user, nil
}

func (s Storage) ListUsers() ([]User, error) {
	query := `
		SELECT id, email, created_at, updated_at, deleted_at, role, name
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
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.Role,
			&user.Name,
		); err != nil {
			return nil, err
		}

		users = append(users, user)

	}

	return users, nil
}
