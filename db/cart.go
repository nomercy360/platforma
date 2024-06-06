package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type CartItem struct {
	VariantID   int    `json:"variant_id"`
	Quantity    int    `json:"quantity"`
	VariantName string `json:"variant_name"`
	Price       int    `json:"price"`
}

type CartItemInput struct {
	VariantID int `json:"variant_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,min=1"`
}

type CartItems []CartItem

type CartItemsInput []CartItemInput

type Cart struct {
	ID         int64           `json:"id" db:"id"`
	Items      CartItems       `json:"items" db:"items"`
	CustomerID int             `json:"customer_id" db:"customer_id"`
	Total      int             `json:"total" db:"total"`
	Subtotal   int             `json:"subtotal" db:"subtotal"`
	Discount   int             `json:"discount" db:"discount"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time      `json:"deleted_at" db:"deleted_at"`
	Context    CustomerContext `json:"context" db:"context"`
}

type CustomerContext struct {
	UserAgent string `json:"user_agent" db:"user_agent"`
	IP        string `json:"ip" db:"ip"`
}

func (ci CartItemsInput) Value() (driver.Value, error) {
	if len(ci) == 0 {
		return nil, nil
	}

	b, err := json.Marshal(ci)

	if err != nil {
		return nil, err
	}

	return string(b), nil
}

func (ci *CartItems) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		*ci = make([]CartItem, 0)
	case []byte:
		if err := json.Unmarshal(src, ci); err != nil {
			return err
		}
	case string:
		if err := json.Unmarshal([]byte(src), ci); err != nil {
			return err
		}
	default:
		return errors.New("unsupported type")
	}

	return nil
}

func (cc CustomerContext) Value() (driver.Value, error) {
	if cc == (CustomerContext{}) {
		return nil, nil
	}

	b, err := json.Marshal(cc)

	if err != nil {
		return nil, err
	}

	return string(b), nil
}

func (cc *CustomerContext) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		*cc = CustomerContext{}
	case []byte:
		if err := json.Unmarshal(src, cc); err != nil {
			return err
		}
	case string:
		if err := json.Unmarshal([]byte(src), cc); err != nil {
			return err
		}
	default:
		return errors.New("unsupported type")
	}

	return nil
}

func (s Storage) GetCartByID(id int64) (*Cart, error) {
	var cart Cart
	row := s.db.QueryRow("SELECT id, customer_id, items, total, subtotal, discount, created_at, updated_at, deleted_at, context FROM cart WHERE id = ?", id)

	err := row.Scan(
		&cart.ID,
		&cart.CustomerID,
		&cart.Items,
		&cart.Total,
		&cart.Subtotal,
		&cart.Discount,
		&cart.CreatedAt,
		&cart.UpdatedAt,
		&cart.DeletedAt,
		&cart.Context,
	)

	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (s Storage) CreateCart(cart Cart) (*Cart, error) {
	res, err := s.db.Exec("INSERT INTO cart (customer_id, items, total, subtotal, discount, context) VALUES (?, ?, ?, ?, ?, ?)",
		cart.CustomerID, cart.Items, cart.Total, cart.Subtotal, cart.Discount, cart.Context)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	return s.GetCartByID(id)
}
