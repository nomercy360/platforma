package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"rednit/db"
	"rednit/terrors"
	"strconv"
)

type CreateCartRequest struct {
	CustomerID int           `json:"customer_id"`
	Items      []db.CartItem `json:"items"  validate:"required,dive,required"`
}

func (h Handler) CreateCart(c echo.Context) error {
	var req CreateCartRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	cart := db.Cart{
		CustomerID: req.CustomerID,
		Items:      req.Items,
	}

	ip := c.RealIP()
	ua := c.Request().UserAgent()

	cart.Context = db.CustomerContext{
		IP:        ip,
		UserAgent: ua,
	}

	createdCart, err := h.st.CreateCart(cart)

	if err != nil {
		return terrors.InternalServerError(err, "failed to create cart")
	}

	return c.JSON(http.StatusCreated, createdCart)
}

func (h Handler) GetCart(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return terrors.BadRequest(err, "invalid cart id")
	}

	cart, err := h.st.GetCartByID(id)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "cart not found")
	} else if err != nil {
		return terrors.InternalServerError(err, "failed to get cart")
	}

	return c.JSON(http.StatusOK, cart)
}
