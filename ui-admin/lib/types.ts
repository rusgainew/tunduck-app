// Типы данных, соответствующие Go структурам

export type Role = 'admin' | 'user' | 'viewer';

export interface User {
  id: string;
  username: string;
  email: string;
  fullName: string;
  phone: string;
  role: Role;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface EsfOrganization {
  id: string;
  name: string;
  description: string;
  token: string;
  dbName: string;
  createdAt: string;
  updatedAt: string;
  deletedAt?: string;
}

export interface EsfEntry {
  id: string;
  documentID: string;
  catalogCode: string;
  catalogName: string;
  measureCode: string;
  quantity: number;
  price: number;
  turnoverSize: number;
  vatRate: number;
  vatAmount: number;
  exciseRate: number;
  exciseAmount: number;
}

export interface EsfDocument {
  id: string;
  createdAt: string;
  updatedAt: string;
  
  // Основные поля
  foreignName: string;
  isBranchDataSent: boolean;
  isPriceWithoutTaxes: boolean;
  affiliateTin: string;
  isIndustry: boolean;
  ownedCrmReceiptCode: string;
  operationTypeCode: string;
  deliveryDate: string;
  deliveryTypeCode: string;
  isResident: boolean;
  contractorTin: string;
  supplierBankAccount: string;
  contractorBankAccount: string;
  currencyCode: string;
  countryCode: string;
  currencyRate: number;
  totalCurrencyValue: number;
  totalCurrencyValueWithoutTaxes: number;
  supplyContractNumber: string;
  contractStartDate: string;
  comment: string;
  deliveryCode: string;
  paymentCode: string;
  taxRateVATCode: string;
  
  // Дополнительные поля
  catalogEntries: EsfEntry[];
  openingBalances: number;
  assessedContributionsAmount: number;
  paidAmount: number;
  penaltiesAmount: number;
  finesAmount: number;
  closingBalances: number;
  amountToBePaid: number;
  personalAccountNumber: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  fullName: string;
  phone: string;
  password: string;
  confirmPassword: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  totalPages?: number;
}

export interface ApiError {
  error: string;
  message: string;
  details?: string;
  timestamp?: string;
}
