package example

import (
	"jykj-cmbp-dev-platform/server/global"
)

// file struct, 文件结构体
type ExaFile struct {
	global.CmbpModel
	FileName     string
	FileMd5      string
	FilePath     string
	ExaFileChunk []ExaFileChunk
	ChunkTotal   int
	IsFinish     bool
}

// file chunk struct, 切片结构体
type ExaFileChunk struct {
	global.CmbpModel
	ExaFileID       string
	FileChunkNumber int
	FileChunkPath   string
}
