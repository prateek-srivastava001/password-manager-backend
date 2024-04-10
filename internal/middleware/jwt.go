package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
		}

		tokenString := strings.Split(authHeader, " ")[1]
		if tokenString == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}
			return []byte("mostestsecretkeyever"), nil
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		if !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}
		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().UTC().After(expirationTime) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token has expired")
		}

		c.Set("user", claims["email"])

		return next(c)
	}
}
