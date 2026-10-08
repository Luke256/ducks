package stockitem

import (
	"errors"
	"mime/multipart"
	"testing"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"github.com/Luke256/ducks/utils/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type imageRepository struct {
	repository.Repository
	item      model.StockItem
	updateErr error
}

func (r *imageRepository) GetStockItemByID(uuid.UUID) (model.StockItem, error) {
	return r.item, nil
}

func (r *imageRepository) UpdateStockItem(id uuid.UUID, name, description, category, imageID string) (model.StockItem, error) {
	if r.updateErr != nil {
		return model.StockItem{}, r.updateErr
	}
	r.item.ImageID = imageID
	return r.item, nil
}

type imageStorage struct {
	storage.Storage
	files     map[string]bool
	uploadErr error
	deleteErr error
}

func (s *imageStorage) UploadFile(*multipart.FileHeader) (string, error) {
	if s.uploadErr != nil {
		return "", s.uploadErr
	}
	s.files["new"] = true
	return "new", nil
}

func (s *imageStorage) DeleteFile(id string) error {
	if s.deleteErr != nil {
		return s.deleteErr
	}
	delete(s.files, id)
	return nil
}

func TestUpdateImagePreservesStoredImageOnFailure(t *testing.T) {
	failure := errors.New("test failure")
	for _, tc := range []struct {
		name                            string
		uploadErr, updateErr, deleteErr error
		wantErr                         bool
		wantImage                       string
		wantFiles                       map[string]bool
	}{
		{name: "success", wantImage: "new", wantFiles: map[string]bool{"new": true}},
		{name: "upload failure", uploadErr: failure, wantErr: true, wantImage: "old", wantFiles: map[string]bool{"old": true}},
		{name: "database failure", updateErr: failure, wantErr: true, wantImage: "old", wantFiles: map[string]bool{"old": true}},
		{name: "cleanup failure", deleteErr: failure, wantImage: "new", wantFiles: map[string]bool{"old": true, "new": true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &imageRepository{item: model.StockItem{ID: uuid.New(), ImageID: "old"}, updateErr: tc.updateErr}
			files := &imageStorage{files: map[string]bool{"old": true}, uploadErr: tc.uploadErr, deleteErr: tc.deleteErr}
			err := NewManagerImpl(repo, files).UpdateImage(repo.item.ID, &multipart.FileHeader{})
			if tc.wantErr {
				require.ErrorIs(t, err, failure)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.wantImage, repo.item.ImageID)
			require.Equal(t, tc.wantFiles, files.files)
		})
	}
}
