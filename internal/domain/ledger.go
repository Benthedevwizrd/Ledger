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

const LedgerIDPrefix string = "ledger_"
