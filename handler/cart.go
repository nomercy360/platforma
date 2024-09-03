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
	VariantID int64  `json:"variant_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
	Currency  string `json:"currency_code" validate:"required,iso4217"`
}

func getCountryFromContext(c echo.Context) *string {
	if countryHeader := c.Request().Header.Get("Cf-Ipcountry"); countryHeader != "" {
		return &countryHeader
	}

	return nil
}

func getIPFromContext(c echo.Context) string {
	if ip := c.Request().Header.Get("Cf-Connecting-Ip"); ip != "" {
		return ip
	}

	return c.RealIP()
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
		CurrencyCode: req.Currency,
		Items: []db.LineItem{
			{
				VariantID: req.VariantID,
				Quantity:  req.Quantity,
			},
		},
	}

	ua := c.Request().UserAgent()

	cart.Context = db.CustomerContext{
		IP:        getIPFromContext(c),
		UserAgent: ua,
		Country:   getCountryFromContext(c),
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

	var req CreateCartRequest
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
	}); err != nil && errors.Is(err, db.ErrAlreadyExists) {
		return terrors.Conflict(err, "failed to save line item")
	} else if err != nil {
		return terrors.InternalServerError(err, "failed to save line item")
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

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}

func (h Handler) UpdateCartItem(c echo.Context) error {
	cartID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	itemID, _ := strconv.ParseInt(c.Param("item_id"), 10, 64)

	if itemID == 0 || cartID == 0 {
		return terrors.BadRequest(errors.New("invalid cart or item id"), "invalid cart or item id")
	}

	var req UpdateCartItemRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	if err := h.st.UpdateLineItemQuantity(itemID, req.Quantity); err != nil {
		return terrors.InternalServerError(err, "failed to update item quantity")
	}

	cart, err := h.st.GetCartByID(cartID, langFromContext(c))

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "cart not found")
	} else if err != nil {
		return terrors.InternalServerError(err, "failed to get cart")
	}

	return c.JSON(http.StatusOK, cart)
}

func (h Handler) RemoveCartItem(c echo.Context) error {
	cartID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	itemID, _ := strconv.ParseInt(c.Param("item_id"), 10, 64)

	if itemID == 0 || cartID == 0 {
		return terrors.BadRequest(errors.New("invalid cart or item id"), "invalid cart or item id")
	}

	if err := h.st.RemoveLineItem(itemID); err != nil {
		return terrors.InternalServerError(err, "failed to remove item")
	}

	cart, err := h.st.GetCartByID(cartID, langFromContext(c))

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "cart not found")
	} else if err != nil {
		return terrors.InternalServerError(err, "failed to get cart")
	}

	return c.JSON(http.StatusOK, cart)
}

func (h Handler) SaveCartCustomer(c echo.Context) error {
	cartID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if cartID == 0 {
		return terrors.BadRequest(errors.New("invalid cart id"), "invalid cart id")
	}

	var req db.Customer
	if err := c.Bind(&req); err != nil {
		return err
	}

	if req.Email == "" {
		return terrors.BadRequest(errors.New("invalid email"), "invalid email")
	}

	customer, err := h.st.GetCustomerByEmail(req.Email)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		customer, err = h.st.AddCustomer(req)
		if err != nil {
			return terrors.InternalServerError(err, "failed to add customer")
		}
	} else if err != nil {
		return terrors.InternalServerError(err, "failed to get customer")
	}

	if err := h.st.UpdateCartCustomer(cartID, customer.ID); err != nil {
		return terrors.InternalServerError(err, "failed to update cart customer")
	}

	cart, err := h.st.GetCartByID(cartID, langFromContext(c))

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "cart not found")
	}

	if err != nil {
		return terrors.InternalServerError(err, "failed to get cart")
	}

	return c.JSON(http.StatusOK, cart)
}

func (h Handler) UpdateCartCurrency(c echo.Context) error {
	cartID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if cartID == 0 {
		return terrors.BadRequest(errors.New("invalid cart id"), "invalid cart id")
	}

	var req struct {
		Currency string `json:"currency_code" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	if err := h.st.UpdateCartCurrency(cartID, req.Currency); err != nil {
		return terrors.InternalServerError(err, "failed to update cart currency")
	}

	cart, err := h.st.GetCartByID(cartID, langFromContext(c))

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "cart not found")
	} else if err != nil {
		return terrors.InternalServerError(err, "failed to get cart")
	}

	return c.JSON(http.StatusOK, cart)
}
