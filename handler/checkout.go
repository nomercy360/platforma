package handler

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"rednit/db"
	"rednit/payment"
	"rednit/terrors"
	"strconv"
)

type CartItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
	VariantID int64 `json:"variant_id"`
}

type CheckoutRequest struct {
	CartID     int64                  `json:"cart_id" validate:"required"`
	Provider   string                 `json:"provider" validate:"required"`
	Name       string                 `json:"name" validate:"required"`
	CustomerID int64                  `json:"customer_id" validate:"required"`
	Phone      string                 `json:"phone" validate:"required"`
	Country    string                 `json:"country" validate:"required"`
	Address    string                 `json:"address" validate:"required"`
	ZIP        string                 `json:"zip" validate:"required"`
	PromoCode  *string                `json:"promo_code"`
	Metadata   map[string]interface{} `json:"metadata"`
}

type CheckoutResponse struct {
	Order       db.Order `json:"order"`
	PaymentLink string   `json:"payment_link"`
}

func (h Handler) Checkout(c echo.Context) error {
	var req CheckoutRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	// get locale from header
	locale := langFromContext(c)

	// locale := "ru"

	customer, err := h.st.GetCustomerByID(req.CustomerID)

	if err != nil {
		return terrors.InternalServerError(err, "failed to get customer")
	} else {
		// update customer info
		customer.Name = &req.Name
		customer.Phone = &req.Phone
		customer.Country = &req.Country
		customer.Address = &req.Address
		customer.ZIP = &req.ZIP

		customer, err = h.st.UpdateCustomer(customer)

		if err != nil {
			return terrors.InternalServerError(err, "failed to update customer")
		}

		log.Infof("Customer updated: %v", customer)
	}

	// Никита попросил использовать только BYN для Bepaid - USD не поддерживается
	cart, err := h.st.GetCartByID(req.CartID, locale)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "cart not found")
	} else if err != nil {
		return terrors.InternalServerError(err, "failed to get cart")
	}

	newOrder := db.Order{
		CustomerID:     customer.ID,
		Status:         "created",
		PaymentStatus:  "pending",
		ShippingStatus: "pending",
		Metadata:       req.Metadata,
		CartID:         cart.ID,
		Total:          cart.Total,
		Subtotal:       cart.Subtotal,
		CurrencyCode:   cart.CurrencyCode,
	}

	order, err := h.st.CreateOrder(newOrder)

	if err != nil {
		return terrors.InternalServerError(err, "failed to create order")
	}

	if err := h.st.UpdateLineItemsOrderID(cart.ID, order.ID); err != nil {
		return terrors.InternalServerError(err, "failed to update line items order id")
	}

	var itemsString string
	for _, item := range cart.Items {
		itemsString += fmt.Sprintf("%s(%s) x %d", item.ProductName, item.VariantName, item.Quantity)
		if item != cart.Items[len(cart.Items)-1] {
			itemsString += ", "
		}
	}

	paymentRequest := payment.BepaidTokenRequest{
		Checkout: payment.BepaidCheckout{
			Attempts:        1,
			Test:            h.config.Bepaid.TestMode,
			TransactionType: "payment",
			Settings: payment.BepaidSettings{
				NotificationUrl: fmt.Sprintf("%s/webhook/bepaid", h.config.ExternalURL),
				SuccessUrl:      fmt.Sprintf("%s/en/orders?orderId=%d", h.config.WebURL, order.ID),
				Language:        locale,
				AutoReturn:      "0",
				WidgetStyle: map[string]interface{}{
					"widget": map[string]interface{}{
						"backgroundColor": "#ffffff",
						"buttonsColor":    "#262626",
						"backgroundType":  "2",
						"color":           "#262626",
						"fontSize":        "15px",
						"fontWeight":      "400",
						"lineHeight":      "21px",
					},
					"inputs": map[string]interface{}{
						"backgroundColor": "#f8f8f8",
						"holder": map[string]interface{}{
							"backgroundColor": "#f8f8f8",
						},
					},
					"button": map[string]interface{}{
						"backgroundColor": "#262626",
						"pay": map[string]interface{}{
							"color": "#ffffff",
						},
						"card": map[string]interface{}{
							"color": "#ffffff",
						},
						"brands": map[string]interface{}{
							"color": "#ffffff",
						},
					},
				},
			},
			Order: payment.BepaidOrder{
				Amount:      order.Total * 100,
				Currency:    order.CurrencyCode,
				Description: fmt.Sprintf("#%d: %s", order.ID, itemsString),
				TrackingID:  strconv.FormatInt(order.ID, 10),
			},
			Customer: payment.BepaidCustomer{
				Email:     customer.Email,
				FirstName: *customer.Name,
				LastName:  *customer.Name,
				Address:   *customer.Address,
				ZIP:       *customer.ZIP,
				Country:   *customer.Country,
				Phone:     *customer.Phone,
			},
		},
	}

	tokenResp, err := payment.CreatePaymentToken(
		paymentRequest,
		fmt.Sprintf("%s/ctp/api/checkouts", h.config.Bepaid.ApiURL),
		h.config.Bepaid.ShopID,
		h.config.Bepaid.SecretKey,
	)

	if err != nil {
		return terrors.InternalServerError(err, "failed to create payment token")
	}

	log.Infof("Payment token created: %s", tokenResp.Checkout.RedirectUrl)

	cr := CheckoutResponse{
		Order:       *order,
		PaymentLink: tokenResp.Checkout.RedirectUrl,
	}

	return c.JSON(http.StatusCreated, cr)
}
