package storage

import (
	"io"
	"mime/multipart"
)

type Storage interface {
	// UnloadFile ファイルをアップロードし、そのファイル名を返します
	UnloadFile(fileHeader *multipart.FileHeader) (string, error)

	// DownloadFile ファイル名をもとにファイルを取得します
	DownloadFile(fileName string) (io.ReadSeekCloser, error)

	// DeleteFile ファイル名をもとにファイルを削除します
	DeleteFile(fileName string) error
}