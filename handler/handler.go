package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/plutov/paypal/v4"
	"rednit/config"
	"rednit/payment"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"rednit/db"
)

type Handler struct {
	st     storage
	config config.Default
	paypal paymentPaypal
}

func New(st storage, config config.Default, p paymentPaypal) Handler {
	return Handler{st: st, config: config, paypal: p}
}

type paymentPaypal interface {
	CreatePaypalOrder(request payment.PayPalRequest) (*paypal.Order, error)
	CapturePaypalOrder(orderID string) (*paypal.CaptureOrderResponse, error)
}

type storage interface {
	ListProducts(locale string) ([]db.Product, error)
	GetProduct(query db.GetProductQuery) (*db.Product, error)
	CreateCart(cart db.Cart, lang string) (*db.Cart, error)
	GetCartByID(cartID int64, locale string) (*db.Cart, error)
	SaveLineItem(li db.LineItem) error
	GetCustomerByEmail(email string) (*db.Customer, error)
	GetCustomerByID(id int64) (*db.Customer, error)
	AddCustomer(c db.Customer) (*db.Customer, error)
	CreateOrder(o db.Order) (*db.Order, error)
	GetDiscount(query db.DiscountQuery) (*db.Discount, error)
	UpdateDiscountUsageCount(id int64) error
	UpdateOrder(o *db.Order) (*db.Order, error)
	GetOrder(query db.GetOrderQuery) (*db.Order, error)
	UpdateLineItemsOrderID(cartID, orderID int64) error
	UpdateCartDiscount(cartID, discountID int64) error
	DropCartDiscount(cartID int64) error
	UpdateLineItemQuantity(li int64, quantity int) error
	RemoveLineItem(li int64) error
	UpdateCustomer(c *db.Customer) (*db.Customer, error)
	UpdateCartCustomer(cartID int64, customerID int64) error
	UpdateCartCurrency(cartID int64, currency string) error
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
