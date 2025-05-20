package common

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(
				http.StatusUnauthorized,
				echo.Map{"error": "Authorization header is required"},
			)
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(
				http.StatusUnauthorized,
				echo.Map{"error": "Authorization header format must be Bearer {token}"},
			)
		}

		// Get the token string
		tokenString := parts[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid token signing method")
			}
			return []byte("secret"), nil
		})

		if err != nil {
			return c.JSON(
				http.StatusUnauthorized,
				echo.Map{"error": "Invalid or expired token"},
			)
		}

		// Check if the token is valid
		if !token.Valid {
			return c.JSON(
				http.StatusUnauthorized,
				echo.Map{"error": "Invalid token"},
			)
		}

		// Get the claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(
				http.StatusUnauthorized,
				echo.Map{"error": "Invalid token claims"},
			)
		}

		// Check if the token is expired
		exp, err := claims.GetExpirationTime()
		if err != nil || exp.Before(time.Now()) {
			return c.JSON(
				http.StatusUnauthorized,
				echo.Map{"error": "Token has expired"},
			)
		}

		// Get the username from claims
		username, ok := claims["username"].(string)
		if !ok {
			return c.JSON(
				http.StatusUnauthorized,
				echo.Map{"error": "Invalid token claims"},
			)
		}

		// Set the user information in the context
		c.Set("user", token)
		c.Set("username", username)

		return next(c)
	}
}
