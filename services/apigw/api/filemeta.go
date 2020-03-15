package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-disk/common"
	"go-disk/common/rpcinterface/fileinterface"
	"go-disk/services/apigw/rpc"
	"go-disk/services/apigw/vo"
	"log"
	"net/http"
)


func GetFileMeta() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req vo.GetFileMetaReq

		if err := ctx.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := rpc.GetFileCli().GetFileMeta(context.TODO(), &fileinterface.GetFileMetaReq{
			FileHash:             req.FileHash,
		})
		if err != nil || resp.Code != int64(common.RespCodeSuccess.Code){
			log.Printf("rpc call (get metat ) error : %v", err)
			ctx.JSON(http.StatusBadRequest, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)
	}
}

func UpdateFileMeta() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req vo.UpdateFileMetaReq
		if err := ctx.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := rpc.GetFileCli().UpdateFileMeta(context.TODO(), &fileinterface.UpdateFileMetaReq{})

		if err != nil || resp.Code != int64(common.RespCodeSuccess.Code){
			log.Printf("rpc call ( update metat ) error : %v", err)
			ctx.JSON(http.StatusBadRequest, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)
	}
}

func GetFileList() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req vo.UserFileReq
		if err := ctx.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := rpc.GetFileCli().GetFileList(context.TODO(), &fileinterface.GetFileListReq{
			Username: req.Username,
			Limit: int64(req.Limit),
		})

		if err != nil || resp.Code != int64(common.RespCodeSuccess.Code){
			log.Printf("rpc call  get metat list) error : %v", err)
			ctx.JSON(http.StatusBadRequest, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)
	}
}

func DeleteFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req vo.DeleteFileReq
		if err := ctx.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := rpc.GetFileCli().DeleteFile(context.TODO(), &fileinterface.DeleteFileReq{
			FileHash: req.FileHash,
		})

		if err != nil || resp.Code != int64(common.RespCodeSuccess.Code){
			log.Printf("rpc call  get metat list) error : %v", err)
			ctx.JSON(http.StatusBadRequest, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)

	}
}

