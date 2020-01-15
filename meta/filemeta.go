package meta

import "log"

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetaMaps = make(map[string]FileMeta)

func UpdateFileMeta(fileSha1 string, meta FileMeta) {
	fileMetaMaps[fileSha1] = meta
}

func GetFileMeta(fileSha1 string) *FileMeta {
	if meta, ok := fileMetaMaps[fileSha1]; ok {
		return &meta
	}
	log.Printf("can't get file meta of : %s", fileSha1)
	return nil
}

func RemoveMeta(fileSha1 string) {
	delete(fileMetaMaps, fileSha1)
}



