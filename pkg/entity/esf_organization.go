package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EstOrganization struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string
	Description string
	Token       string
	DBName      string `gorm:"column:db_name"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

// Validate проверяет валидность данных организации
func (o *EstOrganization) Validate() error {
	if o.Name == "" {
		return fmt.Errorf("organization name cannot be empty")
	}
	if len(o.Name) < 2 {
		return fmt.Errorf("organization name must be at least 2 characters long")
	}
	if len(o.Name) > 255 {
		return fmt.Errorf("organization name must be at most 255 characters long")
	}

	if o.Token == "" {
		return fmt.Errorf("organization token cannot be empty")
	}
	if len(o.Token) < 10 {
		return fmt.Errorf("organization token must be at least 10 characters long")
	}

	if o.DBName == "" {
		return fmt.Errorf("database name cannot be empty")
	}
	if !isValidDBName(o.DBName) {
		return fmt.Errorf("invalid database name format")
	}

	return nil
}

// isValidDBName проверяет корректность имени БД
// Должно соответствовать PostgreSQL правилам имен
func isValidDBName(dbName string) bool {
	if len(dbName) == 0 || len(dbName) > 63 {
		return false
	}
	// PostgreSQL требует, чтобы имена начинались с буквы или underscore
	// и содержали только буквы, цифры и underscores
	for i, ch := range dbName {
		if i == 0 {
			if !(ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_') {
				return false
			}
		} else {
			if !(ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9' || ch == '_') {
				return false
			}
		}
	}
	return true
}
