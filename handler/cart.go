package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"rednit/db"
	"rednit/terrors"
	"strconv"
)

type CartItemInput struct {
	VariantID int64 `json:"variant_id" validate:"required"`
	Quantity  int   `json:"quantity" validate:"required,min=1"`
}

func (h Handler) CreateCart(c echo.Context) error {
	var req CartItemInput
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	cart := db.Cart{
		Items: []db.LineItem{
			{
				VariantID: req.VariantID,
				Quantity:  req.Quantity,
			},
		},
	}

	ip := c.RealIP()
	ua := c.Request().UserAgent()

	cart.Context = db.CustomerContext{
		IP:        ip,
		UserAgent: ua,
	}

	createdCart, err := h.st.CreateCart(cart, langFromContext(c))

	if err != nil {
		return terrors.InternalServerError(err, "failed to create cart")
	}

	return c.JSON(http.StatusCreated, createdCart)
}

func (h Handler) AddItemToCart(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return terrors.BadRequest(err, "invalid cart id")
	}

	var req CartItemInput
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	item := db.LineItem{
		VariantID: req.VariantID,
		Quantity:  req.Quantity,
	}

	if err := h.st.SaveLineItem(db.LineItem{
		CartID:    &id,
		VariantID: item.VariantID,
		Quantity:  item.Quantity,
	}); err != nil {
		return terrors.InternalServerError(err, "failed to add item to cart")
	}

	cart, err := h.st.GetCartByID(id, langFromContext(c))

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "cart not found")
	} else if err != nil {
		return terrors.InternalServerError(err, "failed to get cart")
	}

	return c.JSON(http.StatusOK, cart)
}

func (h Handler) GetCart(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return terrors.BadRequest(err, "invalid cart id")
	}

	cart, err := h.st.GetCartByID(id, langFromContext(c))

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "cart not found")
	} else if err != nil {
		return terrors.InternalServerError(err, "failed to get cart")
	}

	return c.JSON(http.StatusOK, cart)
}

type ApplyDiscountRequest struct {
	Code string `json:"code" validate:"required"`
}

func (h Handler) ApplyDiscount(c echo.Context) error {
	cartID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return terrors.BadRequest(err, "invalid cart id")
	}

	var req ApplyDiscountRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	discount, err := h.st.GetDiscount(db.DiscountQuery{
		Code: req.Code,
	})

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "discount not found")
	} else if err != nil {
		return terrors.InternalServerError(err, "failed to get discount")
	}

	if !discount.IsValid() {
		return terrors.BadRequest(errors.New("invalid discount"), "discount is not valid")
	}

	if err := h.st.UpdateDiscountUsageCount(discount.ID); err != nil {
		return terrors.InternalServerError(err, "failed to update discount usage count")
	}

	if err := h.st.UpdateCartDiscount(cartID, discount.ID); err != nil {
		return terrors.InternalServerError(err, "failed to update cart discount")
	}

	cart, err := h.st.GetCartByID(cartID, langFromContext(c))

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "cart not found")
	} else if err != nil {
		return terrors.InternalServerError(err, "failed to get cart")
	}

	return c.JSON(http.StatusOK, cart)
}

func (h Handler) DropDiscount(c echo.Context) error {
	cartID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return terrors.BadRequest(err, "invalid cart id")
	}

	if err := h.st.DropCartDiscount(cartID); err != nil {
		return terrors.InternalServerError(err, "failed to drop cart discount")
	}

	cart, err := h.st.GetCartByID(cartID, langFromContext(c))

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "cart not found")
	} else if err != nil {
		return terrors.InternalServerError(err, "failed to get cart")
	}

	return c.JSON(http.StatusOK, cart)
}
