package dao_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"payment-gateway/cmd/infra/dao"
)

func TestOrderDao_FindById(t *testing.T) {
	t.Run("should find payment by ID successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		now := time.Now()
		expectedID := int64(1)
		rows := sqlmock.NewRows([]string{"id", "status", "amount", "created_at", "updated_at"}).
			AddRow(expectedID, "approved", 100.5, now, now)

		mock.ExpectQuery(`SELECT id, status, amount, created_at, updated_at FROM orders WHERE id = ?`).
			WithArgs(expectedID).
			WillReturnRows(rows)

		dao := dao.NewOrderDao(db)
		result, err := dao.FindById(expectedID)

		assert.NoError(t, err)
		if assert.NotNil(t, result) {
			assert.Equal(t, expectedID, result.Id())
			assert.Equal(t, "approved", result.Status())
			assert.Equal(t, 100.5, result.Amount())
			assert.NotNil(t, result.CreatedAt())
			assert.NotNil(t, result.UpdatedAt())
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when query fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		expectedID := int64(1)
		mock.ExpectQuery(`SELECT id, status, amount, created_at, updated_at FROM orders WHERE id = ?`).
			WithArgs(expectedID).
			WillReturnError(assert.AnError)

		dao := dao.NewOrderDao(db)
		result, err := dao.FindById(expectedID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return empty entity when no rows found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		expectedID := int64(1)
		rows := sqlmock.NewRows([]string{"id", "status", "amount", "created_at", "updated_at"})

		mock.ExpectQuery(`SELECT id, status, amount, created_at, updated_at FROM orders WHERE id = ?`).
			WithArgs(expectedID).
			WillReturnRows(rows)

		dao := dao.NewOrderDao(db)
		result, err := dao.FindById(expectedID)

		assert.NoError(t, err)
		if assert.NotNil(t, result) {
			assert.Equal(t, int64(0), result.Id())
			assert.Equal(t, "", result.Status())
			assert.Equal(t, 0.0, result.Amount())
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when scan fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		expectedID := int64(1)
		rows := sqlmock.NewRows([]string{"id", "order_id"}).
			AddRow(expectedID, 123)

		mock.ExpectQuery(`SELECT id, status, amount, created_at, updated_at FROM orders WHERE id = ?`).
			WithArgs(expectedID).
			WillReturnRows(rows)

		dao := dao.NewOrderDao(db)
		result, err := dao.FindById(expectedID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
