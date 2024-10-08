package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
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

type PaymentStatus string

const (
	PaymentPending    PaymentStatus = "pending"
	PaymentProcessing PaymentStatus = "processing"
	PaymentPaid       PaymentStatus = "paid"
	PaymentFailed     PaymentStatus = "failed"
	PaymentCanceled   PaymentStatus = "canceled"
	PaymentRefunded   PaymentStatus = "refunded"
)

var ValidPaymentStatuses = []PaymentStatus{
	PaymentPending, PaymentProcessing, PaymentPaid, PaymentFailed, PaymentCanceled,
	PaymentRefunded,
}

func (status PaymentStatus) IsValid() error {
	for _, v := range ValidPaymentStatuses {
		if v == status {
			return nil
		}
	}
	return errors.New("invalid payment status")
}

type OrderStatus string

const (
	OrderNew           OrderStatus = "new"
	OrderApproved      OrderStatus = "approved"
	OrderWaiting       OrderStatus = "waiting"
	OrderProduction    OrderStatus = "production"
	OrderAssembled     OrderStatus = "assembled"
	OrderReady         OrderStatus = "ready"
	OrderShipping      OrderStatus = "shipping"
	OrderShipped       OrderStatus = "shipped"
	OrderReadyToPickup OrderStatus = "ready_to_pickup"
	OrderCompleted     OrderStatus = "completed"
	OrderReturned      OrderStatus = "returned"
	OrderRefunded      OrderStatus = "refunded"
	OrderOutOfStock    OrderStatus = "out_of_stock"
	OrderCancelled     OrderStatus = "cancelled"
)

var ValidOrderStatuses = []OrderStatus{
	OrderNew, OrderApproved, OrderWaiting, OrderProduction, OrderAssembled, OrderReady, OrderShipping, OrderShipped,
	OrderReadyToPickup, OrderCompleted, OrderReturned, OrderRefunded, OrderOutOfStock, OrderCancelled,
}

func (status OrderStatus) IsValid() error {
	for _, v := range ValidOrderStatuses {
		if v == status {
			return nil
		}
	}
	return errors.New("invalid order status")
}

type Order struct {
	ID              int64         `db:"id" json:"id"`
	CustomerID      int64         `db:"customer_id" json:"customer_id"`
	CartID          int64         `db:"cart_id" json:"cart_id"`
	Status          OrderStatus   `db:"status" json:"status"`
	PaymentStatus   PaymentStatus `db:"payment_status" json:"payment_status"`
	Total           int           `db:"total" json:"total"`
	Subtotal        int           `db:"subtotal" json:"subtotal"`
	DiscountID      *int64        `db:"discount_id" json:"discount_id"`
	CurrencyCode    string        `db:"currency_code" json:"currency_code"`
	Metadata        Object        `db:"metadata" json:"metadata"`
	CreatedAt       time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time     `db:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time    `db:"deleted_at" json:"deleted_at"`
	PaymentID       *string       `db:"payment_id" json:"payment_id"`
	PaymentProvider string        `db:"payment_provider" json:"payment_provider"`
	Customer        *Customer     `json:"customer"`
	Items           []LineItem    `json:"items"`
}

func (o *Order) ToString() string {
	var itemsString string
	for _, item := range o.Items {
		itemsString += fmt.Sprintf("%s(%s) x %d", item.ProductName, item.VariantName, item.Quantity)
		if item != o.Items[len(o.Items)-1] {
			itemsString += ", "
		}
	}

	return fmt.Sprintf("#%d: %s", o.ID, itemsString)
}

type GetOrderQuery struct {
	ID        *int64
	PaymentID *string
}

func (s Storage) GetOrder(params GetOrderQuery) (*Order, error) {
	order := new(Order)

	query := `
		SELECT o.id,
			   o.customer_id,
			   o.cart_id,
			   o.discount_id,
			   o.status,
			   o.payment_status,
			   o.total,
			   o.subtotal,
			   o.created_at,
			   o.updated_at,
			   o.deleted_at,
			   o.currency_code,
			   o.metadata,
			   o.payment_id,
			   o.payment_provider
		FROM orders o`

	var args []interface{}
	if params.ID != nil {
		query += " WHERE o.id = ?"
		args = append(args, *params.ID)
	} else if params.PaymentID != nil {
		query += " WHERE o.payment_id = ?"
		args = append(args, *params.PaymentID)
	} else {
		return nil, errors.New("either ID or PaymentID must be provided")
	}

	row := s.db.QueryRow(query, args...)

	err := row.Scan(
		&order.ID,
		&order.CustomerID,
		&order.CartID,
		&order.DiscountID,
		&order.Status,
		&order.PaymentStatus,
		&order.Total,
		&order.Subtotal,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.DeletedAt,
		&order.CurrencyCode,
		&order.Metadata,
		&order.PaymentID,
		&order.PaymentProvider,
	)

	if err != nil && IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	order.Customer, err = s.GetCustomerByID(order.CustomerID)

	if err != nil {
		return nil, err
	}

	itemsParams := LineItemQuery{
		OrderID:  order.ID,
		Currency: order.CurrencyCode,
		Locale:   "en",
	}

	order.Items, err = s.GetLineItems(itemsParams)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s Storage) CreateOrder(o Order) (*Order, error) {
	query := `
		INSERT INTO orders (customer_id, cart_id, status, payment_status, total, subtotal, discount_id, currency_code, metadata, payment_id, payment_provider)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	res, err := s.db.Exec(query,
		o.CustomerID,
		o.CartID,
		o.Status,
		o.PaymentStatus,
		o.Total,
		o.Subtotal,
		o.DiscountID,
		o.CurrencyCode,
		o.Metadata,
		o.PaymentID,
		o.PaymentProvider,
	)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	return s.GetOrder(GetOrderQuery{ID: &id})
}

func (s Storage) UpdateOrder(o *Order) (*Order, error) {
	query := `
		UPDATE orders
		SET customer_id = ?, cart_id = ?, status = ?, payment_status = ?, total = ?, subtotal = ?, discount_id = ?, metadata = ?, payment_id = ?
		WHERE id = ?;
	`

	_, err := s.db.Exec(query, o.CustomerID, o.CartID, o.Status, o.PaymentStatus, o.Total, o.Subtotal, o.DiscountID, o.Metadata, o.PaymentID, o.ID)
	if err != nil {
		return nil, err
	}

	return s.GetOrder(GetOrderQuery{ID: &o.ID})
}

func (s Storage) ListOrders() ([]Order, error) {
	orders := make([]Order, 0)

	query := `
		SELECT o.id,
			   o.customer_id,
			   o.cart_id,
			   o.discount_id,
			   o.status,
			   o.payment_status,
			   o.total,
			   o.subtotal,
			   o.created_at,
			   o.updated_at,
			   o.deleted_at,
			   o.currency_code,
			   o.metadata,
			   o.payment_id,
			   o.payment_provider
		FROM orders o
		ORDER BY o.created_at DESC;
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		order := new(Order)
		err := rows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.CartID,
			&order.DiscountID,
			&order.Status,
			&order.PaymentStatus,
			&order.Total,
			&order.Subtotal,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.DeletedAt,
			&order.CurrencyCode,
			&order.Metadata,
			&order.PaymentID,
			&order.PaymentProvider,
		)

		if err != nil {
			return nil, err
		}

		order.Customer, err = s.GetCustomerByID(order.CustomerID)

		if err != nil {
			return nil, err
		}

		itemsParams := LineItemQuery{
			OrderID:  order.ID,
			Currency: order.CurrencyCode,
			Locale:   "en",
		}

		order.Items, err = s.GetLineItems(itemsParams)

		if err != nil {
			return nil, err
		}

		orders = append(orders, *order)
	}

	return orders, nil
}
