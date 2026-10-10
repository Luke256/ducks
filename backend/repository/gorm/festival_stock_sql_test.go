package gorm

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TestRegisterFestivalStockIncludesZeroPrice(t *testing.T) {
	t.Parallel()
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root@tcp(localhost:3307)/ducks",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun:                 true,
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)

	var values clause.Values
	require.NoError(t, db.Callback().Create().After("gorm:create").Register("test:capture_values", func(tx *gorm.DB) {
		values = tx.Statement.Clauses["VALUES"].Expression.(clause.Values)
	}))

	repo := &GormRepository{db: db}
	_, err = repo.RegisterFestivalStock(uuid.New(), uuid.New(), 0, "Free item", true)
	require.NoError(t, err)
	require.Len(t, values.Values, 1)
	for i, column := range values.Columns {
		if column.Name == "price" {
			require.Equal(t, 0, values.Values[0][i])
			return
		}
	}
	t.Fatal("INSERT does not include the price column")
}
