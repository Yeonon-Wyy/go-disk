package meta

import (
	mysqldb "go-disk/db"
	"log"
	"time"
)

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt time.Time
	UpdateAt time.Time
	Status int
}

func UploadFileMetaDB(meta FileMeta) {
	mysqldb.OnFileUploadFinished(
		meta.FileSha1,
		meta.FileName,
		meta.Location,
		meta.FileSize,
		meta.Status,
		meta.UploadAt,
		meta.UpdateAt)
}

func UpdateFileMetaDB(meta FileMeta) {
	mysqldb.OnFileUpdateFinished(meta.FileSha1, meta.FileName)
}

func GetFileMetaDB(fileSha1 string) *FileMeta {
	tf, err := mysqldb.GetFileMeta(fileSha1)
	if err != nil {
		log.Printf("get file meta error : %v", err)
	}

	if tf == nil {
		return nil
	}

	return &FileMeta{
		FileSha1: tf.FileSha1,
		FileName: tf.Filename,
		FileSize: tf.FileSize,
		Location: tf.FileLocation,
		UploadAt: tf.CreateAt,
		UpdateAt: tf.UpdateAt,
		Status:   tf.Status,
	}


}

func RemoveMetaDB(fileSha1 string) {
	mysqldb.DeleteFileMeta(fileSha1)
}



