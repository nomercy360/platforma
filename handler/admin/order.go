package admin

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"rednit/terrors"
)

func (a Admin) ListOrders(c echo.Context) error {
	orders, err := a.s.ListOrders()
	if err != nil {
		return terrors.InternalServerError(err, "failed to list orders")
	}

	return c.JSON(http.StatusOK, orders)
}
