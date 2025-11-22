package gorm

import (
	"testing"

	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateSaleRecord(t *testing.T) {
	repo := setup(t, common)

	fes := mustCreateFestival(t, repo, "Test Festival", "A festival for testing")
	stockItem := mustCreateStockItem(t, repo, "Test Stock Item", "An item for testing", "Test Category", "")
	fesStock := mustCreateFestivalStock(t, repo, fes.ID, stockItem.ID, 100)

	t.Run("Create Sale Record", func(t *testing.T) {
		saleRecord, err := repo.CreateSaleRecord(fesStock.ID, 5)
		assert.NoError(t, err)
		assert.Equal(t, fesStock.ID, saleRecord.FestivalStockID)
		assert.Equal(t, 5, saleRecord.Quantity)
	})
}

func TestGetSaleRecordByID(t *testing.T) {
	repo := setup(t, common)

	fes := mustCreateFestival(t, repo, "Test Festival", "A festival for testing")
	stockItem := mustCreateStockItem(t, repo, "Test Stock Item", "An item for testing", "Test Category", "")
	fesStock := mustCreateFestivalStock(t, repo, fes.ID, stockItem.ID, 100)
	saleRecord := mustCreateSaleRecord(t, repo, fesStock.ID, 10)

	t.Run("Get Sale Record By ID", func(t *testing.T) {
		got, err := repo.GetSaleRecordByID(saleRecord.ID)
		assert.NoError(t, err)
		assert.Equal(t, saleRecord.ID, got.ID)
		assert.Equal(t, saleRecord.FestivalStockID, got.FestivalStockID)
		assert.Equal(t, saleRecord.Quantity, got.Quantity)
	})

	t.Run("Get Non-Existent Sale Record By ID", func(t *testing.T) {
		_, err := repo.GetSaleRecordByID(uuid.New())
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestGetSaleRecordsByFestivalStockID(t *testing.T) {
	repo := setup(t, common)

	fes := mustCreateFestival(t, repo, "Test Festival", "A festival for testing")
	stockItem := mustCreateStockItem(t, repo, "Test Stock Item", "An item for testing", "Test Category", "")
	fesStock := mustCreateFestivalStock(t, repo, fes.ID, stockItem.ID, 100)
	fesStock2 := mustCreateFestivalStock(t, repo, fes.ID, stockItem.ID, 200)

	saleRecord1 := mustCreateSaleRecord(t, repo, fesStock.ID, 10)
	saleRecord2 := mustCreateSaleRecord(t, repo, fesStock.ID, 20)
	mustCreateSaleRecord(t, repo, fesStock2.ID, 30)

	t.Run("Get Sale Records By Festival Stock ID", func(t *testing.T) {
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
		_, err := repo.GetSaleRecordsByFestivalStockID(uuid.New())
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestQuerySaleRecords(t *testing.T) {
	repo := setup(t, s3)

	fes1 := mustCreateFestival(t, repo, "Festival One", "First festival")
	fes2 := mustCreateFestival(t, repo, "Festival Two", "Second festival")

	itemA := mustCreateStockItem(t, repo, "Item A", "First item", "Category 1", "")
	itemB := mustCreateStockItem(t, repo, "Item B", "Second item", "Category 2", "")

	fes1StockA := mustCreateFestivalStock(t, repo, fes1.ID, itemA.ID, 150)
	fes1StockB := mustCreateFestivalStock(t, repo, fes1.ID, itemB.ID, 200)
	fes2StockA := mustCreateFestivalStock(t, repo, fes2.ID, itemA.ID, 250)

	saleRecord1 := mustCreateSaleRecord(t, repo, fes1StockA.ID, 3)
	saleRecord2 := mustCreateSaleRecord(t, repo, fes1StockB.ID, 5)
	saleRecord3 := mustCreateSaleRecord(t, repo, fes2StockA.ID, 7)

	t.Run("Query All Sale Records", func(t *testing.T) {
		records, err := repo.QuerySaleRecords(uuid.Nil, uuid.Nil)
		assert.NoError(t, err)
		assert.Len(t, records, 3)
	})

	t.Run("Query Sale Records by Festival ID", func(t *testing.T) {
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
		records, err := repo.QuerySaleRecords(fes1.ID, itemB.ID)
		assert.NoError(t, err)
		assert.Len(t, records, 1)
		assert.Equal(t, saleRecord2.ID, records[0].ID)
	})
}

func TestDeleteSaleRecord(t *testing.T) {
	repo := setup(t, common)

	fes := mustCreateFestival(t, repo, "Test Festival", "A festival for testing")
	stockItem := mustCreateStockItem(t, repo, "Test Stock Item", "An item for testing", "Test Category", "")
	fesStock := mustCreateFestivalStock(t, repo, fes.ID, stockItem.ID, 100)
	saleRecord := mustCreateSaleRecord(t, repo, fesStock.ID, 10)

	t.Run("Delete Sale Record", func(t *testing.T) {
		err := repo.DeleteSaleRecord(saleRecord.ID)
		assert.NoError(t, err)

		_, err = repo.GetSaleRecordByID(saleRecord.ID)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("Delete Non-Existent Sale Record", func(t *testing.T) {
		err := repo.DeleteSaleRecord(uuid.New())
		assert.Equal(t, repository.ErrNotFound, err)
	})
}
