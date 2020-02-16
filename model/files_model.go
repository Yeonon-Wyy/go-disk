package model

import "time"

type GetFileMetaReq struct {
	FileHash string `form:"file_hash" bind:"required"`
}

type DownloadFileReq struct {
	FileHash string `form:"file_hash" bind:"required"`
}

type UpdateFileMetaReq struct {
	Username string `form:"username" bind:"required"`
	FileHash string `form:"file_hash" bind:"required"`
	Filename string `form:"filename" bind:"required"`
	Status string `form:"status"`
}

type DeleteFileReq struct {
	FileHash string `form:"file_hash" bind:"required"`
}

type UserFileResp struct {
	Username string `json:"username"`
	Filename string `json:"filename"`
	FileHash string `json:"file_hash"`
	FileSize int64 `json:"file_size"`
	UploadAt time.Time `json:"upload_at"`
	LastUpdate time.Time `json:"last_update"`
}

type UserFileReq struct {
	Username string `form:"username" bind:"required"`
	Limit int `form:"limit" bind:"required"`
}

type FastUploadReq struct {
	Username string `form:"username" bind:"required"`
	FileSize int64 `form:"file_size" bind:"required"`
	FileHash string `form:"file_hash" bind:"required"`
	Filename string `form:"filename" bind:"required"`
}

//for multipart upload
type MPUploadInfo struct {
	FileHash string `json:"file_hash"`
	FileSize int64 `json:"file_size"`
	UploadId string `json:"upload_id"`
	ChunkSize int `json:"chunk_size"`
	ChunkCount int `json:"chunk_count"`
}

type MPUploadInitReq struct {
	Username string `form:"username" bind:"required"`
	FileSize int64 `form:"file_size" bind:"required"`
	FileHash string `form:"file_hash" bind:"required"`
}

type MPUploadPartReq struct {
	UploadId string `from:"uploadid" bind:"required"`
	Index string  `from:"index" bind:"required"`
}

type MPUploadCompleteReq struct {
	UploadId string `from:"uploadid" bind:"required"`
	Username string `form:"username" bind:"required"`
	FileSize int64 `form:"filesize" bind:"required"`
	FileHash string `form:"filehash" bind:"required"`
	Filename string `form:"filename" bind:"required"`
}