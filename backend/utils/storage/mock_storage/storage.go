package mockstorage

import (
	"io"
	"mime/multipart"

	"github.com/google/uuid"
)

// MockStorage は Storage インターフェースのモック実装
// S3を使わず、ローカルファイルシステムに保存します

type MockStorage struct {}


func (s *MockStorage) UploadFile(fileHeader *multipart.FileHeader) (string, error) {
	return uuid.NewString(), nil
}

func (s *MockStorage) DownloadFile(fileName string) (io.ReadSeekCloser, error) {
	return nil, nil
}

func (s *MockStorage) DeleteFile(fileName string) error {
	return nil
}

func (s *MockStorage) GetFileURL(fileName string) string {
	return "https://www.luke256.dev/favicon.ico"
}
