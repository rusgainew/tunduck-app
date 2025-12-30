package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EsfEntries struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	DocumentID uuid.UUID      `gorm:"type:uuid;not null;index" json:"documentId"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// UnitClassificationCode - код единицы измерения по классификатору
	// (например: шт, кг, л и т.д.)
	UnitClassificationCode string `gorm:"size:20;not null" json:"unitClassificationCode" valid:"required"`

	// SalesTaxCode - код товара или услуги по классификатору
	// для целей налогообложения
	SalesTaxCode string `gorm:"size:50;not null" json:"salesTaxCode" valid:"required"`

	// CustomsAuthorityCode - код таможенного органа, через который
	// прошло оформление товара (если применимо)
	CustomsAuthorityCode string `gorm:"size:50" json:"customsAuthorityCode"`

	// Quantity - количество товара или объем услуги
	// в указанных единицах измерения
	Quantity float64 `gorm:"type:decimal(15,4);not null" json:"quantity" valid:"required"`

	// Price - цена за единицу товара или услуги
	// без учета налогов
	Price float64 `gorm:"type:decimal(15,2);not null" json:"price" valid:"required"`

	// VatAmount - сумма налога на добавленную стоимость (НДС)
	// для данной позиции
	VatAmount float64 `gorm:"type:decimal(15,2);default:0" json:"vatAmount"`

	// SalesTaxAmount - сумма акцизного налога
	// для подакцизных товаров
	SalesTaxAmount float64 `gorm:"type:decimal(15,2);default:0" json:"salesTaxAmount"`

	// AmountWithoutTaxes - общая сумма за позицию
	// без учета НДС и акцизов
	AmountWithoutTaxes float64 `gorm:"type:decimal(15,2);not null" json:"amountWithoutTaxes" valid:"required"`

	// TotalAmount - итоговая сумма за позицию
	// с учетом всех налогов
	TotalAmount float64 `gorm:"type:decimal(15,2);not null" json:"totalAmount" valid:"required"`
}

func (EsfEntries) TableName() string {
	return "esf_entries"
}

// CatalogEntriesModels представляет список товаров и услуг
// в электронной счет-фактуре (ЭСФ).
// Используется для передачи массива позиций в документе.
type CatalogEntriesModels []EsfEntries
