package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/RhoNit/file_storage_api/common"
	"github.com/RhoNit/file_storage_api/internal/models"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	DefaultStorageQuota = 50 * 1024 * 1024 // 50MB (in bytes)
	DefaultPageSize     = 10               // Default no. of items per page
)

// @Summary Upload a file
// @Description Upload a file (requires authentication)
// @Tags files
// @Accept mpfd
// @Produce json
// @Param file formData file true "File to be uploaded"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /upload [post]
func (h *Handler) UploadFileHandler(c echo.Context) error {
	// username := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	username := c.Get("username").(string)
	user := common.UsersStore[username]

	file, err := c.FormFile("file")
	if err != nil {
		h.ZapLogger.Error(
			"Error while returning the multipart form file",
			zap.Any("associated user", user),
			zap.String("associated filename", file.Filename),
			zap.Error(err),
		)

		return c.JSON(
			http.StatusBadRequest,
			echo.Map{"error": "No file uploaded"},
		)
	}

	if user.StorageUsed+file.Size > DefaultStorageQuota {
		h.ZapLogger.Error(
			"Error due to user's storage quota exceeded",
			zap.Any("associated user", user),
			zap.String("associated filename", file.Filename),
			zap.Int64("associated file size", file.Size),
			zap.Error(err),
		)

		return c.JSON(
			http.StatusBadRequest,
			echo.Map{"error": "Storage quota exceeded"},
		)
	}

	src, err := file.Open()
	if err != nil {
		h.ZapLogger.Error(
			"Error while opening the file",
			zap.Any("associated user", user),
			zap.String("associated filename", file.Filename),
			zap.Error(err),
		)

		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": "Failed to open file"},
		)
	}
	defer src.Close()

	// Create user directory if it doesn't exist
	userDir := filepath.Join("file_store", username)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		h.ZapLogger.Error(
			"Error while creating user storage dir",
			zap.Any("associated user", user),
			zap.Error(err),
		)

		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": "Failed to create user directory"},
		)
	}

	// Create destination file
	dst, err := os.Create(filepath.Join(userDir, file.Filename))
	if err != nil {
		h.ZapLogger.Error(
			"Error while creating destination file",
			zap.Any("associated user", user),
			zap.String("associated filename", dst.Name()),
			zap.Error(err),
		)

		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": "Failed to create destination file"},
		)
	}
	defer dst.Close()

	// Copy file contents
	if _, err = io.Copy(dst, src); err != nil {
		h.ZapLogger.Error(
			"Error while copying file contents",
			zap.Any("associated user", user),
			zap.String("associated filename", file.Filename),
			zap.Error(err),
		)

		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": "Failed to save file"},
		)
	}

	// load location as per Asia/Kolkata time
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		h.ZapLogger.Error(
			"Error while loading location",
			zap.Error(err),
		)
	}

	timeNow := time.Now().In(loc)

	// Update user storage and metadata
	user.StorageUsed += file.Size
	metadata := models.FileStruct{
		Filename:     file.Filename,
		OriginalName: file.Filename,
		Size:         file.Size,
		UploadTime:   timeNow,
		Username:     username,
	}
	common.FileMetadataStore[username] = append(common.FileMetadataStore[username], metadata)

	h.ZapLogger.Info(
		"Successfully uploaded file",
		zap.Any("associated user", user),
		zap.String("associated filename", file.Filename),
	)

	return c.JSON(
		http.StatusOK,
		echo.Map{"message": "File uploaded successfully"},
	)
}

// @Summary Get remaining storage
// @Description Get remaining storage for authenticated user
// @Tags user_storage_info
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.StorageInfo
// @Router /storage/remaining [get]
func (h *Handler) GetRemainingStorageHandler(c echo.Context) error {
	// username := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	username := c.Get("username").(string)
	user := common.UsersStore[username]

	info := models.StorageInfoStruct{
		TotalStorage:     DefaultStorageQuota,
		UsedStorage:      user.StorageUsed,
		RemainingStorage: DefaultStorageQuota - user.StorageUsed,
	}

	h.ZapLogger.Info(
		"Remaining Storage",
		zap.Any("associated user data", user),
		zap.Any("associated storage", info),
	)

	return c.JSON(
		http.StatusOK,
		echo.Map{
			"username":     username,
			"storage info": info,
		},
	)
}

// @Summary Get user files
// @Description Get list of files uploaded by authenticated user with pagination
// @Tags files
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Number of items per page (default: 10)"
// @Success 200 {object} models.PaginatedResponse
// @Router /files [get]
func (h *Handler) GetUserFilesHandler(c echo.Context) error {
	// username := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	username := c.Get("username").(string)

	h.ZapLogger.Info(
		"all files associated with an user",
		zap.Any("user", username),
	)

	page := 1
	pageSize := DefaultPageSize

	if pageStr := c.QueryParam("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.QueryParam("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	// get user's files
	files := common.FileMetadataStore[username]
	totalItems := len(files)
	totalPages := (totalItems + pageSize - 1) / pageSize

	// calculate start and end indices
	start := (page - 1) * pageSize
	end := start + pageSize

	if end > totalItems {
		end = totalItems
	}

	// ensure start index is valid
	if start >= totalItems {
		start = 0
		end = 0
	}

	// construct paginated response
	resp := models.GellAllFilesPaginatedResponseStruct{
		Data:       files[start:end],
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	return c.JSON(
		http.StatusOK,
		echo.Map{"paginated_response": resp},
	)
}
