package poster

import (
	"errors"
	"mime/multipart"
	"slices"
	"strings"
	"testing"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"github.com/Luke256/ducks/utils/compressor"
	"github.com/Luke256/ducks/utils/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type imageRepository struct {
	repository.Repository
	poster                                    model.Poster
	getErr, registerErr, updateErr, deleteErr error
}

func (r *imageRepository) GetPosterByID(uuid.UUID) (model.Poster, error) {
	return r.poster, r.getErr
}

func (r *imageRepository) GetPosterByFestivalIDAndPosterName(uuid.UUID, string) (model.Poster, error) {
	return model.Poster{}, repository.ErrNotFound
}

func (r *imageRepository) GetFestivalByID(id uuid.UUID) (model.Festival, error) {
	return model.Festival{ID: id}, nil
}

func (r *imageRepository) RegisterPoster(festivalID uuid.UUID, name, description string, ids []string) (model.Poster, error) {
	if r.registerErr != nil {
		return model.Poster{}, r.registerErr
	}
	r.poster = model.Poster{ID: uuid.New(), FestivalID: festivalID, PosterName: name, Description: description}
	for _, id := range ids {
		r.poster.Images = append(r.poster.Images, model.PosterImage{ID: id, PosterID: r.poster.ID})
	}
	return r.poster, nil
}

func (r *imageRepository) UpdatePosterImages(id uuid.UUID, addIDs, deleteIDs []string) ([]model.PosterImage, error) {
	if r.updateErr != nil {
		return nil, r.updateErr
	}
	images := make([]model.PosterImage, 0)
	for _, img := range r.poster.Images {
		if !slices.Contains(deleteIDs, img.ID) {
			images = append(images, img)
		}
	}
	for _, imageID := range addIDs {
		images = append(images, model.PosterImage{ID: imageID, PosterID: id})
	}
	r.poster.Images = images
	return images, nil
}

func (r *imageRepository) DeletePoster(uuid.UUID) error {
	if r.deleteErr == nil {
		r.poster.Images = nil
	}
	return r.deleteErr
}

type imageStorage struct {
	storage.Storage
	files      map[string]bool
	uploads    int
	failUpload int
	uploadErr  error
	failDelete bool
}

func (s *imageStorage) UploadFile(file *multipart.FileHeader) (string, error) {
	s.uploads++
	if s.uploads == s.failUpload {
		return "", s.uploadErr
	}
	s.files[file.Filename] = true
	return file.Filename, nil
}

func (s *imageStorage) DeleteFile(id string) error {
	if s.failDelete {
		return errors.New("cleanup failure")
	}
	delete(s.files, id)
	return nil
}

func (s *imageStorage) GetFileURL(id string) string {
	return "/api/v1/images/" + id
}

func TestUpdateImagesPreservesFilesOnFailure(t *testing.T) {
	t.Parallel()
	failure := errors.New("test failure")
	for _, tc := range []struct {
		name              string
		deleteIDs         []string
		failUpload        int
		updateErr, getErr error
		failDelete        bool
		wantErr           error
		wantIDs           []string
		wantFiles         map[string]bool
	}{
		{name: "replace", deleteIDs: []string{"old1"}, wantIDs: []string{"old2", "new1", "new2"}, wantFiles: map[string]bool{"old2": true, "new1": true, "new2": true}},
		{name: "append", wantIDs: []string{"old1", "old2", "new1", "new2"}, wantFiles: map[string]bool{"old1": true, "old2": true, "new1": true, "new2": true}},
		{name: "second upload fails", deleteIDs: []string{"old1"}, failUpload: 2, wantErr: failure},
		{name: "database fails", deleteIDs: []string{"old1"}, updateErr: failure, wantErr: failure},
		{name: "foreign image", deleteIDs: []string{"foreign"}, wantErr: ErrInvalidImages},
		{name: "duplicate deletion", deleteIDs: []string{"old1", "old1"}, wantErr: ErrInvalidImages},
		{name: "missing poster", getErr: repository.ErrNotFound, wantErr: ErrNotFound},
		{name: "cleanup fails after commit", deleteIDs: []string{"old1"}, failDelete: true, wantIDs: []string{"old2", "new1", "new2"}, wantFiles: map[string]bool{"old1": true, "old2": true, "new1": true, "new2": true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			id := uuid.New()
			repo := &imageRepository{poster: model.Poster{ID: id, Images: []model.PosterImage{{ID: "old1", PosterID: id}, {ID: "old2", PosterID: id}}}, updateErr: tc.updateErr, getErr: tc.getErr}
			files := &imageStorage{files: map[string]bool{"old1": true, "old2": true}, failUpload: tc.failUpload, uploadErr: failure, failDelete: tc.failDelete}
			images, err := NewManagerImpl(repo, files).UpdateImages(id, []*multipart.FileHeader{{Filename: "new1"}, {Filename: "new2"}}, tc.deleteIDs)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				require.Nil(t, images)
				require.Len(t, repo.poster.Images, 2)
				require.Equal(t, map[string]bool{"old1": true, "old2": true}, files.files)
				return
			}
			require.NoError(t, err)
			ids := make([]string, len(images))
			for i, image := range images {
				ids[i] = image.ID
				require.Equal(t, "/api/v1/images/"+image.ID, image.URL)
			}
			require.ElementsMatch(t, tc.wantIDs, ids)
			require.Equal(t, tc.wantFiles, files.files)
		})
	}
}

func TestUpdateImagesRejectsDeletingAllFiles(t *testing.T) {
	t.Parallel()
	for _, ids := range [][]string{{"old1"}, {"old1", "old2"}} {
		t.Run(strings.Join(ids, ","), func(t *testing.T) {
			t.Parallel()
			id := uuid.New()
			repo := &imageRepository{poster: model.Poster{ID: id}}
			files := &imageStorage{files: map[string]bool{}}
			for _, imageID := range ids {
				repo.poster.Images = append(repo.poster.Images, model.PosterImage{ID: imageID, PosterID: id})
				files.files[imageID] = true
			}
			_, err := NewManagerImpl(repo, files).UpdateImages(id, nil, ids)
			require.ErrorIs(t, err, ErrInvalidImages)
			require.Len(t, repo.poster.Images, len(ids))
			for _, imageID := range ids {
				require.True(t, files.files[imageID])
			}
			require.Zero(t, files.uploads)
		})
	}
}

func TestCreateImagesRollsBackUploadedFiles(t *testing.T) {
	t.Parallel()
	failure := errors.New("test failure")
	for _, tc := range []struct {
		name        string
		failUpload  int
		registerErr error
		wantErr     error
	}{
		{name: "success"},
		{name: "second upload fails", failUpload: 2, wantErr: failure},
		{name: "database fails", registerErr: failure, wantErr: failure},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo := &imageRepository{registerErr: tc.registerErr}
			files := &imageStorage{files: map[string]bool{}, failUpload: tc.failUpload, uploadErr: failure}
			poster, err := NewManagerImpl(repo, files).Create("test", uuid.New(), "location", []*multipart.FileHeader{{Filename: "new1"}, {Filename: "new2"}})
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				require.Empty(t, files.files)
				return
			}
			require.NoError(t, err)
			require.Len(t, poster.Images, 2)
			require.Len(t, files.files, 2)
		})
	}
}

func TestInvalidUploads(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		images []*multipart.FileHeader
		want   error
	}{
		{name: "nil file", images: []*multipart.FileHeader{nil}, want: ErrInvalidImages},
		{name: "too many", images: make([]*multipart.FileHeader, repository.MaxPosterImages+1), want: ErrInvalidImages},
		{name: "too large", images: []*multipart.FileHeader{{Size: MaxImageSize + 1}}, want: ErrImageTooLarge},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.ErrorIs(t, validateUploads(tc.images), tc.want)
		})
	}
	files := &imageStorage{files: map[string]bool{}, failUpload: 1, uploadErr: compressor.ErrInvalidImage}
	_, err := NewManagerImpl(&imageRepository{}, files).Create("test", uuid.New(), "location", []*multipart.FileHeader{{Filename: "bad"}})
	require.ErrorIs(t, err, ErrInvalidImages)
}

func TestDeletePosterPreservesFilesOnDatabaseFailure(t *testing.T) {
	t.Parallel()
	failure := errors.New("test failure")
	repo := &imageRepository{poster: model.Poster{Images: []model.PosterImage{{ID: "old"}}}, deleteErr: failure}
	files := &imageStorage{files: map[string]bool{"old": true}}
	require.ErrorIs(t, NewManagerImpl(repo, files).Delete(uuid.New()), failure)
	require.Equal(t, map[string]bool{"old": true}, files.files)
}
