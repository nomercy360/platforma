package main

import (
	"context"
	"errors"
	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"rednit/config"
	"rednit/db"
	"rednit/handler/admin"
	"rednit/handler/store"
	"rednit/payment"
	"rednit/terrors"
	"strings"
	"time"
)

func getLoggerMiddleware(logger *slog.Logger) middleware.RequestLoggerConfig {
	return middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}
}

func localeMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		lang := "en"
		langHeader := c.Request().Header.Get("Accept-Language")
		if langHeader != "" {
			languages := strings.Split(langHeader, ",")
			if len(languages) > 0 {
				primaryLang := strings.SplitN(languages[0], ";", 2)[0]
				if primaryLang != "" {
					lang = primaryLang
				}
			}
		}

		if len(lang) > 2 && lang[2] == '-' {
			lang = lang[:2]
		}

		c.Set("lang", lang)
		return next(c)
	}
}

func getServerErrorHandler(e *echo.Echo) func(err error, context2 echo.Context) {
	return func(err error, c echo.Context) {
		var (
			code = http.StatusInternalServerError
			msg  interface{}
		)

		var he *echo.HTTPError
		var terror *terrors.Error
		switch {
		case errors.As(err, &he):
			code = he.Code
			msg = he.Message
		case errors.As(err, &terror):
			code = terror.Code
			msg = terror.Message
		default:
			msg = err.Error()
		}

		if _, ok := msg.(string); ok {
			msg = map[string]interface{}{"error": msg}
		}

		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				err = c.NoContent(code)
			} else {
				err = c.JSON(code, msg)
			}

			if err != nil {
				e.Logger.Error(err)
			}
		}
	}
}

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	cfg := config.Default{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse config: %v\n", err)
	}

	sql, err := db.ConnectDB(cfg.DBPath)

	if err != nil {
		e.Logger.Fatalf("failed to connect to db: %v", err)
	}

	if err := sql.Migrate(); err != nil {
		e.Logger.Fatalf("failed to migrate db: %v", err)
	}

	paypal, err := payment.NewPaypalClient(cfg.PayPal.ClientID, cfg.PayPal.ClientSecret, cfg.PayPal.LiveMode)
	if err != nil {
		e.Logger.Fatalf("failed to create paypal client: %v", err)
	}

	h := store.New(sql, cfg, paypal)
	a := admin.New(sql, cfg)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "https://clan-api.pages.dev", "https://plumplum.co"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	e.Use(localeMiddleware)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Use(middleware.RequestLoggerWithConfig(getLoggerMiddleware(logger)))

	e.HTTPErrorHandler = getServerErrorHandler(e)

	e.Validator = &customValidator{validator: validator.New()}

	// Routes
	api := e.Group("/api")

	adm := api.Group("/admin")

	adm.Use(admin.AuthMiddleware)

	adm.POST("/sign-in", a.LoginUser)
	adm.POST("/users", a.CreateUser)
	adm.GET("/me", a.GetUserMe)
	adm.GET("/customers", a.ListCustomers)
	adm.GET("/orders", a.ListOrders)
	adm.GET("/discounts", a.ListDiscounts)
	adm.GET("/users", a.ListUsers)

	adm.GET("/products", a.ListProducts)

	st := api.Group("/store")
	st.GET("/products", h.ListProducts)
	st.GET("/products/:handle", h.GetProduct)
	st.POST("/cart", h.CreateCart)
	st.GET("/cart/:id", h.GetCart)
	st.GET("/orders/:id", h.GetOrder)
	st.POST("/checkout", h.Checkout)
	st.POST("/cart/:id/discounts", h.ApplyDiscount)
	st.POST("/cart/:id/items", h.AddItemToCart)
	st.PUT("/cart/:id/items/:item_id", h.UpdateCartItem)
	st.DELETE("/cart/:id/items/:item_id", h.RemoveCartItem)
	st.DELETE("/cart/:id/discounts", h.DropDiscount)
	st.POST("/cart/:id/customer", h.SaveCartCustomer)
	st.POST("/cart/:id/currency", h.UpdateCartCurrency)
	st.POST("/paypal/capture", h.CapturePaypalPayment)
	st.GET("/debug", h.Debug)

	//g.PUT("/cart/:id/products", h.AddProductToCart)
	//g.DELETE("/cart/:id/products/:product_id", h.RemoveProductFromCart)

	wh := e.Group("/webhook")
	wh.POST("/bepaid", h.BepaidNotification)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server: %v", err)
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
