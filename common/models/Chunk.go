package models

import "gorm.io/gorm"

type Chunk struct {
	gorm.Model
	FileId     string `gorm:"column:file_id" json:"fileId"`
	ChunkId    string `gorm:"column:chunk_id" json:"chunkId"`
	ServerNode string `gorm:"column:server_node" json:"serverNode"`
}

func (Chunk) TableName() string {
	return "chunk_info"
}
