package transaction

import (
	"context"

	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Transaction представляет транзакцию базы данных
type Transaction struct {
	tx     *gorm.DB
	logger *logrus.Logger
	ctx    context.Context
}

// Begin начинает новую транзакцию
func Begin(ctx context.Context, db *gorm.DB, logger *logrus.Logger) *Transaction {
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		logger.WithError(tx.Error).Error("Failed to begin transaction")
		return &Transaction{
			tx:     nil,
			logger: logger,
			ctx:    ctx,
		}
	}

	logger.Debug("Transaction started")

	return &Transaction{
		tx:     tx,
		logger: logger,
		ctx:    ctx,
	}
}

// GetDB возвращает GORM DB с активной транзакцией
func (t *Transaction) GetDB() *gorm.DB {
	if t.tx == nil {
		return nil
	}
	return t.tx
}

// Commit коммитит транзакцию
func (t *Transaction) Commit() error {
	if t.tx == nil {
		return apperror.New(apperror.ErrDatabase, "transaction not initialized")
	}

	if err := t.tx.Commit().Error; err != nil {
		t.logger.WithError(err).Error("Failed to commit transaction")
		return apperror.DatabaseError("committing transaction", err)
	}

	t.logger.Debug("Transaction committed successfully")
	return nil
}

// Rollback откатывает транзакцию
func (t *Transaction) Rollback() error {
	if t.tx == nil {
		return apperror.New(apperror.ErrDatabase, "transaction not initialized")
	}

	if err := t.tx.Rollback().Error; err != nil {
		t.logger.WithError(err).Error("Failed to rollback transaction")
		return apperror.DatabaseError("rolling back transaction", err)
	}

	t.logger.Debug("Transaction rolled back successfully")
	return nil
}

// Execute выполняет функцию внутри транзакции с автоматическим commit/rollback
func Execute(ctx context.Context, db *gorm.DB, logger *logrus.Logger, fn func(*gorm.DB) error) error {
	tx := Begin(ctx, db, logger)
	if tx.tx == nil {
		return apperror.New(apperror.ErrDatabase, "failed to initialize transaction")
	}

	// Выполняем функцию
	if err := fn(tx.tx); err != nil {
		// Откатываем при ошибке
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.WithError(rollbackErr).Error("Failed to rollback transaction after error")
		}
		return err
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return err
	}

	return nil
}

// ExecuteWithSavepoint выполняет функцию с savepoint для более гибкого контроля
func ExecuteWithSavepoint(ctx context.Context, db *gorm.DB, logger *logrus.Logger, fn func(*gorm.DB) error) error {
	tx := Begin(ctx, db, logger)
	if tx.tx == nil {
		return apperror.New(apperror.ErrDatabase, "failed to initialize transaction")
	}

	// Создаем savepoint для возможности частичного отката
	savepointTx := tx.tx.SavePoint("sp1")
	if savepointTx.Error != nil {
		tx.Rollback()
		return apperror.DatabaseError("creating savepoint", savepointTx.Error)
	}

	// Выполняем функцию
	if err := fn(tx.tx); err != nil {
		// Откатываем до savepoint
		rollbackTx := tx.tx.RollbackTo("sp1")
		if rollbackTx.Error != nil {
			logger.WithError(rollbackTx.Error).Error("Failed to rollback to savepoint")
		}

		// Коммитим откат
		if commitErr := tx.Commit(); commitErr != nil {
			logger.WithError(commitErr).Error("Failed to commit after savepoint rollback")
		}

		return err
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return err
	}

	return nil
}
