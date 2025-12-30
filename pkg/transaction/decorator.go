package transaction

import (
	"context"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// TxDecorator предоставляет методы для работы с транзакциями в repository паттерне
type TxDecorator struct {
	db     *gorm.DB
	logger *logrus.Logger
}

// NewTxDecorator создает новый декоратор транзакций
func NewTxDecorator(db *gorm.DB, logger *logrus.Logger) *TxDecorator {
	return &TxDecorator{
		db:     db,
		logger: logger,
	}
}

// WithTx выполняет операцию внутри транзакции
// Использование:
//
//	txDec := NewTxDecorator(db, logger)
//	err := txDec.WithTx(ctx, func(txDB *gorm.DB) error {
//	    // выполнить операции с txDB
//	    return nil
//	})
func (t *TxDecorator) WithTx(ctx context.Context, fn func(*gorm.DB) error) error {
	return Execute(ctx, t.db, t.logger, fn)
}

// WithTxSavepoint выполняет операцию внутри транзакции с savepoint
func (t *TxDecorator) WithTxSavepoint(ctx context.Context, fn func(*gorm.DB) error) error {
	return ExecuteWithSavepoint(ctx, t.db, t.logger, fn)
}

// GetDB возвращает базовую БД
func (t *TxDecorator) GetDB() *gorm.DB {
	return t.db
}

// GetLogger возвращает логгер
func (t *TxDecorator) GetLogger() *logrus.Logger {
	return t.logger
}
