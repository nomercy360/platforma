package store

import (
	"github.com/labstack/echo/v4"
	"github.com/plutov/paypal/v4"
	"log"
	"rednit/config"
	"rednit/db"
	"rednit/payment"
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
	ListProducts(params db.ListProductsQuery) ([]db.Product, error)
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

func langFromContext(c echo.Context) string {
	return c.Get("lang").(string)
}

func (h Handler) Debug(c echo.Context) error {
	for name, headers := range c.Request().Header {
		for _, h := range headers {
			log.Printf("Header '%v': '%v'\n", name, h)
		}
	}

	return c.JSON(200, "OK")
}
