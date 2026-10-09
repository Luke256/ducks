package repository

import (
	"errors"
	"github.com/Luke256/ducks/model"

	"github.com/google/uuid"
)

const MaxPosterImages = 10

var ErrInvalidPosterImages = errors.New("invalid poster images")

type PosterRepository interface {
	// RegisterPoster ポスターを登録します
	// 登録に成功した場合、ポスターIDを返します
	RegisterPoster(festivalID uuid.UUID, posterName, description string, imageIDs []string) (model.Poster, error)

	// GetPostersByFestivalID イベントIDからポスター一覧を取得します
	GetPostersByFestivalID(festivalID uuid.UUID) ([]model.Poster, error)

	// GetPosterByID ポスターIDからポスターを取得します
	GetPosterByID(posterID uuid.UUID) (model.Poster, error)

	// GetPosterByFestivalIDAndPosterName イベントIDとポスター名からポスターを取得します
	GetPosterByFestivalIDAndPosterName(festivalID uuid.UUID, posterName string) (model.Poster, error)

	// UpdatePoster ポスター情報を更新します
	UpdatePoster(posterID uuid.UUID, posterName, description string) error

	// UpdatePosterImages 画像の追加と削除をまとめて反映します
	UpdatePosterImages(posterID uuid.UUID, addIDs, deleteIDs []string) ([]model.PosterImage, error)

	// UpdatePosterStatus ポスターのステータスを更新します
	UpdatePosterStatus(posterID uuid.UUID, status string) error

	// DeletePoster ポスターを削除します
	DeletePoster(posterID uuid.UUID) error
}
