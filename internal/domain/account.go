package domain

import (
	"time"

	"github.com/47monad/xuid"
)

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

const AccountIDPrefix string = "account_"
