package db

import "time"

type Customer struct {
	ID        int64      `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Email     string     `db:"email" json:"email"`
	Phone     string     `db:"phone" json:"phone"`
	Country   string     `db:"country" json:"country"`
	Address   string     `db:"address" json:"address"`
	ZIP       string     `db:"zip" json:"zip"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

func (s Storage) GetCustomerByEmail(email string) (*Customer, error) {
	c := new(Customer)

	query := `SELECT * FROM customers WHERE email = ?`
	row := s.db.QueryRow(query, email)

	err := row.Scan(
		&c.ID,
		&c.Name,
		&c.Email,
		&c.Phone,
		&c.Country,
		&c.Address,
		&c.ZIP,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.DeletedAt,
	)

	if err != nil && IsNoRowsError(err) {
		return c, ErrNotFound
	} else if err != nil {
		return c, err
	}

	return c, nil
}

func (s Storage) GetCustomerByID(id int64) (*Customer, error) {
	c := new(Customer)

	query := `SELECT * FROM customers WHERE id = ?`
	row := s.db.QueryRow(query, id)

	err := row.Scan(
		&c.ID,
		&c.Name,
		&c.Email,
		&c.Phone,
		&c.Country,
		&c.Address,
		&c.ZIP,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.DeletedAt,
	)

	if err != nil && IsNoRowsError(err) {
		return c, ErrNotFound
	} else if err != nil {
		return c, err
	}

	return c, nil
}

func (s Storage) SaveCustomer(c Customer) (*Customer, error) {
	query := `
		INSERT INTO customers (name, email, phone, country, address, zip)
		VALUES (?, ?, ?, ?, ?, ?);
	`

	res, err := s.db.Exec(query, c.Name, c.Email, c.Phone, c.Country, c.Address, c.ZIP)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	return s.GetCustomerByID(id)
}
