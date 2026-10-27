package domain

import "github.com/shopspring/decimal"

// AccountType is the accounting classification of an account.
type AccountType string

const (
	AccountTypeAsset     AccountType = "asset"
	AccountTypeLiability AccountType = "liability"
	AccountTypeEquity    AccountType = "equity"
	AccountTypeRevenue   AccountType = "revenue"
	AccountTypeExpense   AccountType = "expense"
)

type Direction string

const (
	DirectionDebit  Direction = "debit"
	DirectionCredit Direction = "credit"
)

// accountStatus is the lifecycle state of an account
type AccountStatus string

const (
	AccountStatusActive AccountStatus = "active"
	AccountStatusFrozen AccountStatus = "frozen"
	AccountStatusClosed AccountStatus = "closed"
)

// TransactionStatus is the lifecycle state of a transaction

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusPosted  TransactionStatus = "posted"
	TransactionStatusVoided  TransactionStatus = "voided"
)

type MonetaryValue struct {
	Amount   decimal.Decimal
	Currency string
}
