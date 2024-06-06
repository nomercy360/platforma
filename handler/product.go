package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"rednit/db"
	"rednit/terrors"
)

func (h Handler) ListProducts(c echo.Context) error {
	locale := c.QueryParam("locale")
	if locale == "" {
		locale = "en"
	}

	products, err := h.st.ListProducts(locale)
	if err != nil {
		return terrors.InternalServerError(err, "failed to list products")
	}

	return c.JSON(http.StatusOK, products)
}

func (h Handler) GetProduct(c echo.Context) error {
	handle := c.Param("handle")
	locale := c.QueryParam("locale")

	if locale == "" {
		locale = "en"
	}

	product, err := h.st.GetProduct(db.GetProductQuery{Handle: handle, Locale: locale})
	if err != nil {
		return terrors.InternalServerError(err, "failed to get product")
	}

	return c.JSON(http.StatusOK, product)
}
