package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"rednit/db"
	"rednit/terrors"
	"strconv"
)

func (h Handler) GetOrder(c echo.Context) error {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return terrors.BadRequest(err, "invalid order id")
	}

	order, err := h.st.GetOrder(db.GetOrderQuery{ID: &orderID})

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return terrors.NotFound(err, "order not found")
		}
		return terrors.InternalServerError(err, "failed to get order")
	}

	return c.JSON(http.StatusOK, order)
}
