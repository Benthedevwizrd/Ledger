package domain

import (
	"time"

	"github.com/47monad/xuid"
)

// ---------------------------------------------------------------------------
// Ledger
// ---------------------------------------------------------------------------

// Ledger is a self-contained book of accounts. Every account and transaction
// belongs to exactly one ledger.

type Ledger struct {
	ID          xuid.XUID
	Name        string
	Description string

	Metadata map[string]string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// ---------------------------------------------------------------------------
// Account
// ---------------------------------------------------------------------------

// Account is a node in the chart of accounts. It holds a balance per currency
// (see Balance) but never stores the balance itself.

type Account struct {
	ID       xuid.XUID
	LedgerID xuid.XUID

	Code string // human-facing code, unique within the ledger e.g. "1010"
	Name string
	Type AccountType

	// NormalBalance is the side that increases this account.
	// Asset/Expense -> debit, Liability/Equity/Revenue -> credit.
	NormalBalance Direction

	Status   AccountStatus
	Metadata map[string]string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// ---------------------------------------------------------------------------
// Transaction
// ---------------------------------------------------------------------------

// Transaction is a balanced, atomic group of entries (a journal entry).
// Once posted it is immutable; corrections are made by a reversing transaction.

type Transaction struct {
	ID       xuid.XUID
	LedgerID xuid.XUID

	// IdempotencyKey is supplied by the caller and unique within the ledger.
	IdempotencyKey string
	Reference      string // external reference, e.g. an order or payment ID
	Description    string

	Status TransactionStatus

	// EffectiveAt is the business date the transaction applies to.
	// PostedAt is the system time it became final.
	EffectiveAt time.Time
	PostedAt    *time.Time

	// ReversesID points to the transaction this one reverses, if any.
	ReversesID *xuid.XUID

	Metadata map[string]string

	Entries []Entry // loaded with the transaction; persisted in their own table

	CreatedAt time.Time
	UpdatedAt time.Time
}

// ---------------------------------------------------------------------------
// Entry
// ---------------------------------------------------------------------------

// Entry is a single debit or credit line against one account.
// A transaction's entries must balance per currency (invariant, not enforced here).
type Entry struct {
	ID            xuid.XUID
	TransactionID xuid.XUID
	LedgerID      xuid.XUID // denormalised for ledger-scoped queries
	AccountID     xuid.XUID

	Line      int32     // position within the transaction, starting at 1
	Direction Direction // debit or credit

	monetaryValue MonetaryValue

	// EffectiveAt is copied from the transaction so statements can be
	// queried without a join.
	EffectiveAt time.Time

	Description string
	Metadata    map[string]string

	CreatedAt time.Time // entries are append-only: no UpdatedAt
}

// ---------------------------------------------------------------------------
// Balance
// ---------------------------------------------------------------------------

// Balance is the running balance of one account in one currency.
// It is derived from entries and kept as a materialised row for fast reads.
type Balance struct {
	AccountID xuid.XUID
	LedgerID  xuid.XUID

	DebitsPosted   MonetaryValue // sum of posted debit entries
	CreditsPosted  MonetaryValue // sum of posted credit entries
	DebitsPending  MonetaryValue // sum of pending debit entries
	CreditsPending MonetaryValue // sum of pending credit entries

	LastEntryID *xuid.XUID // last entry applied; helps reconciliation
	Version     int64      // optimistic concurrency counter

	UpdatedAt time.Time
}
