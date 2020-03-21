package vo

type FastUploadReq struct {
	Username string `form:"username" bind:"required"`
	FileSize int64 `form:"file_size" bind:"required"`
	FileHash string `form:"file_hash" bind:"required"`
	Filename string `form:"filename" bind:"required"`
	FilePath string `form:"file_path" bind:"required"`
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
	UploadId string `form:"uploadid" bind:"required"`
	Index string  `form:"index" bind:"required"`
}

type MPUploadCompleteReq struct {
	UploadId string `form:"uploadid" bind:"required"`
	Username string `form:"username" bind:"required"`
	FileSize int64 `form:"file_size" bind:"required"`
	FileHash string `form:"file_hash" bind:"required"`
	Filename string `form:"filename" bind:"required"`
	FilePath string `form:"file_path" bind:"required"`
}

