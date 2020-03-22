package rpc

import (
	"context"
	"go-disk/common"
	"go-disk/common/rpcinterface/downloadinterface"
	"go-disk/services/download/config"
)

type Download struct {
}

func (d *Download) DownloadEndPoint(ctx context.Context, req *downloadinterface.DownloadEndpointReq, resp *downloadinterface.DownloadEndpointResp) error {
	resp.Code = int64(common.RespCodeSuccess.Code)
	resp.Message = common.RespCodeSuccess.Message
	resp.Data = &downloadinterface.DownloadEndpointResp_Data{Endpoint: config.Conf.Business.DownloadServiceEndpoint}
	return nil
}
