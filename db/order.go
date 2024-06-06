package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Object map[string]interface{}

func (o *Object) Scan(value interface{}) error {
	if value == nil {
		*o = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, o)
	case string:
		return json.Unmarshal([]byte(v), o)
	default:
		return errors.New("unsupported type for Object")
	}
}

func (o Object) Value() (driver.Value, error) {
	if o == nil {
		return nil, nil
	}
	return json.Marshal(o)
}

type Order struct {
	ID             int64      `db:"id" json:"id"`
	CustomerID     int64      `db:"customer_id" json:"customer_id"`
	CartID         int64      `db:"cart_id" json:"cart_id"`
	Status         string     `db:"status" json:"status"`
	PaymentStatus  string     `db:"payment_status" json:"payment_status"`
	ShippingStatus string     `db:"shipping_status" json:"shipping_status"`
	Total          int        `db:"total" json:"total"`
	Subtotal       int        `db:"subtotal" json:"subtotal"`
	DiscountID     *int64     `db:"discount_id" json:"discount_id"`
	Metadata       Object     `db:"metadata" json:"metadata"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deleted_at"`
}

func (s Storage) GetOrderByID(id int64) (*Order, error) {
	order := new(Order)

	query := `
		SELECT o.id,
			   o.customer_id,
			   o.cart_id,
			   o.discount_id,
			   o.status,
			   o.payment_status,
			   o.shipping_status,
			   o.total,
			   o.subtotal,
			   o.created_at,
			   o.updated_at,
			   o.deleted_at,
			   o.metadata
		FROM orders o
		WHERE o.id = ?;`

	row := s.db.QueryRow(query, id)

	err := row.Scan(
		&order.ID,
		&order.CustomerID,
		&order.CartID,
		&order.DiscountID,
		&order.Status,
		&order.PaymentStatus,
		&order.ShippingStatus,
		&order.Total,
		&order.Subtotal,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.DeletedAt,
		&order.Metadata,
	)

	if err != nil && IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return order, nil
}

func (s Storage) CreateOrder(o Order) (*Order, error) {
	query := `
		INSERT INTO orders (customer_id, cart_id, status, payment_status, shipping_status, total, subtotal, discount_id, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	res, err := s.db.Exec(query, o.CustomerID, o.CartID, o.Status, o.PaymentStatus, o.ShippingStatus, o.Total, o.Subtotal, o.DiscountID, o.Metadata)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	return s.GetOrderByID(id)
}

func (s Storage) UpdateOrder(o *Order) (*Order, error) {
	query := `
		UPDATE orders
		SET customer_id = ?, cart_id = ?, status = ?, payment_status = ?, shipping_status = ?, total = ?, subtotal = ?, discount_id = ?, metadata = ?
		WHERE id = ?;
	`

	_, err := s.db.Exec(query, o.CustomerID, o.CartID, o.Status, o.PaymentStatus, o.ShippingStatus, o.Total, o.Subtotal, o.DiscountID, o.Metadata, o.ID)
	if err != nil {
		return nil, err
	}

	return s.GetOrderByID(o.ID)
}
