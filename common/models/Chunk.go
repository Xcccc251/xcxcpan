package models

import "gorm.io/gorm"

type Chunk struct {
	gorm.Model
	FileId   string `gorm:"column:file_id" json:"fileId"`
	ChunkId  string `gorm:"column:chunk_id" json:"chunkId"`
	ServerId int    `gorm:"column:server_id" json:"serverId"`
}

func (Chunk) TableName() string {
	return "chunk_info"
}
