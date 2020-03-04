package rpc

import (
	"context"
	"go-disk/common"
	"go-disk/services/download/config"
	"go-disk/services/download/proto"
)

type Download struct {

}

func (d *Download) DownloadEndPoint(ctx context.Context, req *proto.DownloadEndpointReq, resp *proto.DownloadEndpointResp) error {
	resp.Code = int64(common.RespCodeSuccess.Code)
	resp.Message = common.RespCodeSuccess.Message
	resp.Data = &proto.DownloadEndpointResp_Data{Endpoint: config.DownloadServiceEndpoint}
	return nil
}

