package common

import "github.com/RhoNit/file_storage_api/internal/models"

// In-memory storage
var (
	UsersStore        = make(map[string]*models.User)
	FileMetadataStore = make(map[string][]models.FileMetadata)
)
