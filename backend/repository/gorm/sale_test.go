package gorm

import (
	"testing"
	"time"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestSaleRecordUTC(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)
	fixed := time.Date(2026, 10, 8, 0, 5, 12, 123000000, time.FixedZone("JST", 9*60*60))
	db := repo.db.Session(&gorm.Session{NowFunc: func() time.Time { return fixed }})
	utcRepo, _, err := NewGormRepository(db, false)
	require.NoError(t, err)
	require.Equal(t, fixed, db.NowFunc(), "repository setup must not mutate the caller's clock")

	fes := mustCreateFestival(t, utcRepo, "UTC Festival", "UTC test")
	item := mustCreateStockItem(t, utcRepo, "UTC Item", "UTC test", "Test", "")
	stock := mustCreateFestivalStock(t, utcRepo, fes.ID, item.ID, 100, "")
	created := mustCreateSaleRecord(t, utcRepo, stock.ID, 1)
	require.Equal(t, fixed.UTC(), created.CreatedAt)

	got, err := utcRepo.GetSaleRecordByID(created.ID)
	require.NoError(t, err)
	byStock, err := utcRepo.GetSaleRecordsByFestivalStockID(stock.ID)
	require.NoError(t, err)
	queried, err := utcRepo.QuerySaleRecords(fes.ID, item.ID)
	require.NoError(t, err)
	for _, records := range [][]model.SaleRecord{{got}, byStock, queried} {
		require.Len(t, records, 1)
		require.Equal(t, created.CreatedAt, records[0].CreatedAt)
	}

	var stored string
	err = utcRepo.db.Raw("SELECT DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s.%f') FROM sale_records WHERE id = ?", created.ID).Scan(&stored).Error
	require.NoError(t, err)
	require.Equal(t, "2026-10-07 15:05:12.123000", stored)
	var sessionZone string
	require.NoError(t, utcRepo.db.Raw("SELECT @@session.time_zone").Scan(&sessionZone).Error)
	require.Equal(t, "+00:00", sessionZone)
}

func TestCreateSaleRecord(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	fes := mustCreateFestival(t, repo, "Test Festival", "A festival for testing")
	stockItem := mustCreateStockItem(t, repo, "Test Stock Item", "An item for testing", "Test Category", "")
	fesStock := mustCreateFestivalStock(t, repo, fes.ID, stockItem.ID, 100, "Stock Description")

	t.Run("Create Sale Record", func(t *testing.T) {
		t.Parallel()
		saleRecord, err := repo.CreateSaleRecords(repository.SaleData{
			FestivalStockID: fesStock.ID,
			Quantity:        5,
		})
		assert.NoError(t, err)
		assert.Equal(t, fesStock.ID, saleRecord[0].FestivalStockID)
		assert.Equal(t, 5, saleRecord[0].Quantity)
	})

	t.Run("Create Multiple Sale Records", func(t *testing.T) {
		t.Parallel()
		saleRecords, err := repo.CreateSaleRecords(
			repository.SaleData{
				FestivalStockID: fesStock.ID,
				Quantity:        3,
			},
			repository.SaleData{
				FestivalStockID: fesStock.ID,
				Quantity:        7,
			},
		)
		assert.NoError(t, err)
		assert.Len(t, saleRecords, 2)
		assert.Equal(t, 3, saleRecords[0].Quantity)
		assert.Equal(t, 7, saleRecords[1].Quantity)
	})

	t.Run("Create Sale Record with Non-Existent Festival Stock", func(t *testing.T) {
		t.Parallel()
		_, err := repo.CreateSaleRecords(repository.SaleData{
			FestivalStockID: uuid.New(),
			Quantity:        5,
		})
		assert.Error(t, err)
		assert.Equal(t, repository.ErrForeignKey, err)
	})
}

func TestGetSaleRecordByID(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	fes := mustCreateFestival(t, repo, "Test Festival", "A festival for testing")
	stockItem := mustCreateStockItem(t, repo, "Test Stock Item", "An item for testing", "Test Category", "")
	fesStock := mustCreateFestivalStock(t, repo, fes.ID, stockItem.ID, 100, "Stock Description")
	saleRecord := mustCreateSaleRecord(t, repo, fesStock.ID, 10)

	t.Run("Get Sale Record By ID", func(t *testing.T) {
		t.Parallel()
		got, err := repo.GetSaleRecordByID(saleRecord.ID)
		assert.NoError(t, err)
		assert.Equal(t, saleRecord.ID, got.ID)
		assert.Equal(t, saleRecord.FestivalStockID, got.FestivalStockID)
		assert.Equal(t, saleRecord.Quantity, got.Quantity)
	})

	t.Run("Get Non-Existent Sale Record By ID", func(t *testing.T) {
		t.Parallel()
		_, err := repo.GetSaleRecordByID(uuid.New())
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestGetSaleRecordsByFestivalStockID(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	fes := mustCreateFestival(t, repo, "Test Festival", "A festival for testing")
	stockItem := mustCreateStockItem(t, repo, "Test Stock Item", "An item for testing", "Test Category", "")
	fesStock := mustCreateFestivalStock(t, repo, fes.ID, stockItem.ID, 100, "Stock Description")
	fesStock2 := mustCreateFestivalStock(t, repo, fes.ID, stockItem.ID, 200, "Stock Description")

	saleRecord1 := mustCreateSaleRecord(t, repo, fesStock.ID, 10)
	saleRecord2 := mustCreateSaleRecord(t, repo, fesStock.ID, 20)
	mustCreateSaleRecord(t, repo, fesStock2.ID, 30)

	t.Run("Get Sale Records By Festival Stock ID", func(t *testing.T) {
		t.Parallel()
		records, err := repo.GetSaleRecordsByFestivalStockID(fesStock.ID)
		assert.NoError(t, err)
		assert.Len(t, records, 2)

		var recordIDs []uuid.UUID
		for _, r := range records {
			recordIDs = append(recordIDs, r.ID)
		}
		assert.Contains(t, recordIDs, saleRecord1.ID)
		assert.Contains(t, recordIDs, saleRecord2.ID)
	})

	t.Run("Get Sale Records By Non-Existent Festival Stock ID", func(t *testing.T) {
		t.Parallel()
		_, err := repo.GetSaleRecordsByFestivalStockID(uuid.New())
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestQuerySaleRecords(t *testing.T) {
	t.Parallel()
	repo := setup(t, s3)

	fes1 := mustCreateFestival(t, repo, "Festival One", "First festival")
	fes2 := mustCreateFestival(t, repo, "Festival Two", "Second festival")

	itemA := mustCreateStockItem(t, repo, "Item A", "First item", "Category 1", "")
	itemB := mustCreateStockItem(t, repo, "Item B", "Second item", "Category 2", "")

	fes1StockA := mustCreateFestivalStock(t, repo, fes1.ID, itemA.ID, 150, "Stock Description A")
	fes1StockB := mustCreateFestivalStock(t, repo, fes1.ID, itemB.ID, 200, "Stock Description B")
	fes2StockA := mustCreateFestivalStock(t, repo, fes2.ID, itemA.ID, 250, "Stock Description A")

	saleRecord1 := mustCreateSaleRecord(t, repo, fes1StockA.ID, 3)
	saleRecord2 := mustCreateSaleRecord(t, repo, fes1StockB.ID, 5)
	saleRecord3 := mustCreateSaleRecord(t, repo, fes2StockA.ID, 7)

	t.Run("Query All Sale Records", func(t *testing.T) {
		t.Parallel()
		records, err := repo.QuerySaleRecords(uuid.Nil, uuid.Nil)
		assert.NoError(t, err)
		assert.Len(t, records, 3)
	})

	t.Run("Query Sale Records by Festival ID", func(t *testing.T) {
		t.Parallel()
		records, err := repo.QuerySaleRecords(fes1.ID, uuid.Nil)
		assert.NoError(t, err)
		assert.Len(t, records, 2)

		// Check that the correct records are returned
		var recordIDs []uuid.UUID
		for _, r := range records {
			recordIDs = append(recordIDs, r.ID)
		}
		assert.Contains(t, recordIDs, saleRecord1.ID)
		assert.Contains(t, recordIDs, saleRecord2.ID)
	})

	t.Run("Query Sale Records by Stock Item ID", func(t *testing.T) {
		t.Parallel()
		records, err := repo.QuerySaleRecords(uuid.Nil, itemA.ID)
		assert.NoError(t, err)
		assert.Len(t, records, 2)

		// Check that the correct records are returned
		var recordIDs []uuid.UUID
		for _, r := range records {
			recordIDs = append(recordIDs, r.ID)
		}
		assert.Contains(t, recordIDs, saleRecord1.ID)
		assert.Contains(t, recordIDs, saleRecord3.ID)
	})

	t.Run("Query Sale Records by Festival ID and Stock Item ID", func(t *testing.T) {
		t.Parallel()
		records, err := repo.QuerySaleRecords(fes1.ID, itemB.ID)
		assert.NoError(t, err)
		assert.Len(t, records, 1)
		assert.Equal(t, saleRecord2.ID, records[0].ID)
	})
}

func TestDeleteSaleRecord(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	fes := mustCreateFestival(t, repo, "Test Festival", "A festival for testing")
	stockItem := mustCreateStockItem(t, repo, "Test Stock Item", "An item for testing", "Test Category", "")
	fesStock := mustCreateFestivalStock(t, repo, fes.ID, stockItem.ID, 100, "Stock Description")
	saleRecord := mustCreateSaleRecord(t, repo, fesStock.ID, 10)

	t.Run("Delete Sale Record", func(t *testing.T) {
		t.Parallel()
		err := repo.DeleteSaleRecord(saleRecord.ID)
		assert.NoError(t, err)

		_, err = repo.GetSaleRecordByID(saleRecord.ID)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("Delete Non-Existent Sale Record", func(t *testing.T) {
		t.Parallel()
		err := repo.DeleteSaleRecord(uuid.New())
		assert.Equal(t, repository.ErrNotFound, err)
	})
}
