package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Cart struct {
	ID             int64           `json:"id" db:"id"`
	Items          []LineItem      `json:"items" db:"items"`
	CustomerID     *int            `json:"customer_id" db:"customer_id"`
	CurrencyCode   string          `json:"currency_code" db:"currency_code"`
	CurrencySymbol string          `json:"currency_symbol" db:"currency_symbol"`
	Total          int             `json:"total" db:"total"`
	Count          int             `json:"count" db:"count"`
	Subtotal       int             `json:"subtotal" db:"subtotal"`
	Discount       *Discount       `json:"discount" db:"discount"`
	DiscountID     *int64          `json:"discount_id" db:"discount_id"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time      `json:"deleted_at" db:"deleted_at"`
	Context        CustomerContext `json:"context" db:"context"`
	DiscountAmount int             `json:"discount_amount" db:"-"`
}

type CustomerContext struct {
	UserAgent string `json:"user_agent" db:"user_agent"`
	IP        string `json:"ip" db:"ip"`
}

func currencyFromLocale(locale string) string {
	switch locale {
	case "en":
		return "USD"
	case "ru", "by":
		return "BYN"
	default:
		return "USD"
	}
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

func (s Storage) GetCartByID(id int64, locale, currency string) (*Cart, error) {
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
			COALESCE(p.currency_code, ?) AS currency,
			COALESCE(cr.symbol, '$') AS currency_symbol,
			c.discount_id
		FROM
			cart c
		LEFT JOIN
			line_items li ON c.id = li.cart_id
		LEFT JOIN
			product_variants pv ON li.variant_id = pv.id
		LEFT JOIN
			product_prices p ON pv.product_id = p.product_id AND p.currency_code = ?
		LEFT JOIN
			currencies cr ON p.currency_code = cr.code
		WHERE
			c.id = ?
		GROUP BY
			c.id, p.currency_code;`

	row := s.db.QueryRow(q, currency, currency, id)

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
		&cart.CurrencyCode,
		&cart.CurrencySymbol,
		&cart.DiscountID,
	)

	if err != nil && IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	items, err := s.GetLineItems(LineItemQuery{Locale: locale, CartID: id})

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
				cart.DiscountAmount = cart.Subtotal - cart.Total
			case "fixed":
				cart.Total = cart.Total - discount.Value
				cart.DiscountAmount = discount.Value
			}
		}
	}

	// only for testing purposes
	if len(cart.Items) == 1 && cart.Items[0].Price == 1 {
		cart.Total = 1
		cart.Subtotal = 1
	} else {
		// delivery
		if cart.CurrencyCode == "USD" {
			cart.Total += 10
		} else {
			cart.Total += 25
		}
	}

	return &cart, nil
}

func lineItemQuery(locale string) string {
	switch locale {
	case "ru", "by":
		return `
			SELECT li.id,
				   li.cart_id,
				   li.order_id,
				   li.variant_id,
				   li.quantity,
				   li.created_at,
				   li.updated_at,
				   li.deleted_at,
				   pv.name AS variant_name,
				   pt.name AS product_name,
				   p.cover_image_url AS image_url,
				   pp.price
			FROM line_items li
			JOIN main.product_variants pv on li.variant_id = pv.id
			JOIN main.products p on pv.product_id = p.id
			JOIN product_prices pp on p.id = pp.product_id AND pp.currency_code = ?
			JOIN product_translations pt on p.id = pt.product_id AND pt.language = 'ru'`
	default:
		return `
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
		JOIN product_variants pv on li.variant_id = pv.id
		JOIN products p on pv.product_id = p.id
		JOIN product_prices pp on p.id = pp.product_id AND pp.currency_code = ?`
	}
}

type LineItemQuery struct {
	Locale   string
	CartID   int64
	OrderID  int64
	Currency string
}

func (s Storage) GetLineItems(query LineItemQuery) ([]LineItem, error) {
	q := lineItemQuery(query.Locale)

	var currency string
	if query.Currency != "" {
		currency = query.Currency
	} else {
		currency = currencyFromLocale(query.Locale)
	}

	args := []interface{}{currency}

	if query.CartID > 0 {
		q = fmt.Sprintf("%s WHERE li.cart_id = %d", q, query.CartID)
		args = append(args, query.CartID)
	} else if query.OrderID > 0 {
		q = fmt.Sprintf("%s WHERE li.order_id = %d", q, query.OrderID)
		args = append(args, query.OrderID)
	}

	rows, err := s.db.Query(q, args...)
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

func (s Storage) CreateCart(cart Cart, locale string) (*Cart, error) {
	res, err := s.db.Exec("INSERT INTO cart (customer_id, context) VALUES (?, ?)",
		cart.CustomerID, cart.Context)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	if len(cart.Items) > 0 {
		for _, item := range cart.Items {
			li := LineItem{
				CartID:    &id,
				VariantID: item.VariantID,
				Quantity:  item.Quantity,
			}

			if err := s.SaveLineItem(li); err != nil {
				return nil, err
			}
		}
	}

	return s.GetCartByID(id, locale, currencyFromLocale(locale))
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

func (s Storage) UpdateCartCustomer(cartID int64, customerID int64) error {
	_, err := s.db.Exec("UPDATE cart SET customer_id = ? WHERE id = ?", customerID, cartID)

	return err
}
