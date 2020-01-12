package common

var (
	//common code (0 ~ 99)
	RespCodeSuccess  = RespCode{0, "success"}
	RespCodeError = RespCode{Code: 1, Message:"error"}

	//file service code (100 ~ 199)
	RespCodeReadFileError =  RespCode{100, "read file error"}
	RespCodeOpenFileError =  RespCode{101, "open file error"}
	RespCodeCreateFileError =  RespCode{101, "create file error"}
	RespCodeCopyFileError =  RespCode{102, "copy file error"}
)

type RespCode struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

type ServiceResp struct {
	RespCode RespCode `json:"resp_code"`
	Data interface{} `json:"data"`
}

func NewServiceResp(respCode RespCode, data interface{}) *ServiceResp {
	return &ServiceResp{
		RespCode: respCode,
		Data: data,
	}
}


