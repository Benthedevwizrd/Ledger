package domain

import (
	"time"

	"github.com/47monad/xuid"
)

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

const TransactionIDPrefix string = "transaction_"
