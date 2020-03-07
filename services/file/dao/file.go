package dao

import "time"

type TableFile struct {
	FileSha1 string
	Filename string
	FileSize int64
	FileLocation string
	CreateAt time.Time
	UpdateAt time.Time
	Status int
}

type UserFileDao struct {
	Username string
	Filename string
	FileHash string
	FileSize int64
	UploadAt time.Time
	LastUpdate time.Time
}

