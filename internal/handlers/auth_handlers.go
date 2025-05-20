package handlers

import (
	"net/http"
	"time"

	"github.com/RhoNit/file_storage_api/common"
	"github.com/RhoNit/file_storage_api/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Login user
// @Description Login with username and password
// @Tags user_login
// @Accept json
// @Produce json
// @Param login_request body models.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /login [post]
func (h *Handler) LoginUserHandler(c echo.Context) error {
	var req models.LoginRequest

	// unmarshalling the request payload into go-struct type
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			echo.Map{"error": "Invalid request"},
		)
	}

	// check if the username is present in the data store or not
	user, exists := common.UsersStore[req.Username]
	if !exists {
		h.ZapLogger.Warn(
			"Username doesn't exist in the data store",
			zap.String("associated username", req.Username),
		)

		return c.JSON(
			http.StatusUnauthorized,
			echo.Map{"error": "Invalid credentials.. username doesn't exist in the data store"},
		)
	}

	// compare the request password and the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		h.ZapLogger.Error(
			"Error while comparing the request password and hashed password",
			zap.String("associated user", req.Username),
			zap.Error(err),
		)

		return c.JSON(
			http.StatusUnauthorized,
			echo.Map{"error": "Invalid credentials"},
		)
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		h.ZapLogger.Error(
			"Error while signing the secret key",
			zap.String("associated user", req.Username),
			zap.Error(err),
		)

		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": "Failed to sign the string token"},
		)
	}

	h.ZapLogger.Info(
		"Successfully generated token and logged in",
		zap.String("associated user", req.Username),
	)

	return c.JSON(
		http.StatusOK,
		echo.Map{"jwt_token": tokenString},
	)
}
