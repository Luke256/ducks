package storage

import (
	"mime/multipart"
)

type Storage interface {
	// UploadFile ファイルをアップロードし、そのファイル名を返します
	UploadFile(fileHeader *multipart.FileHeader) (string, error)

	// DeleteFile ファイル名をもとにファイルを削除します
	DeleteFile(fileName string) error

	// GetFileURL ファイル名をもとにファイルのURLを取得します
	GetFileURL(fileName string) string
}
