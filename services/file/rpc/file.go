package rpc

import (
	"context"
	"go-disk/common"
	"go-disk/common/log4disk"
	"go-disk/common/rpcinterface/fileinterface"
	"go-disk/services/file/db"
	"log"
)

type FileService struct {
}

func (f FileService) GetFileMeta(ctx context.Context, req *fileinterface.GetFileMetaReq, resp *fileinterface.GetFileMetaResp) error {

	if existInUserFile := db.ExistByFileHashAndUsername(req.FileHash, req.Username); !existInUserFile {
		resp.Code = int64(common.RespCodeNotFoundFileError.Code)
		resp.Message = common.RespCodeNotFoundFileError.Message
		return nil
	}


	tblFile, err := db.GetFileMeta(req.FileHash)
	if err != nil {
		resp.Code = int64(common.RespCodeNotFoundFileError.Code)
		resp.Message = common.RespCodeNotFoundFileError.Message
		return nil
	}
	resp.Code = int64(common.RespCodeSuccess.Code)
	resp.Message = common.RespCodeSuccess.Message
	resp.Data = &fileinterface.GetFileMetaResp_Data{
		FileSha1: tblFile.FileHash,
		Filename: tblFile.FileName,
		FileSize: tblFile.FileSize,
		Location: tblFile.FileAddr,
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
		tblFile.FileName = req.Filename
		//唯一文件表不需要也不应该修改文件名，文件hash与文件名无关

		//更新到user file关联表
		db.UpdateUserFilename(req.Username, req.FileHash, req.Filename)

		resp.Code = int64(common.RespCodeSuccess.Code)
		resp.Message = common.RespCodeSuccess.Message
		return nil
	}

	log4disk.E("filename cann't equals empty")
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
			Filename:             userFile.FileName,
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

	_, err := db.GetFileMeta(req.FileHash)
	if err != nil {
		resp.Code = int64(common.RespCodeNotFoundFileError.Code)
		resp.Message = common.RespCodeNotFoundFileError.Message
		return nil
	}

	if suc := db.DeleteFileMeta(req.FileHash, req.Filename, req.Username); !suc {
		resp.Code = int64(common.RespCodeRemoveFileError.Code)
		resp.Message = common.RespCodeRemoveFileError.Message
	}

	resp.Code = int64(common.RespCodeSuccess.Code)
	resp.Message = common.RespCodeSuccess.Message

	return nil
}
