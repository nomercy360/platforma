package handler

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"rednit/db"
	"rednit/notification"
	"rednit/payment"
	"rednit/terrors"
	"strconv"
)

func (h Handler) telegramOrderPaid(order db.Order) {
	var delivery string
	if order.CurrencyCode == "BYN" {
		delivery = "Сдэком по СНГ"
	} else {
		delivery = "Международная Экспресс доставка"
	}
	msg := fmt.Sprintf(`%s
Оплата: %s
Тип доставки: %s
Адрес доставки: %s
Сумма заказа: %d %s
Итого (включая доставку и дискаунт): %d %s

Покупатель:
Имя: %s
Email: %s
Телефон: %s
Страна: %s
Индекс: %s
Адрес: %s`,
		order.ToString(), order.PaymentProvider, delivery, *order.Customer.Address, order.Subtotal, order.CurrencyCode, order.Total, order.CurrencyCode, *order.Customer.Name, order.Customer.Email, *order.Customer.Phone, *order.Customer.Country, *order.Customer.ZIP, *order.Customer.Address)

	msg = notification.EscapeMarkdown(msg)

	if err := notification.NotifyTelegram(h.config.Notifications.Telegram.BotToken, h.config.Notifications.Telegram.ChatID, msg); err != nil {
		log.Printf("failed to send notification to telegram: %v", err)
	}
}

func (h Handler) BepaidNotification(c echo.Context) error {
	req := new(payment.BepaidNotification)

	if err := c.Bind(req); err != nil {
		return err
	}

	username, password, ok := c.Request().BasicAuth()
	if !ok {
		return terrors.Unauthorized(errors.New("bepaid: missing credentials"), "missing credentials")
	}

	if username != h.config.Bepaid.ShopID || password != h.config.Bepaid.SecretKey {
		return terrors.Unauthorized(errors.New("bepaid: invalid credentials"), "invalid credentials")
	}

	// get order by tracking_id
	id, err := strconv.ParseInt(req.Transaction.TrackingId, 10, 64)

	if err != nil {
		return terrors.BadRequest(err, "invalid tracking_id")
	}

	order, err := h.st.GetOrder(db.GetOrderQuery{ID: &id})

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
			h.telegramOrderPaid(*order)
		}()
	}

	return c.NoContent(http.StatusOK)
}
