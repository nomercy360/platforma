package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// type ICartItem = {
//  id: number;
//  name: string;
//  price: number;
//  quantity: number;
//  size: string;
//};

type Cart struct {
	ID         int64           `json:"id" db:"id"`
	Items      []LineItem      `json:"items" db:"items"`
	CustomerID *int            `json:"customer_id" db:"customer_id"`
	Currency   string          `json:"currency" db:"currency"`
	Total      int             `json:"total" db:"total"`
	Count      int             `json:"count" db:"count"`
	Subtotal   int             `json:"subtotal" db:"subtotal"`
	Discount   *Discount       `json:"discount" db:"discount"`
	DiscountID *int64          `json:"discount_id" db:"discount_id"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time      `json:"deleted_at" db:"deleted_at"`
	Context    CustomerContext `json:"context" db:"context"`
}

type CustomerContext struct {
	UserAgent string `json:"user_agent" db:"user_agent"`
	IP        string `json:"ip" db:"ip"`
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
	q := `
		SELECT 
			c.id,
			c.customer_id,
			c.created_at,
			c.updated_at,
			c.deleted_at,
			c.context,
			COALESCE(SUM(p.price * li.quantity), 0) AS subtotal,
			COALESCE(SUM(p.price * li.quantity), 0) AS total,
			COALESCE(SUM(li.quantity), 0) AS count,
			p.currency,
			c.discount_id
		FROM
			cart c
		LEFT JOIN
			line_items li ON c.id = li.cart_id
		LEFT JOIN
			product_variants pv ON li.variant_id = pv.id
		LEFT JOIN
			product_prices p ON pv.id = p.product_id
		WHERE
			c.id = ?
		GROUP BY
			c.id, p.currency;`

	row := s.db.QueryRow(q, id)

	err := row.Scan(
		&cart.ID,
		&cart.CustomerID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
		&cart.DeletedAt,
		&cart.Context,
		&cart.Subtotal,
		&cart.Total,
		&cart.Count,
		&cart.Currency,
		&cart.DiscountID,
	)

	if err != nil && IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	items, err := s.GetLineItemsByCartID(id, "en")

	if err != nil {
		return nil, err
	}

	cart.Items = items

	if cart.DiscountID != nil {
		discount, err := s.GetDiscount(DiscountQuery{ID: *cart.DiscountID})
		if err != nil {
			return nil, err
		}

		cart.Discount = discount

		if discount.Value > 0 {
			switch discount.Type {
			case "percentage":
				cart.Total = cart.Total - (cart.Total * discount.Value / 100)
			case "fixed":
				cart.Total = cart.Total - discount.Value
			}
		}
	}

	return &cart, nil
}

func (s Storage) GetLineItemsByCartID(cartID int64, locale string) ([]LineItem, error) {
	q := `
		SELECT li.id,
			   li.cart_id,
			   li.order_id,
			   li.variant_id,
			   li.quantity,
			   li.created_at,
			   li.updated_at,
			   li.deleted_at,
			   pv.name AS variant_name,
			   p.name AS product_name,
			   p.cover_image_url AS image_url,
			   pp.price
		FROM line_items li
		JOIN main.product_variants pv on li.variant_id = pv.id
		JOIN main.products p on pv.product_id = p.id
		JOIN product_prices pp on pv.id = pp.product_id AND pp.currency = ?
		WHERE li.cart_id = ?
		  AND li.deleted_at IS NULL;`

	var currency string
	switch locale {
	case "en":
		currency = "USD"
	case "ru", "by":
		currency = "BYN"
	}

	rows, err := s.db.Query(q, currency, cartID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []LineItem
	for rows.Next() {
		var item LineItem
		if err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.OrderID,
			&item.VariantID,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
			&item.VariantName,
			&item.ProductName,
			&item.ImageURL,
			&item.Price,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (s Storage) CreateCart(cart Cart) (*Cart, error) {
	res, err := s.db.Exec("INSERT INTO cart (customer_id, context) VALUES (?, ?)",
		cart.CustomerID, cart.Context)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	for _, item := range cart.Items {
		_, err := s.db.Exec("INSERT INTO line_items (cart_id, variant_id, quantity) VALUES (?, ?, ?)",
			id, item.VariantID, item.Quantity)

		if err != nil {
			return nil, err
		}
	}

	return s.GetCartByID(id)
}

func (s Storage) DeleteCart(id int64) error {
	_, err := s.db.Exec("UPDATE cart SET deleted_at = ? WHERE id = ?", time.Now(), id)

	return err
}

func (s Storage) UpdateCartDiscount(cartID int64, discountID int64) error {
	_, err := s.db.Exec("UPDATE cart SET discount_id = ? WHERE id = ?", discountID, cartID)

	return err
}

func (s Storage) DropCartDiscount(cartID int64) error {
	_, err := s.db.Exec("UPDATE cart SET discount_id = NULL WHERE id = ?", cartID)

	return err
}
