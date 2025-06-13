package dao_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"payment-gateway/cmd/domain/payment"
	"payment-gateway/cmd/infra/dao"
)

func TestPaymentDao_Insert(t *testing.T) {
	paymentEntity := payment.NewPayment(123, 100.5, "credit_card")

	t.Run("should insert payment successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		createdAt := paymentEntity.CreatedAt().Format("2006-01-02 15:04:05")
		updatedAt := paymentEntity.UpdatedAt().Format("2006-01-02 15:04:05")

		mock.ExpectExec(`INSERT INTO payments`).
			WithArgs(
				paymentEntity.OrderID(),
				paymentEntity.Status(),
				paymentEntity.Type(),
				paymentEntity.Amount(),
				createdAt,
				updatedAt,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		dao := dao.NewPaymentDao(db)
		result, err := dao.Insert(paymentEntity)

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

		mock.ExpectExec(`INSERT INTO payments`).
			WillReturnError(assert.AnError)

		dao := dao.NewPaymentDao(db)
		result, err := dao.Insert(paymentEntity)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when failing to get last insert ID", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(`INSERT INTO payments`).
			WillReturnResult(sqlmock.NewErrorResult(assert.AnError))

		dao := dao.NewPaymentDao(db)
		result, err := dao.Insert(paymentEntity)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPaymentDao_FindById(t *testing.T) {
	t.Run("should find payment by ID successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		now := time.Now()
		expectedID := int64(1)
		rows := sqlmock.NewRows([]string{"id", "order_id", "status", "payment_type", "created_at", "updated_at", "details", "amount"}).
			AddRow(expectedID, 123, "approved", "credit_card", now, now, "test details", 100.5)

		mock.ExpectQuery(`SELECT id, order_id, status, payment_type, created_at, updated_at, IFNULL\(details, ''\) as details, amount FROM payments WHERE id = \?`).
			WithArgs(expectedID).
			WillReturnRows(rows)

		dao := dao.NewPaymentDao(db)
		result, err := dao.FindById(expectedID)

		assert.NoError(t, err)
		if assert.NotNil(t, result) {
			assert.Equal(t, expectedID, result.Id())
			assert.Equal(t, int64(123), result.OrderID())
			assert.Equal(t, "approved", result.Status())
			assert.Equal(t, "credit_card", result.Type())
			assert.Equal(t, "test details", result.Details())
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
		mock.ExpectQuery(`SELECT id, order_id, status, payment_type, created_at, updated_at, IFNULL\(details, ''\) as details, amount FROM payments WHERE id = \?`).
			WithArgs(expectedID).
			WillReturnError(assert.AnError)

		dao := dao.NewPaymentDao(db)
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

		mock.ExpectQuery(`SELECT id, order_id, status, payment_type, created_at, updated_at, IFNULL\(details, ''\) as details, amount FROM payments WHERE id = \?`).
			WithArgs(expectedID).
			WillReturnRows(rows)

		dao := dao.NewPaymentDao(db)
		result, err := dao.FindById(expectedID)

		assert.NoError(t, err)
		if assert.NotNil(t, result) {
			assert.Equal(t, int64(0), result.Id())
			assert.Equal(t, int64(0), result.OrderID())
			assert.Equal(t, "", result.Status())
			assert.Equal(t, "", result.Type())
			assert.Equal(t, "", result.Details())
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

		mock.ExpectQuery(`SELECT id, order_id, status, payment_type, created_at, updated_at, IFNULL\(details, ''\) as details, amount FROM payments WHERE id = \?`).
			WithArgs(expectedID).
			WillReturnRows(rows)

		dao := dao.NewPaymentDao(db)
		result, err := dao.FindById(expectedID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPaymentDao_FindByOrderId(t *testing.T) {
	t.Run("should find payments by order ID successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		now := time.Now()
		orderID := int64(123)
		rows := sqlmock.NewRows([]string{"id", "order_id", "status", "payment_type", "created_at", "updated_at", "details", "amount"}).
			AddRow(1, orderID, "approved", "credit_card", now, now, "test details 1", 100.5).
			AddRow(2, orderID, "pending", "pix", now, now, "test details 2", 200.0)

		mock.ExpectQuery(`SELECT id, order_id, status, payment_type, created_at, updated_at, IFNULL\(details, ''\) as details, amount FROM payments WHERE order_id = \?`).
			WithArgs(orderID).
			WillReturnRows(rows)

		paymentDao := dao.NewPaymentDao(db)
		result, err := paymentDao.FindByOrderId(orderID)

		assert.NoError(t, err)
		if assert.Len(t, result, 2) {
			assert.Equal(t, int64(1), result[0].Id())
			assert.Equal(t, orderID, result[0].OrderID())
			assert.Equal(t, "approved", result[0].Status())
			assert.Equal(t, "credit_card", result[0].Type())
			assert.Equal(t, "test details 1", result[0].Details())
			assert.Equal(t, 100.5, result[0].Amount())

			assert.Equal(t, int64(2), result[1].Id())
			assert.Equal(t, orderID, result[1].OrderID())
			assert.Equal(t, "pending", result[1].Status())
			assert.Equal(t, "pix", result[1].Type())
			assert.Equal(t, "test details 2", result[1].Details())
			assert.Equal(t, 200.0, result[1].Amount())
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return empty slice when no payments found for order", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		orderID := int64(999)
		rows := sqlmock.NewRows([]string{"id", "order_id", "status", "payment_type", "created_at", "updated_at", "details", "amount"})

		mock.ExpectQuery(`SELECT id, order_id, status, payment_type, created_at, updated_at, IFNULL\(details, ''\) as details, amount FROM payments WHERE order_id = \?`).
			WithArgs(orderID).
			WillReturnRows(rows)

		paymentDao := dao.NewPaymentDao(db)
		result, err := paymentDao.FindByOrderId(orderID)

		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when query fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		orderID := int64(123)

		mock.ExpectQuery(`SELECT id, order_id, status, payment_type, created_at, updated_at, IFNULL\(details, ''\) as details, amount FROM payments WHERE order_id = \?`).
			WithArgs(orderID).
			WillReturnError(assert.AnError)

		paymentDao := dao.NewPaymentDao(db)
		result, err := paymentDao.FindByOrderId(orderID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when scan fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		orderID := int64(123)
		rows := sqlmock.NewRows([]string{"id", "order_id"}).AddRow(1, orderID)

		mock.ExpectQuery(`SELECT id, order_id, status, payment_type, created_at, updated_at, IFNULL\(details, ''\) as details, amount FROM payments WHERE order_id = \?`).
			WithArgs(orderID).
			WillReturnRows(rows)

		paymentDao := dao.NewPaymentDao(db)
		result, err := paymentDao.FindByOrderId(orderID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPaymentDao_Update(t *testing.T) {
	t.Run("should update payment successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		paymentEntity := payment.NewPayment(123, 100.5, "credit_card")
		paymentEntity.SetId(1)
		paymentEntity.SetStatus("approved")

		mock.ExpectExec(`UPDATE payments`).
			WithArgs(
				paymentEntity.Status(),
				sqlmock.AnyArg(), // updated_at
				paymentEntity.Id(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		paymentDao := dao.NewPaymentDao(db)
		result, err := paymentDao.Update(paymentEntity)

		assert.NoError(t, err)
		if assert.NotNil(t, result) {
			assert.Equal(t, paymentEntity.Id(), result.Id())
			assert.Equal(t, paymentEntity.Status(), result.Status())
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		paymentEntity := payment.NewPayment(123, 100.5, "credit_card")
		paymentEntity.SetId(1)
		paymentEntity.SetStatus("approved")

		mock.ExpectExec(`UPDATE payments`).
			WillReturnError(assert.AnError)

		paymentDao := dao.NewPaymentDao(db)
		result, err := paymentDao.Update(paymentEntity)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
