package models

import "time"

type FileMetadata struct {
	Filename     string    `json:"filename"`
	OriginalName string    `json:"originalName"`
	Size         int64     `json:"size"`
	UploadTime   time.Time `json:"uploadTime"`
	Username     string    `json:"username"`
}

type StorageInfo struct {
	TotalStorage     int64 `json:"totalStorage"`
	UsedStorage      int64 `json:"usedStorage"`
	RemainingStorage int64 `json:"remainingStorage"`
}

type PaginatedResponse struct {
	Data       []FileMetadata `json:"data"`
	Page       int            `json:"page"`
	PageSize   int            `json:"pageSize"`
	TotalItems int            `json:"totalItems"`
	TotalPages int            `json:"totalPages"`
}
