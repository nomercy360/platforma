package admin

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"rednit/db"
	"rednit/terrors"
	"time"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UID    int64 `json:"uid"`
	ChatID int64 `json:"chat_id"`
}

func generateJWT(secret string, uid int64) (string, error) {
	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		UID: uid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (a Admin) LoginUser(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return terrors.BadRequest(err, "Invalid request")
	}

	if err := c.Validate(req); err != nil {
		return terrors.BadRequest(err, "Invalid request")
	}

	user, err := a.s.GetUser(db.UserQuery{Email: req.Email})
	if err != nil {
		return terrors.Unauthorized(err, "Unauthorized")
	}

	if err = CheckPassword(user.Password, req.Password); err != nil {
		return terrors.Unauthorized(err, "Unauthorized")
	}

	token, err := generateJWT("your_secret_key", user.ID)
	if err != nil {
		return terrors.InternalServerError(err, "Failed to generate JWT")
	}

	// Set JWT in HttpOnly cookie
	cookie := new(http.Cookie)
	cookie.Name = "clan_cookie"
	cookie.Value = token
	cookie.Secure = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Path = "/"
	cookie.MaxAge = 86400
	cookie.HttpOnly = true

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, user)
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (a Admin) CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return terrors.BadRequest(err, "Invalid request")
	}

	if err := c.Validate(req); err != nil {
		return terrors.BadRequest(err, "Invalid request")
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return terrors.InternalServerError(err, "Failed to hash password")
	}

	user, err := a.s.CreateUser(req.Email, hashedPassword)
	if err != nil {
		return terrors.InternalServerError(err, "Failed to create user")
	}

	return c.JSON(http.StatusCreated, user)
}

func (a Admin) GetUserMe(c echo.Context) error {
	uid := getUserID(c)

	user, err := a.s.GetUser(db.UserQuery{ID: uid})
	if err != nil {
		return terrors.InternalServerError(err, "Failed to get user")
	}

	return c.JSON(http.StatusOK, user)
}
