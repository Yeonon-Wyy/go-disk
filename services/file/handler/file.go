package handler

import (
	"context"
	"errors"
	"go-disk/common"
	"go-disk/common/rpcinterface/fileinterface"
	"go-disk/db"
	"go-disk/meta"
	"log"
	"os"
)

type FileService struct {
}

func (f FileService) GetFileMeta(ctx context.Context, req *fileinterface.GetFileMetaReq, resp *fileinterface.GetFileMetaResp) error {

	fileMeta := meta.GetFileMetaDB(req.FileHash)
	if fileMeta.FileSha1 == "" {
		resp.Code = int64(common.RespCodeNotFoundFileError.Code)
		resp.Message = common.RespCodeNotFoundFileError.Message
		return errors.New("can't find file")
	}
	resp.Code = int64(common.RespCodeSuccess.Code)
	resp.Message = common.RespCodeSuccess.Message
	resp.Data = &fileinterface.GetFileMetaResp_Data{
		FileSha1: fileMeta.FileSha1,
		Filename: fileMeta.FileName,
		FileSize: fileMeta.FileSize,
		Location: fileMeta.Location,
		UploadAt: fileMeta.UploadAt.String(),
		UpdateAt: fileMeta.UpdateAt.String(),
		Status:   int64(fileMeta.Status),
	}

	return nil
}

func (f FileService) UpdateFileMeta(ctx context.Context, req *fileinterface.UpdateFileMetaReq, resp *fileinterface.UpdateFileMetaResp) error {

	fm := meta.GetFileMetaDB(req.FileHash)
	if fm.FileSha1 == "" {
		resp.Code = int64(common.RespCodeNotFoundFileError.Code)
		resp.Message = common.RespCodeNotFoundFileError.Message
		return errors.New("can't find file")
	}

	if req.Filename != "" {
		fm.FileName = req.Filename
		meta.UpdateFileMetaDB(fm)

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
		return errors.New("can't find file")
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
	resp.DataList = dataList

	return nil
}

func (f FileService) DeleteFile(ctx context.Context, req *fileinterface.DeleteFileReq, resp *fileinterface.DeleteFileResp) error {

	fm := meta.GetFileMetaDB(req.FileHash)
	if fm.FileSha1 == "" {
		resp.Code = int64(common.RespCodeNotFoundFileError.Code)
		resp.Message = common.RespCodeNotFoundFileError.Message
		return errors.New("can't find file")
	}

	err := os.Remove(fm.Location)
	if err != nil {
		resp.Code = int64(common.RespCodeRemoveFileError.Code)
		resp.Message = common.RespCodeRemoveFileError.Message
		return errors.New("remove file error")
	}

	meta.RemoveMetaDB(req.FileHash)

	resp.Code = int64(common.RespCodeSuccess.Code)
	resp.Message = common.RespCodeSuccess.Message

	return nil
}
