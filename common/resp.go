package common

var (
	//common code (0 ~ 99)
	RespCodeSuccess  = RespCode{1000, "success"}
	RespCodeBindReParamError = RespCode{Code: 1, Message: "bind request parameters error"}
	RespCodeUnauthorizedError = RespCode{Code: 2, Message: "user unauthorized"}


	//file service code (100 ~ 199)
	RespCodeReadFileError =  RespCode{100, "read file error"}
	RespCodeOpenFileError =  RespCode{101, "open file error"}
	RespCodeCreateFileError =  RespCode{102, "create file error"}
	RespCodeCopyFileError =  RespCode{103, "copy file error"}
	RespCodeFilenameError = RespCode{104, "filename error"}
	RespCodeRemoveFileError = RespCode{105, "remove file error"}
	RespCodeNotFoundFileError = RespCode{106, "not found file error"}
	RespCodeUploadFileError = RespCode{107, "upload file error"}
	RespCodeQueryFileError = RespCode{108, "query file error"}
	RespCodeFastUploadFailed = RespCode{109, "fast upload failed"}
	RespCodeConnectFSRedisServerError = RespCode{110, "connect to redis server error"}
	RespCodeCompleteUploadError = RespCode{111, "complete upload error error"}
	RespCodeReadDataError = RespCode{112, "read data error"}
	RespCodeWriteFileError = RespCode{113, "write file error"}
	RespCodePutDataToCephError = RespCode{114, "put data to ceph error"}

	//user service code (200 ~ 299)
	RespCodeUserRegisterError = RespCode{200, "user register error"}
	RespCodeUserNotFound = RespCode{201, "user not found error"}
	RespCodeUserLoginError = RespCode{202, "user login error"}
	RespCodeUserAlreadyRegistered = RespCode{203, "user already registered"}

	//auth service code(300 ~ 399)
	RespCodeGenTokenError = RespCode{301, "gen token error"}
	RespCodeValidateTokenError = RespCode{302, "validate token error"}
	RespCodeDeleteTokenError = RespCode{303, "delete token error"}
)

type RespCode struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

type ServiceResp struct {
	code RespCode `json:"code"`
	Data interface{} `json:"data"`
}

func NewServiceResp(respCode RespCode, data interface{}) *ServiceResp {
	return &ServiceResp{
		code: respCode,
		Data: data,
	}
}


