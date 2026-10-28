package domain

import (
	"time"

	"github.com/47monad/xuid"
)

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
