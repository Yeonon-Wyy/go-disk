package constant

//status
const (
	FileStatusAvailable = 1
	FileStatusDisable = 0
	FileStatusDelete = -1
)

//for multipart upload
const (
	FileMPUploadChunkSize = 5 * 1024 * 1024 //5MB

)
