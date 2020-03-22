package rpc

import (
	"context"
	"go-disk/common"
	"go-disk/common/rpcinterface/uploadinterface"
	"go-disk/services/upload/config"
)

var (
	businessConfig = config.Conf.Business
)

type EndPoint struct {
}

func (e *EndPoint) UploadEndPoint(ctx context.Context, req *uploadinterface.UploadEndPointReq, resp *uploadinterface.UploadEndPointResp) error {
	resp.Code = int64(common.RespCodeSuccess.Code)
	resp.Message = common.RespCodeSuccess.Message
	resp.Data = &uploadinterface.UploadEndPointResp_Data{Endpoint: businessConfig.UploadServiceEndpoint}
	return nil
}
