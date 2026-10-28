package domain

import (
	"errors"
	"time"

	"github.com/47monad/xuid"
	"github.com/shopspring/decimal"
)

var (
	ErrEmptyEntryID            = errors.New("entry id is required")
	ErrEmptyEntryTransactionID = errors.New("entry transaction id is requried")
	ErrEmptyEntryLedgerID      = errors.New("entry ledger id is required")
	ErrEmptyEntryAccountID     = errors.New("entry account id is necessary")
)

type entryFactory struct{}

func NewEntryFactory() entryFactory {
	return entryFactory{}
}

type EntryOption func(*Entry)

func (e entryFactory) Make(id, transactionID, ledgerID, AccountId xuid.XUID, opts ...EntryOption) (*Entry, error) {
	if xuid.IsEmpty(id) {
		return nil, ErrEmptyEntryID
	}

	if xuid.IsEmpty(transactionID) {
		return nil, ErrEmptyEntryTransactionID
	}

	if xuid.IsEmpty(ledgerID) {
		return nil, ErrEmptyEntryLedgerID
	}

	if xuid.IsEmpty(AccountId) {
		return nil, ErrEmptyEntryAccountID
	}

	now := time.Now().UTC()
	entry := &Entry{
		ID:            id,
		TransactionID: transactionID,
		LedgerID:      ledgerID,
		AccountID:     AccountId,
		Direction:     DirectionCredit,
		MonetaryValue: MonetaryValue{
			Amount:   decimal.Zero,
			Currency: "",
		},
		EffectiveAt: now,
		CreatedAt:   now,
	}

	for _, opt := range opts {
		opt(entry)
	}

	return entry, nil
}
