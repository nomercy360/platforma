package admin

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"rednit/terrors"
)

func (a Admin) ListCustomers(c echo.Context) error {
	customers, err := a.s.ListCustomers()
	if err != nil {
		return terrors.InternalServerError(err, "failed to list customers")
	}

	return c.JSON(http.StatusOK, customers)
}
