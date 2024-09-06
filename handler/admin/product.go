package admin

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"rednit/db"
	"rednit/terrors"
)

func (a Admin) ListProducts(c echo.Context) error {
	products, err := a.s.ListProducts(db.ListProductsQuery{})

	if err != nil {
		return terrors.InternalServerError(err, "failed to list customers")
	}

	return c.JSON(http.StatusOK, products)
}
