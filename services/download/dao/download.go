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

func (TableFileDao) TableName() string {
	return "tbl_file"
}
