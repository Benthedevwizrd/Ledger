package domain

import (
	"time"

	"github.com/47monad/xuid"
)

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
