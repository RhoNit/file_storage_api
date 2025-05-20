package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/RhoNit/file_storage_api/common"
	"github.com/RhoNit/file_storage_api/internal/models"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Register a new user
// @Description Register a new user with username and password
// @Tags user_registration
// @Accept json
// @Produce json
// @Param register_request body models.RegisterRequest true "User registration info"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /register [post]
func (h *Handler) RegisterUserHandler(c echo.Context) error {
	var req models.RegisterRequest

	// de-serialize or unmarshall the request body into go-native type
	if err := c.Bind(&req); err != nil {
		h.ZapLogger.Error(
			"Failed while de-serializing the request payload into Go-struct type",
			zap.Error(err),
		)

		return c.JSON(
			http.StatusBadRequest,
			echo.Map{"error": "Invalid request"},
		)
	}

	// check whether the req.Username is already present in the database or not
	if _, exists := common.UsersStore[req.Username]; exists {
		h.ZapLogger.Warn(
			"Username already exists",
			zap.String("Username", req.Username),
		)

		return c.JSON(
			http.StatusBadRequest,
			echo.Map{"error": "Username already exists"},
		)
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.ZapLogger.Error(
			"Error while hashing the password",
			zap.String("associated username", req.Username),
			zap.Error(err),
		)

		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": "Failed to process password"},
		)
	}

	// load location for setting time as per IST
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		h.ZapLogger.Error(
			"Error while loadinf location",
			zap.Error(err),
		)
	}

	nowTime := time.Now().In(loc)

	user := &models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		StorageUsed:  0,
		CreatedAt:    nowTime,
	}

	// store the user into the in-memory store
	common.UsersStore[req.Username] = user
	common.FileMetadataStore[req.Username] = []models.FileMetadata{}

	// Create user directory
	userDir := filepath.Join("file_store", req.Username)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		h.ZapLogger.Error(
			"Error while creating user storage directory",
			zap.String("associated user", user.Username),
			zap.Error(err),
		)

		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": "Failed to create user directory"},
		)
	}

	h.ZapLogger.Info(
		"User has been registered successfully",
		zap.String("associated user", user.Username),
	)

	return c.JSON(
		http.StatusCreated,
		echo.Map{"message": "User registered successfully"},
	)
}
