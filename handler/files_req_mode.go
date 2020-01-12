package handler

type GetFileMetaReq struct {
	FileHash string `form:"file_hash" bind:"required"`
}

