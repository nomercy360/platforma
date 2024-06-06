package handler

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
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

	_, err = h.st.UpdateOrder(order)

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
