package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"go-disk/common"
	"go-disk/common/rpcinterface/fileinterface"
	"go-disk/config"
	"go-disk/model"
	"log"
	"net/http"
)

var fileMetaCli fileinterface.FileService

func init() {
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			config.ConsulAddress,
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.file"),
	)

	service.Init()

	fileMetaCli = fileinterface.NewFileService("go.micro.service.file", service.Client())
}

func GetFileMeta() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.GetFileMetaReq

		if err := ctx.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := fileMetaCli.GetFileMeta(context.TODO(), &fileinterface.GetFileMetaReq{
			FileHash:             req.FileHash,
		})
		if err != nil {
			log.Printf("rpc call (get metat ) error : %v", err)
			ctx.JSON(http.StatusBadRequest, *resp)
			return
		}

		ctx.JSON(http.StatusOK, common.NewServiceResp(common.RespCodeSuccess, *resp))
	}
}

func UpdateFileMeta() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.UpdateFileMetaReq
		if err := ctx.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := fileMetaCli.UpdateFileMeta(context.TODO(), &fileinterface.UpdateFileMetaReq{})

		if err != nil {
			log.Printf("rpc call ( update metat ) error : %v", err)
			ctx.JSON(http.StatusBadRequest, *resp)
			return
		}

		ctx.JSON(http.StatusOK, common.NewServiceResp(common.RespCodeSuccess, *resp))
	}
}

func GetFileList() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.UserFileReq
		if err := ctx.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := fileMetaCli.GetFileList(context.TODO(), &fileinterface.GetFileListReq{
			Username: req.Username,
			Limit: int64(req.Limit),
		})

		if err != nil {
			log.Printf("rpc call  get metat list) error : %v", err)
			ctx.JSON(http.StatusBadRequest, *resp)
			return
		}

		ctx.JSON(http.StatusOK, common.NewServiceResp(common.RespCodeSuccess, *resp))
	}
}

func DeleteFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.DeleteFileReq
		if err := ctx.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := fileMetaCli.DeleteFile(context.TODO(), &fileinterface.DeleteFileReq{
			FileHash: req.FileHash,
		})

		if err != nil {
			log.Printf("rpc call  get metat list) error : %v", err)
			ctx.JSON(http.StatusBadRequest, *resp)
			return
		}

		ctx.JSON(http.StatusOK, common.NewServiceResp(common.RespCodeSuccess, *resp))

	}
}

