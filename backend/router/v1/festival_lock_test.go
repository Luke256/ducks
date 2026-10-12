package v1

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"github.com/stretchr/testify/require"
)

func TestLockedFestivalMutationEndpoints(t *testing.T) {
	t.Parallel()
	env := setup(t, common)
	festival := env.mustCreateFestival(t, "Locked event", "description")
	item := env.mustCreateStockItem(t, "Item", "description", "category")
	stock := env.mustCreateFestivalStock(t, festival.ID, item.ID, 100, "description")
	poster := env.mustCreatePoster(t, festival.ID, "Poster", "description")
	sale := env.mustCreateSaleRecord(t, stock.ID, 2)
	env.mustIncrementVisitorCount(t, festival.ID, 5)
	otherFestival := env.mustCreateFestival(t, "Open event", "description")
	otherStock := env.mustCreateFestivalStock(t, otherFestival.ID, item.ID, 100, "description")
	require.NoError(t, env.DB.Model(&model.Festival{}).Where("id = ?", festival.ID).Update("locked", true).Error)

	festivalPath := "/api/festivals/" + festival.ID.String()
	posterPath := "/api/posters/" + poster.ID.String()
	stockPath := "/api/stocks/" + stock.ID.String()
	itemPath := "/api/items/" + item.ID.String()
	salePath := "/api/sales/" + sale.ID.String()
	visitorPath := "/api/visitors/" + festival.ID.String()
	for _, tc := range []struct {
		name, method, path string
		json               any
		form               map[string]string
		image              bool
	}{
		{name: "edit event", method: http.MethodPut, path: festivalPath, json: map[string]any{"name": "changed", "description": "changed"}},
		{name: "delete event", method: http.MethodDelete, path: festivalPath},
		{name: "create poster", method: http.MethodPost, path: "/api/posters", form: map[string]string{"festival_id": festival.ID.String(), "name": "New poster", "description": "changed"}, image: true},
		{name: "edit poster", method: http.MethodPut, path: posterPath, json: map[string]any{"name": "changed", "description": "changed"}},
		{name: "change poster images", method: http.MethodPatch, path: posterPath + "/images", image: true},
		{name: "change poster status", method: http.MethodPatch, path: posterPath + "/status", json: map[string]any{"status": "collected"}},
		{name: "delete poster", method: http.MethodDelete, path: posterPath},
		{name: "create stock", method: http.MethodPost, path: festivalPath + "/stocks", json: map[string]any{"item_id": item.ID, "price": 100, "for_sale": true}},
		{name: "edit stock", method: http.MethodPut, path: stockPath, json: map[string]any{"description": "changed", "for_sale": false}},
		{name: "delete stock", method: http.MethodDelete, path: stockPath},
		{name: "create mixed sales", method: http.MethodPost, path: "/api/sales", json: map[string]any{"items": []map[string]any{{"stock_id": otherStock.ID, "quantity": 1}, {"stock_id": stock.ID, "quantity": 1}}}},
		{name: "delete sale", method: http.MethodDelete, path: salePath},
		{name: "increment visitors", method: http.MethodPost, path: visitorPath, json: map[string]any{"amount": 1}},
		{name: "decrement visitors", method: http.MethodPost, path: visitorPath, json: map[string]any{"amount": -1}},
		{name: "edit shared item", method: http.MethodPut, path: itemPath, form: map[string]string{"name": "changed", "description": "changed", "category": "changed"}},
		{name: "change shared item image", method: http.MethodPut, path: itemPath + "/image", image: true},
		{name: "delete shared item", method: http.MethodDelete, path: itemPath},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := env.R(t).Request(tc.method, tc.path)
			if tc.json != nil {
				req.WithJSON(tc.json)
			}
			if tc.form != nil || tc.image {
				req.WithMultipart()
				if tc.form != nil {
					req.WithForm(tc.form)
				}
				if tc.image {
					req.WithFile("image", "new.png", strings.NewReader(""))
				}
			}
			req.Expect().Status(http.StatusForbidden).JSON().Object().Value("message").String().IsEqual(repository.ErrEventLocked.Error())
		})
	}

	for _, path := range []string{
		"/api/festivals", festivalPath, festivalPath + "/posters", posterPath,
		"/api/posters/" + festival.ID.String() + "/Poster",
		festivalPath + "/stocks", stockPath, itemPath, salePath,
		stockPath + "/sales", visitorPath,
	} {
		t.Run("read "+path, func(t *testing.T) { env.R(t).GET(path).Expect().Status(http.StatusOK) })
	}
	env.R(t).GET("/api/sales").WithQuery("festival_id", festival.ID.String()).Expect().Status(http.StatusOK)
	sales, err := env.Repo.GetSaleRecordsByFestivalStockID(otherStock.ID)
	require.NoError(t, err)
	require.Empty(t, sales)

	require.NoError(t, env.DB.Model(&model.Festival{}).Where("id = ?", festival.ID).Update("locked", false).Error)
	env.R(t).PUT(festivalPath).WithJSON(map[string]any{"name": "Unlocked event", "description": "updated"}).Expect().Status(http.StatusOK)
}
