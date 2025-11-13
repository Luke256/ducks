package storage

import (
	"errors"
	"io"
	"mime/multipart"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

type Storage interface {
	// UploadFile ファイルをアップロードし、そのファイル名を返します
	UploadFile(fileHeader *multipart.FileHeader) (string, error)

	// DeleteFile ファイル名をもとにファイルを削除します
	DeleteFile(fileName string) error

	// DownloadFile ファイル名をもとにファイルをダウンロードします
	DownloadFile(fileName string) (io.ReadSeekCloser, error)

	// GetFileURL ファイル名をもとにファイルのURLを取得します
	GetFileURL(fileName string) string
}
