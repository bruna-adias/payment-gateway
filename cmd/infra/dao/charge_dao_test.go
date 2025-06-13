package dao_test

import (
	"payment-gateway/cmd/domain/charge"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"payment-gateway/cmd/infra/dao"
)

func TestChargeDao_Insert(t *testing.T) {
	chargeEntity := charge.NewChargeBuilder().
		WithAmount(10.0).
		WithCategory("financial_fee").
		WithPaymentId(1).
		WithCreatedAt(time.Now()).
		WithUpdatedAt(time.Now()).
		Build()

	t.Run("should insert charge successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		createdAt := chargeEntity.CreatedAt().Format("2006-01-02 15:04:05")
		updatedAt := chargeEntity.UpdatedAt().Format("2006-01-02 15:04:05")

		mock.ExpectExec(`INSERT INTO charges`).
			WithArgs(
				chargeEntity.Amount(),
				chargeEntity.Category(),
				chargeEntity.PaymentId(),
				createdAt,
				updatedAt,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		dao := dao.NewChargeDao(db)
		result, err := dao.Insert(chargeEntity)

		assert.NoError(t, err)
		if assert.NotNil(t, result) {
			assert.Equal(t, int64(1), result.Id())
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when database operation fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(`INSERT INTO charges`).
			WillReturnError(assert.AnError)

		dao := dao.NewChargeDao(db)
		result, err := dao.Insert(chargeEntity)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when failing to get last insert ID", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(`INSERT INTO charges`).
			WillReturnResult(sqlmock.NewErrorResult(assert.AnError))

		dao := dao.NewChargeDao(db)
		result, err := dao.Insert(chargeEntity)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestChargeDao_FindById(t *testing.T) {
	t.Run("should find payment by ID successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		now := time.Now()
		expectedID := int64(1)
		rows := sqlmock.NewRows([]string{"id", "amount", "category", "payment_id", "created_at", "updated_at"}).
			AddRow(expectedID, 100.5, "finance_fee", 123, now, now)

		mock.ExpectQuery(`SELECT id, amount, category, payment_id, created_at, updated_at FROM charges WHERE id = ?`).
			WithArgs(expectedID).
			WillReturnRows(rows)

		dao := dao.NewChargeDao(db)
		result, err := dao.FindById(expectedID)

		assert.NoError(t, err)
		if assert.NotNil(t, result) {
			assert.Equal(t, expectedID, result.Id())
			assert.Equal(t, int64(123), result.PaymentId())
			assert.Equal(t, "finance_fee", result.Category())
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
		mock.ExpectQuery(`SELECT id, amount, category, payment_id, created_at, updated_at FROM charges WHERE id = ?`).
			WithArgs(expectedID).
			WillReturnError(assert.AnError)

		dao := dao.NewChargeDao(db)
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
		rows := sqlmock.NewRows([]string{"id", "order_id", "status", "payment_type", "created_at", "updated_at", "details", "amount"})

		mock.ExpectQuery(`SELECT id, amount, category, payment_id, created_at, updated_at FROM charges WHERE id = ?`).
			WithArgs(expectedID).
			WillReturnRows(rows)

		dao := dao.NewChargeDao(db)
		result, err := dao.FindById(expectedID)

		assert.NoError(t, err)
		if assert.NotNil(t, result) {
			assert.Equal(t, int64(0), result.Id())
			assert.Equal(t, int64(0), result.PaymentId())
			assert.Equal(t, "", result.Category())
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

		mock.ExpectQuery(`SELECT id, amount, category, payment_id, created_at, updated_at FROM charges WHERE id = ?`).
			WithArgs(expectedID).
			WillReturnRows(rows)

		dao := dao.NewChargeDao(db)
		result, err := dao.FindById(expectedID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
