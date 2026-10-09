package poster

import (
	"errors"
	"mime/multipart"

	"github.com/Luke256/ducks/repository"
	"github.com/Luke256/ducks/service/festival"
	"github.com/google/uuid"
)

const (
	PosterStatusUnCollected = "uncollected"
	PosterStatusCollected   = "collected"
	PosterStatusLost        = "lost"
	MaxImageSize            = 10 << 20
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidImages = repository.ErrInvalidPosterImages
	ErrImageTooLarge = errors.New("image exceeds 10 MiB")
)

type Image struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type Poster struct {
	ID          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Images      []Image           `json:"image"`
	Status      string            `json:"status"`
	Festival    festival.Festival `json:"festival"`
}

type Manager interface {
	// Create ポスターを作成します
	Create(name string, festivalID uuid.UUID, description string, images []*multipart.FileHeader) (Poster, error)

	// Get 指定されたIDのポスターを取得します
	Get(id uuid.UUID) (Poster, error)

	// GetByFestival 指定されたイベントIDのポスターを取得します
	GetByFestival(festivalID uuid.UUID) ([]Poster, error)

	// GetByName 指定されたイベントの、ポスター名でポスターを取得します
	GetByName(festivalID uuid.UUID, name string) (Poster, error)

	// Edit 指定されたIDのポスター情報を更新します
	Edit(id uuid.UUID, name, description string) error

	// UpdateImages 画像を追加・削除します
	UpdateImages(id uuid.UUID, images []*multipart.FileHeader, deleteIDs []string) ([]Image, error)

	// ChangeStatus 指定されたIDのポスターのステータスを変更します
	ChangeStatus(id uuid.UUID, status string) error

	// Delete 指定されたIDのポスターを削除します
	Delete(id uuid.UUID) error
}
