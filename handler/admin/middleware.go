package admin

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"rednit/terrors"
	"strings"
)

func getUserID(c echo.Context) int64 {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	return claims.UID
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// skip /sign-in and /sign-up
		if strings.HasSuffix(c.Path(), "/sign-in") || strings.HasSuffix(c.Path(), "/sign-up") {
			return next(c)
		}

		// Get the token from the cookie
		cookie, err := c.Cookie("clan_cookie")
		if err != nil {
			return terrors.Unauthorized(err, "Unauthorized")
		}

		// Validate the JWT
		token, err := jwt.ParseWithClaims(cookie.Value, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("your_secret_key"), nil
		})

		if err != nil || !token.Valid {
			return terrors.Unauthorized(err, "Unauthorized")
		}

		c.Set("user", token)

		// Token is valid, proceed to the next handler
		return next(c)
	}
}
