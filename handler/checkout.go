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
	CartID    int64                  `json:"cart_id"`
	Provider  string                 `json:"provider" validate:"required"`
	Name      string                 `json:"name" validate:"required"`
	Email     string                 `json:"email" validate:"required,email"`
	Phone     string                 `json:"phone" validate:"required"`
	Country   string                 `json:"country" validate:"required"`
	Address   string                 `json:"address" validate:"required"`
	ZIP       string                 `json:"zip" validate:"required"`
	PromoCode *string                `json:"promo_code"`
	Metadata  map[string]interface{} `json:"metadata"`
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

	customer, err := h.st.GetCustomerByEmail(req.Email)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		newCustomer := db.Customer{
			Name:    req.Name,
			Email:   req.Email,
			Phone:   req.Phone,
			Country: req.Country,
			Address: req.Address,
			ZIP:     req.ZIP,
		}

		customer, err = h.st.SaveCustomer(newCustomer)

		if err != nil {
			return terrors.InternalServerError(err, "failed to save customer")
		}
	} else if err != nil {
		return terrors.InternalServerError(err, "failed to get customer")
	}
	cart, err := h.st.GetCartByID(req.CartID)

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
		Currency:       cart.Currency,
	}

	order, err := h.st.CreateOrder(newOrder)

	if err != nil {
		return terrors.InternalServerError(err, "failed to create order")
	}

	if err := h.st.UpdateLineItemsOrderID(cart.ID, order.ID); err != nil {
		return terrors.InternalServerError(err, "failed to update line items order id")
	}

	//if req.PromoCode != nil {
	//	disc, err := h.st.GetDiscountByCode(*req.PromoCode)
	//
	//	if err != nil && errors.Is(err, db.ErrNotFound) {
	//		return terrors.NotFound(err, "discount not found")
	//	} else if err != nil {
	//		return terrors.InternalServerError(err, "failed to get discount")
	//	}
	//
	//	// calculate total
	//	if disc.Value > 0 {
	//		switch disc.Type {
	//		case "percentage":
	//			total = total - (total * disc.Value / 100)
	//		case "fixed":
	//			total = total - disc.Value
	//		}
	//
	//		// update usage count
	//		if err := h.st.UpdateDiscountUsageCount(disc.ID); err != nil {
	//			return terrors.InternalServerError(err, "failed to update discount usage count")
	//		}
	//
	//		// save order discount
	//		order.DiscountID = &disc.ID
	//	}
	//}
	//
	//order.Total = total
	//order.Subtotal = subtotal

	paymentRequest := payment.BepaidTokenRequest{
		Checkout: payment.BepaidCheckout{
			Attempts:        1,
			Test:            true,
			TransactionType: "payment",
			Settings: payment.BepaidSettings{
				NotificationUrl: "https://d421-125-24-110-63.ngrok-free.app/webhook/bepaid",
			},
			Order: payment.BepaidOrder{
				Amount:      order.Total * 100,
				Currency:    order.Currency,
				Description: fmt.Sprintf("Order #%d", order.ID),
				TrackingID:  strconv.FormatInt(order.ID, 10),
			},
			Customer: payment.BepaidCustomer{
				Email:     req.Email,
				FirstName: req.Name,
				LastName:  req.Name,
				Address:   req.Address,
				ZIP:       req.ZIP,
				Country:   req.Country,
				Phone:     req.Phone,
			},
		},
	}

	tokenResp, err := payment.CreatePaymentToken(paymentRequest, "https://checkout.bepaid.by/ctp/api/checkouts", "12498", "8d34e129dcba1cb5570e42ec0ebde0131c10169db2bec39b6e085b000e32ed3a")

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
