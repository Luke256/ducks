package v1

import (
	"encoding/json"
	"testing"
)

func TestRegisterFestivalStockPriceValidation(t *testing.T) {
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
			req := RegisterFestivalStockRequest{FestivalID: "festival"}
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
