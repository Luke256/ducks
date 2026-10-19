package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"github.com/Luke256/ducks/service/festival"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type festivalDateRepository struct {
	repository.FestivalRepository
	stored model.Festival
	writes int
}

func (r *festivalDateRepository) RegisterFestival(name, description string, date time.Time) (model.Festival, error) {
	r.stored = model.Festival{ID: uuid.New(), Name: name, Description: description, Date: date}
	r.writes++
	return r.stored, nil
}

func (r *festivalDateRepository) UpdateFestival(id uuid.UUID, name, description string, date time.Time) error {
	r.stored = model.Festival{ID: id, Name: name, Description: description, Date: date}
	r.writes++
	return nil
}

func (r *festivalDateRepository) GetFestivalByID(uuid.UUID) (model.Festival, error) {
	return r.stored, nil
}

func TestFestivalDateRequests(t *testing.T) {
	t.Parallel()
	for _, method := range []string{http.MethodPost, http.MethodPut} {
		for _, tc := range []struct {
			name  string
			date  any
			omit  bool
			valid bool
		}{
			{name: "valid", date: "2026-10-18", valid: true},
			{name: "leap day", date: "2028-02-29", valid: true},
			{name: "minimum date", date: "1000-01-01", valid: true},
			{name: "maximum date", date: "9999-12-31", valid: true},
			{name: "missing", omit: true},
			{name: "null", date: nil},
			{name: "empty", date: ""},
			{name: "invalid leap day", date: "2026-02-29"},
			{name: "invalid day", date: "2026-04-31"},
			{name: "wrong format", date: "2026/10/18"},
			{name: "timestamp", date: "2026-10-18T00:00:00Z"},
			{name: "out of database range", date: "0999-12-31"},
			{name: "number", date: 20261018},
		} {
			t.Run(method+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				repo := &festivalDateRepository{}
				h := &Handler{festivalManager: festival.NewManagerImpl(repo)}
				e := echo.New()
				e.POST("/festivals", h.CreateFestival)
				e.PUT("/festivals/:id", h.EditFestival)
				path := "/festivals"
				status := http.StatusCreated
				if method == http.MethodPut {
					path += "/" + uuid.NewString()
					status = http.StatusOK
				}
				body := map[string]any{"name": "Festival", "description": "Description"}
				if !tc.omit {
					body["date"] = tc.date
				}
				payload, err := json.Marshal(body)
				require.NoError(t, err)
				req := httptest.NewRequest(method, path, bytes.NewReader(payload))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)
				if !tc.valid {
					require.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())
					require.Zero(t, repo.writes)
					return
				}
				require.Equal(t, status, rec.Code, rec.Body.String())
				require.Equal(t, 1, repo.writes)
				require.Equal(t, tc.date, repo.stored.Date.Format(time.DateOnly))
				var response map[string]any
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
				require.Equal(t, tc.date.(string)+"T00:00:00Z", response["date"])
			})
		}
	}
}
