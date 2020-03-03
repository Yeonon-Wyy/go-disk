package rpc

import (
	"context"
	"go-disk/common"
	"go-disk/services/upload/config"
	"go-disk/services/upload/proto"
)

type EndPoint struct {

}

func (e *EndPoint) UploadEndPoint(ctx context.Context, req *proto.UploadEndPointReq, resp *proto.UploadEndPointResp) error {
	resp.Code = int64(common.RespCodeSuccess.Code)
	resp.Message = common.RespCodeSuccess.Message
	resp.Data = &proto.UploadEndPointResp_Data{Endpoint:config.UploadServiceEndpoint}
	return nil
}

