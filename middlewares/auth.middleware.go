package middlewares

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	auth *auth.Client
}

func NewAuth(client *auth.Client) (auth *Auth) {
	return &Auth{
		auth: client,
	}
}

func (m *Auth) Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		authToken, err := m.auth.VerifyIDToken(context.Background(), token)
		if err != nil {
			fmt.Println(err.Error())
			c.Error(echo.ErrUnauthorized)
			return nil
		}
		c.Set("authToken", authToken)
		if err = next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}
