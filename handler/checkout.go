package handler

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/plutov/paypal/v4"
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
	CartID          int64                  `json:"cart_id" validate:"required"`
	PaymentProvider string                 `json:"payment_provider" validate:"required,oneof=bepaid paypal"`
	Name            string                 `json:"name" validate:"required"`
	CustomerID      int64                  `json:"customer_id" validate:"required"`
	Phone           string                 `json:"phone" validate:"required"`
	Country         string                 `json:"country" validate:"required"`
	Address         string                 `json:"address" validate:"required"`
	ZIP             string                 `json:"zip" validate:"required"`
	PromoCode       *string                `json:"promo_code"`
	Metadata        map[string]interface{} `json:"metadata"`
}

var (
	PaymentProviderBePaid = "bepaid"
	PaymentProviderPayPal = "paypal"
)

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

	currencyCode := "BYN"
	if req.PaymentProvider == PaymentProviderPayPal {
		currencyCode = "USD"
	}

	newOrder := db.Order{
		CustomerID:      customer.ID,
		Status:          "created",
		PaymentStatus:   "pending",
		ShippingStatus:  "pending",
		Metadata:        req.Metadata,
		CartID:          cart.ID,
		Total:           cart.Total,
		Subtotal:        cart.Subtotal,
		CurrencyCode:    currencyCode,
		PaymentProvider: req.PaymentProvider,
	}

	order, err := h.st.CreateOrder(newOrder)

	if err != nil {
		return terrors.InternalServerError(err, "failed to create order")
	}

	if err := h.st.UpdateLineItemsOrderID(cart.ID, order.ID); err != nil {
		return terrors.InternalServerError(err, "failed to update line items order id")
	}

	var paymentLink string
	if req.PaymentProvider == PaymentProviderBePaid {
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
					Description: order.ToString(),
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
		paymentLink = tokenResp.Checkout.RedirectUrl
	} else if req.PaymentProvider == PaymentProviderPayPal {
		// Populate PayPal-specific fields
		paypalRequest := payment.PayPalRequest{
			PurchaseUnits: []paypal.PurchaseUnitRequest{
				{
					Amount: &paypal.PurchaseUnitAmount{
						Currency: order.CurrencyCode,
						Value:    strconv.Itoa(order.Total),
					},
					Description: order.ToString(),
					CustomID:    strconv.FormatInt(order.ID, 10), // Order ID as tracking ID
				},
			},
			ApplicationContext: &paypal.ApplicationContext{
				BrandName:   "PLUM<3",
				LandingPage: "BILLING",
				UserAction:  "PAY_NOW",
				ReturnURL:   fmt.Sprintf("%s/en/orders?orderId=%d", h.config.WebURL, order.ID),
				CancelURL:   fmt.Sprintf("%s/en/orders/cancel?orderId=%d", h.config.WebURL, order.ID),
			},
			Payer: &paypal.Payer{
				PayerInfo: &paypal.PayerInfo{
					Email:       customer.Email,
					FirstName:   req.Name,
					Phone:       req.Phone,
					CountryCode: req.Country,
					PayerID:     strconv.FormatInt(customer.ID, 10),
					ShippingAddress: &paypal.ShippingAddress{
						RecipientName: req.Name,
						Line1:         req.Address,
						PostalCode:    req.ZIP,
						CountryCode:   req.Country,
					},
				},
			},
		}

		// Create PayPal order
		paypalResp, err := h.paypal.CreatePaypalOrder(paypalRequest)
		if err != nil {
			return terrors.InternalServerError(err, "failed to create PayPal payment")
		}
		id := paypalResp.ID
		// save payment id
		order.PaymentID = &id
		order, err = h.st.UpdateOrder(order)
		if err != nil {
			return terrors.InternalServerError(err, "failed to update order")
		}

	} else {
		return terrors.BadRequest(errors.New("unsupported payment provider"), "unsupported payment provider")
	}

	cr := CheckoutResponse{
		Order:       *order,
		PaymentLink: paymentLink,
	}

	return c.JSON(http.StatusCreated, cr)
}

type CapturePaypalPaymentRequest struct {
	OrderID string `json:"order_id" validate:"required"`
}

func (h Handler) CapturePaypalPayment(c echo.Context) error {
	var req CapturePaypalPaymentRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := h.paypal.CapturePaypalOrder(req.OrderID)
	if err != nil {
		return terrors.InternalServerError(err, "failed to capture PayPal payment")
	}

	log.Infof("PayPal payment captured: %v", resp)

	if resp.Status != "COMPLETED" {
		return terrors.BadRequest(errors.New("payment not completed"), "payment not completed")
	}

	order, err := h.st.GetOrder(db.GetOrderQuery{PaymentID: &req.OrderID})

	if err != nil {
		return terrors.InternalServerError(err, "failed to get order")
	}

	order.PaymentStatus = "paid"

	if _, err = h.st.UpdateOrder(order); err != nil {
		return terrors.InternalServerError(err, "failed to update order")
	}

	go func() {
		h.telegramOrderPaid(*order)
	}()

	return c.JSON(http.StatusOK, order)
}
