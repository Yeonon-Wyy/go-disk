package rpc

import (
	"context"
	"go-disk/common"
	"go-disk/common/rpcinterface/fileinterface"
	"go-disk/services/file/config"
	"go-disk/services/file/db"
	"go-disk/services/file/store/ceph"
	"log"
)

type FileService struct {
}

func (f FileService) GetFileMeta(ctx context.Context, req *fileinterface.GetFileMetaReq, resp *fileinterface.GetFileMetaResp) error {

	tblFile, err := db.GetFileMeta(req.FileHash)
	if err != nil {
		resp.Code = int64(common.RespCodeNotFoundFileError.Code)
		resp.Message = common.RespCodeNotFoundFileError.Message
		return nil
	}
	resp.Code = int64(common.RespCodeSuccess.Code)
	resp.Message = common.RespCodeSuccess.Message
	resp.Data = &fileinterface.GetFileMetaResp_Data{
		FileSha1: tblFile.FileSha1,
		Filename: tblFile.Filename,
		FileSize: tblFile.FileSize,
		Location: tblFile.FileLocation,
		UploadAt: tblFile.CreateAt.String(),
		UpdateAt: tblFile.UpdateAt.String(),
		Status:   int64(tblFile.Status),
	}

	return nil
}

func (f FileService) UpdateFileMeta(ctx context.Context, req *fileinterface.UpdateFileMetaReq, resp *fileinterface.UpdateFileMetaResp) error {

	tblFile, err := db.GetFileMeta(req.FileHash)
	if err != nil {
		resp.Code = int64(common.RespCodeNotFoundFileError.Code)
		resp.Message = common.RespCodeNotFoundFileError.Message
		return nil
	}

	if req.Filename != "" {
		tblFile.Filename = req.Filename
		db.OnFileUpdateFinished(
			tblFile.FileSha1,
			tblFile.Filename,
			)

		//更新到user file关联表
		db.UpdateUserFilename(req.Username, req.FileHash, req.Filename)

		resp.Code = int64(common.RespCodeSuccess.Code)
		resp.Message = common.RespCodeSuccess.Message
		return nil
	}

	log.Printf("filename cann't equals empty")
	resp.Code = int64(common.RespCodeSuccess.Code)
	resp.Message = common.RespCodeSuccess.Message
	return nil
}

func (f FileService) GetFileList(ctx context.Context, req *fileinterface.GetFileListReq, resp *fileinterface.GetFileListResp) error {

	userFiles, ok := db.QueryUserFileMetas(req.Username, int(req.Limit))
	if !ok {
		resp.Code = int64(common.RespCodeQueryFileError.Code)
		resp.Message = common.RespCodeQueryFileError.Message
		return nil
	}

	//transfer data to rpc data type
	dataList := []*fileinterface.GetFileListResp_Data{}

	for _, userFile := range userFiles {
		data := &fileinterface.GetFileListResp_Data{
			Username:             userFile.Username,
			Filename:             userFile.Filename,
			FileHash:             userFile.FileHash,
			FileSize:             userFile.FileSize,
			UploadAt:             userFile.UploadAt.String(),
			LastUpdate:           userFile.LastUpdate.String(),
		}

		dataList = append(dataList, data)
	}

	resp.Code = int64(common.RespCodeSuccess.Code)
	resp.Message = common.RespCodeSuccess.Message
	resp.Data = dataList

	return nil
}

func (f FileService) DeleteFile(ctx context.Context, req *fileinterface.DeleteFileReq, resp *fileinterface.DeleteFileResp) error {

	tblFile, err := db.GetFileMeta(req.FileHash)
	if err != nil {
		resp.Code = int64(common.RespCodeNotFoundFileError.Code)
		resp.Message = common.RespCodeNotFoundFileError.Message
		return nil
	}

	bucket := ceph.GetCephBucket(config.Conf.Store.Ceph.FileStoreBucketName)
	err = bucket.Del(config.Conf.Store.Ceph.FilePathPrefix + tblFile.FileSha1)


	if err != nil {
		resp.Code = int64(common.RespCodeRemoveFileError.Code)
		resp.Message = common.RespCodeRemoveFileError.Message
		return nil
	}

	db.DeleteFileMeta(req.FileHash)

	resp.Code = int64(common.RespCodeSuccess.Code)
	resp.Message = common.RespCodeSuccess.Message

	return nil
}
