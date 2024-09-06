package store

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"rednit/db"
	"rednit/terrors"
)

func (h Handler) ListProducts(c echo.Context) error {
	products, err := h.st.ListProducts(db.ListProductsQuery{Locale: langFromContext(c), IsPublished: true})
	if err != nil {
		return terrors.InternalServerError(err, "failed to list products")
	}

	return c.JSON(http.StatusOK, products)
}

func (h Handler) GetProduct(c echo.Context) error {
	handle := c.Param("handle")

	product, err := h.st.GetProduct(db.GetProductQuery{Handle: handle, Locale: langFromContext(c)})
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "product not found")
	} else if err != nil {
		return terrors.InternalServerError(err, "failed to get product")
	}

	return c.JSON(http.StatusOK, product)
}
