package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Create document

type EsfDocument struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// false Наименование иностранца или Наименование на иностранном языке
	ForeignName string `gorm:"size:255" json:"foreignName"`
	// true Отправить от имени филиала
	IsBranchDataSent bool `gorm:"not null" json:"isBranchDataSent" valid:"required"`
	// true Цена без налогов
	IsPriceWithoutTaxes bool `gorm:"not null" json:"isPriceWithoutTaxes" valid:"required"`
	// false ИНН филиала
	AffiliateTin string `gorm:"size:14" json:"affiliateTin"`
	// false Отраслевые
	IsIndustry bool `gorm:"default:false" json:"isIndustry"`
	// false Номер учетной системы
	OwnedCrmReceiptCode string `gorm:"size:100" json:"ownedCrmReceiptCode"`
	// true Код вида операции
	OperationTypeCode string `gorm:"size:20;not null" json:"operationTypeCode" valid:"required"`
	// true Дата поставки
	DeliveryDate time.Time `gorm:"not null" json:"deliveryDate" valid:"required"`
	// true Код типа поставки
	DeliveryTypeCode string `gorm:"size:20;not null" json:"deliveryTypeCode" valid:"required"`
	// true Субъект Кыргызской Республики
	IsResident bool `gorm:"not null" json:"isResident" valid:"required"`
	// true ИНН покупателя
	ContractorTin string `gorm:"size:14;not null" json:"contractorTin" valid:"required"`
	// false Номер банковского счета поставщика
	SupplierBankAccount string `gorm:"size:50" json:"supplierBankAccount"`
	// false Номер банковского счета покупателя
	ContractorBankAccount string `gorm:"size:50" json:"contractorBankAccount"`
	// true Код валюты
	CurrencyCode string `gorm:"size:3;not null" json:"currencyCode" valid:"required"`
	// false Код страны
	CountryCode string `gorm:"size:2" json:"countryCode"`
	// false Курс валюты к сому
	CurrencyRate float64 `gorm:"type:decimal(10,4)" json:"currencyRate"`
	// false Общая стоимость в валюте
	TotalCurrencyValue float64 `gorm:"type:decimal(15,2)" json:"totalCurrencyValue"`
	// false Общая стоимость в валюте без налогов
	TotalCurrencyValueWithoutTaxes float64 `gorm:"type:decimal(15,2)" json:"totalCurrencyValueWithoutTaxes"`
	// false Номер договора на поставку
	SupplyContractNumber string `gorm:"size:100" json:"supplyContractNumber"`
	// false Дата договора на поставку
	ContractStartDate time.Time `json:"contractStartDate"`
	// false Дата окончания договора на поставку Комментарий
	Comment string `gorm:"type:text" json:"comment"`
	// false Код способа доставки
	DeliveryCode string `gorm:"size:20" json:"deliveryCode"`
	// true Код формы оплаты
	PaymentCode string `gorm:"size:20;not null" json:"paymentCode" valid:"required"`
	// true Код ставки НДС
	TaxRateVATCode string `gorm:"size:20;not null" json:"taxRateVATCode" valid:"required"`
	// true Товары и услуги
	CatalogEntries []EsfEntries `gorm:"foreignKey:DocumentID;constraint:OnDelete:CASCADE" json:"catalogEntries"`
	// false Начальные остатки,сальдо на начало периода
	OpeningBalances float64 `gorm:"type:decimal(15,2);default:0" json:"openingBalances"`
	// false Начисленные взносы
	AssessedContributionsAmount float64 `gorm:"type:decimal(15,2);default:0" json:"assessedContributionsAmount"`
	// false Поступления,оплачено
	PaidAmount float64 `gorm:"type:decimal(15,2);default:0" json:"paidAmount"`
	// false Штрафы
	PenaltiesAmount float64 `gorm:"type:decimal(15,2);default:0" json:"penaltiesAmount"`
	// false пени
	FinesAmount float64 `gorm:"type:decimal(15,2);default:0" json:"finesAmount"`
	// false Конечные остатки,сальдо на конец периода
	ClosingBalances float64 `gorm:"type:decimal(15,2);default:0" json:"closingBalances"`
	// false Сумма к оплате
	AmountToBePaid float64 `gorm:"type:decimal(15,2);default:0" json:"amountToBePaid"`
	// false Лицевой счет
	PersonalAccountNumber string `gorm:"size:50" json:"personalAccountNumber"`
}

func (EsfDocument) TableName() string {
	return "esf_documents"
}
