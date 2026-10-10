package v1

import (
	"encoding/json"
	"testing"
)

func TestRegisterFestivalStockPriceValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		body  string
		valid bool
	}{
		{"positive", `{"item_id":"item", "price":1500}`, true},
		{"zero", `{"item_id":"item", "price":0}`, true},
		{"negative", `{"item_id":"item", "price":-1}`, false},
		{"missing", `{"item_id":"item"}`, false},
		{"null", `{"item_id":"item", "price":null}`, false},
		{"fractional", `{"item_id":"item", "price":0.5}`, false},
		{"string", `{"item_id":"item", "price":"0"}`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := RegisterFestivalStockRequest{FestivalID: "festival", ForSale: new(bool)}
			err := json.Unmarshal([]byte(tc.body), &req)
			if err == nil {
				err = req.Validate()
			}
			if (err == nil) != tc.valid {
				t.Fatalf("valid = %v, want %v (error: %v)", err == nil, tc.valid, err)
			}
		})
	}
}

func TestFestivalStockForSaleValidation(t *testing.T) {
	t.Parallel()
	env := setup(t, common)
	festival := env.mustCreateFestival(t, "ForSale validation", "test")
	item := env.mustCreateStockItem(t, "ForSale item", "test", "test")

	for _, tc := range []struct {
		name  string
		value any
		valid bool
	}{
		{"true", true, true},
		{"false", false, true},
		{"missing", nil, false},
		{"null", nil, false},
		{"string", "false", false},
		{"number", 0, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			e := env.R(t)
			payload := map[string]any{"item_id": item.ID, "price": 100, "description": "New description"}
			if tc.name != "missing" {
				payload["for_sale"] = tc.value
			}
			status := 400
			if tc.valid {
				status = 201
			}
			created := e.POST("/api/festivals/{festival_id}/stocks", festival.ID).
				WithJSON(payload).Expect().Status(status)
			if tc.valid {
				stock := created.JSON().Object()
				stock.Value("for_sale").IsEqual(tc.value)
				e.GET("/api/stocks/{festival_stock_id}", stock.Value("id").String().Raw()).
					Expect().Status(200).JSON().Object().Value("for_sale").IsEqual(tc.value)
			}

			stock := env.mustCreateFestivalStock(t, festival.ID, item.ID, 100, "Old description")
			update := map[string]any{"description": "New description"}
			if tc.name != "missing" {
				update["for_sale"] = tc.value
			}
			if tc.valid {
				status = 204
			}
			e.PUT("/api/stocks/{festival_stock_id}", stock.ID).
				WithJSON(update).Expect().Status(status)
			fetched := e.GET("/api/stocks/{festival_stock_id}", stock.ID).
				Expect().Status(200).JSON().Object()
			if tc.valid {
				fetched.Value("for_sale").IsEqual(tc.value)
				fetched.Value("description").IsEqual("New description")
			} else {
				fetched.Value("for_sale").IsEqual(true)
				fetched.Value("description").IsEqual("Old description")
			}
		})
	}
}
