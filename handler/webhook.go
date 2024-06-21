package handler

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"rednit/notification"
	"rednit/payment"
	"rednit/terrors"
	"strconv"
)

func (h Handler) BepaidNotification(c echo.Context) error {
	req := new(payment.BepaidNotification)

	if err := c.Bind(req); err != nil {
		return err
	}

	username, password, ok := c.Request().BasicAuth()
	if !ok {
		return terrors.Unauthorized(errors.New("bepaid: missing credentials"), "missing credentials")
	}

	if username != "12498" || password != "8d34e129dcba1cb5570e42ec0ebde0131c10169db2bec39b6e085b000e32ed3a" {
		return terrors.Unauthorized(errors.New("bepaid: invalid credentials"), "invalid credentials")
	}

	// get order by tracking_id
	id, err := strconv.ParseInt(req.Transaction.TrackingId, 10, 64)

	if err != nil {
		return terrors.BadRequest(err, "invalid tracking_id")
	}

	order, err := h.st.GetOrderByID(id)

	if err != nil {
		return err
	}

	switch req.Transaction.Status {
	case "successful":
		order.PaymentStatus = "paid"
	case "failed":
		order.PaymentStatus = "failed"
	case "incomplete":
		order.PaymentStatus = "pending"
	case "expired":
		order.PaymentStatus = "pending"
	default:
		return terrors.BadRequest(errors.New(fmt.Sprintf("bepaid: invalid status %s", req.Transaction.Status)), "invalid status")
	}

	order.PaymentID = &req.Transaction.ID

	_, err = h.st.UpdateOrder(order)

	if err != nil {
		return err
	}

	if order.PaymentStatus == "paid" {
		go func() {
			msg := fmt.Sprintf(`Order #%d
Order paid
%s
%s
Shipping address:
%s
Payment Amount: %d %s
Payment Amount With Discount: %d %s

Purchaser information:
name: %s
email: %s
phone: %s
country: %s
postcode: %s
address: %s`,
				order.ID, "bepaid", "courier", order.Customer.Address, order.Total, order.CurrencyCode, order.Total, order.CurrencyCode, order.Customer.Name, order.Customer.Email, order.Customer.Phone, order.Customer.Country, order.Customer.ZIP, order.Customer.Address)

			msg = notification.EscapeMarkdown(msg)

			if err := notification.NotifyTelegram(h.config.Notifications.Telegram.BotToken, h.config.Notifications.Telegram.ChatID, msg); err != nil {
				log.Printf("failed to send notification to telegram: %v", err)
			}
		}()
	}

	return c.NoContent(http.StatusOK)
}
