package poster

import (
	"errors"
	"mime/multipart"

	"github.com/Luke256/ducks/service/festival"
	"github.com/google/uuid"
)

const (
	PosterStatusUnCollected = "uncollected"
	PosterStatusCollected   = "collected"
	PosterStatusLost        = "lost"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type Poster struct {
	ID          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ImageURL    string            `json:"image_url"`
	Status      string            `json:"status"`
	Festival    festival.Festival `json:"festival"`
}

type Manager interface {
	// Create ポスターを作成します
	Create(name string, festivalID uuid.UUID, description string, image *multipart.FileHeader) (Poster, error)

	// Get 指定されたIDのポスターを取得します
	Get(id uuid.UUID) (Poster, error)

	// GetByFestival 指定されたイベントIDのポスターを取得します
	GetByFestival(festivalID uuid.UUID) ([]Poster, error)

	// GetByName 指定されたイベントの、ポスター名でポスターを取得します
	GetByName(festivalID uuid.UUID, name string) (Poster, error)

	// Edit 指定されたIDのポスター情報を更新します
	Edit(id uuid.UUID, name, description string) error

	// ChangeStatus 指定されたIDのポスターのステータスを変更します
	ChangeStatus(id uuid.UUID, status string) error

	// Delete 指定されたIDのポスターを削除します
	Delete(id uuid.UUID) error
}
