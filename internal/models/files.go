package models

import "time"

type FileStruct struct {
	Filename     string    `json:"filename"`
	OriginalName string    `json:"originalName"`
	Size         int64     `json:"size"`
	UploadTime   time.Time `json:"uploadTime"`
	Username     string    `json:"username"`
}

type StorageInfoStruct struct {
	TotalStorage     int64 `json:"totalStorage"`
	UsedStorage      int64 `json:"usedStorage"`
	RemainingStorage int64 `json:"remainingStorage"`
}

type GellAllFilesPaginatedResponseStruct struct {
	Data       []FileStruct `json:"data"`
	Page       int          `json:"page"`
	PageSize   int          `json:"pageSize"`
	TotalItems int          `json:"totalItems"`
	TotalPages int          `json:"totalPages"`
}
