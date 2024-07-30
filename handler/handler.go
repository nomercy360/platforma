package handler

import (
	"github.com/labstack/echo/v4"
	"rednit/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"rednit/db"
)

type Handler struct {
	st     storage
	config config.Default
}

func New(st storage, config config.Default) Handler {
	return Handler{st: st, config: config}
}

type storage interface {
	ListProducts(locale string) ([]db.Product, error)
	GetProduct(query db.GetProductQuery) (*db.Product, error)
	CreateCart(cart db.Cart, locale string) (*db.Cart, error)
	GetCartByID(cartID int64, locale, currency string) (*db.Cart, error)
	SaveLineItem(li db.LineItem) error
	GetCustomerByEmail(email string) (*db.Customer, error)
	GetCustomerByID(id int64) (*db.Customer, error)
	SaveCustomer(c db.Customer) (*db.Customer, error)
	CreateOrder(o db.Order) (*db.Order, error)
	GetDiscount(query db.DiscountQuery) (*db.Discount, error)
	UpdateDiscountUsageCount(id int64) error
	UpdateOrder(o *db.Order) (*db.Order, error)
	GetOrderByID(id int64) (*db.Order, error)
	UpdateLineItemsOrderID(cartID, orderID int64) error
	UpdateCartDiscount(cartID, discountID int64) error
	DropCartDiscount(cartID int64) error
	UpdateLineItemQuantity(li int64, quantity int) error
	RemoveLineItem(li int64) error
	UpdateCustomer(c *db.Customer) (*db.Customer, error)
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UID    int64 `json:"uid"`
	ChatID int64 `json:"chat_id"`
}

func generateJWT(secret string, uid, chatID int64) (string, error) {
	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		UID:    uid,
		ChatID: chatID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func langFromContext(c echo.Context) string {
	return c.Get("lang").(string)
}
