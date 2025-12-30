package models

// EsfEntriesModel представляет модель записи в электронной счет-фактуре (ЭСФ).
// Содержит информацию о товаре или услуге, включая коды классификации,
// количественные и стоимостные показатели, а также данные о налогах.
type EsfEntriesModel struct {
	// ID - уникальный идентификатор записи в базе данных
	ID int `json:"id" valid:"required"`

	// UnitClassificationCode - код единицы измерения по классификатору
	// (например: шт, кг, л и т.д.)
	UnitClassificationCode string `json:"unitClassificationCode" valid:"required"`

	// SalesTaxCode - код товара или услуги по классификатору
	// для целей налогообложения
	SalesTaxCode string `json:"salesTaxCode" valid:"required"`
	// CustomsAuthorityCode - код таможенного органа, через который
	// прошло оформление товара (если применимо)
	CustomsAuthorityCode string `json:"customsAuthorityCode" valid:"required"`

	// Quantity - количество товара или объем услуги
	// в указанных единицах измерения
	Quantity float64 `json:"quantity" valid:"required"`

	// Price - цена за единицу товара или услуги
	// без учета налогов
	Price float64 `json:"price" valid:"required"`

	// VatAmount - сумма налога на добавленную стоимость (НДС)
	// для данной позиции
	VatAmount float64 `json:"vatAmount" valid:"required"`

	// SalesTaxAmount - сумма акцизного налога
	// для подакцизных товаров
	SalesTaxAmount float64 `json:"salesTaxAmount" valid:"required"`

	// AmountWithoutTaxes - общая сумма за позицию
	// без учета НДС и акцизов
	AmountWithoutTaxes float64 `json:"amountWithoutTaxes" valid:"required"`

	// TotalAmount - итоговая сумма за позицию
	// с учетом всех налогов
	TotalAmount float64 `json:"totalAmount" valid:"required"`
}

// CatalogEntriesModels представляет список товаров и услуг
// в электронной счет-фактуре (ЭСФ).
// Используется для передачи массива позиций в документе.
type CatalogEntriesModels []EsfEntriesModel
