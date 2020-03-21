package dao

import "time"

type TableFileDao struct {
	Id       uint       `gorm:"column:id"`
	FileHash string     `gorm:"column:file_sha1"`
	FileName string     `gorm:"column:file_name"`
	FileSize int64      `gorm:"column:file_size"`
	FileAddr string     `gorm:"column:file_addr"`
	CreateAt *time.Time `gorm:"column:create_at"`
	UpdateAt *time.Time `gorm:"column:update_at"`
	Status   int        `gorm:"column:status"`
}

type UserFileDao struct {
	Id         uint       `gorm:"column:id"`
	Username   string     `gorm:"column:user_name"`
	FileHash   string     `gorm:"column:file_sha1"`
	FileSize   int64      `gorm:"column:file_size"`
	FileName   string     `gorm:"column:file_name"`
	FilePath   string     `gorm:"column:file_path"`
	UploadAt   *time.Time `gorm:"column:upload_at"`
	LastUpdate *time.Time `gorm:"column:last_update"`
	Status     int        `gorm:"column:status"`
}

func (UserFileDao) TableName() string {
	return "tbl_user_file"
}

func (TableFileDao) TableName() string {
	return "tbl_file"
}
