package domain

import (
	"errors"
	"slices"
	"time"

	"github.com/47monad/xuid"
)

var (
	ErrEmptyEntryID            = errors.New("entry id is required")
	ErrEmptyEntryTransactionID = errors.New("entry transaction id is requried")
	ErrEmptyEntryLedgerID      = errors.New("entry ledger id is required")
	ErrEmptyEntryAccountID     = errors.New("entry account id is necessary")
	ErrEmptyEntryDirection     = errors.New("entry direction is required")
	ErrInvalidEntryAmount      = errors.New("entry amount must be a positive number")
	ErrEmptyEntryCurrency      = errors.New("entry currency is required")
	ErrUnsuppEntryDirection    = errors.New("provided entry direction isn't supported")
)

type EntryFactory struct{}

func NewEntryFactory() EntryFactory {
	return EntryFactory{}
}

type EntryOption func(*Entry)

func (e EntryFactory) Make(id, transactionID, ledgerID, AccountID xuid.XUID, direction Direction, money MonetaryValue, opts ...EntryOption) (*Entry, error) {
	if xuid.IsEmpty(id) {
		return nil, ErrEmptyEntryID
	}

	if xuid.IsEmpty(transactionID) {
		return nil, ErrEmptyEntryTransactionID
	}

	if xuid.IsEmpty(ledgerID) {
		return nil, ErrEmptyEntryLedgerID
	}

	if xuid.IsEmpty(AccountID) {
		return nil, ErrEmptyEntryAccountID
	}

	if direction == "" {
		return nil, ErrEmptyEntryDirection
	}

	if !slices.Contains([]Direction{
		DirectionCredit,
		DirectionDebit,
	}, direction) {
		return nil, ErrUnsuppEntryDirection
	}

	if money.Amount.IsZero() || money.Amount.IsNegative() {
		return nil, ErrInvalidEntryAmount
	}

	// TODO: a more sophisticated currency validations must be implemented
	if money.Currency == "" {
		return nil, ErrEmptyEntryCurrency
	}

	now := time.Now().UTC()
	entry := &Entry{
		ID:            id,
		TransactionID: transactionID,
		LedgerID:      ledgerID,
		AccountID:     AccountID,
		Direction:     direction,
		Line:          1,
		MonetaryValue: money,
		EffectiveAt:   now,
		CreatedAt:     now,
	}

	for _, opt := range opts {
		opt(entry)
	}

	return entry, nil
}

func WithEntryLine(line int32) EntryOption {
	return func(e *Entry) {
		e.Line = line
	}
}
func WithEntryEffectiveAt(ea time.Time) EntryOption {
	return func(e *Entry) {
		e.EffectiveAt = ea
	}
}

func WithEntryDescription(description string) EntryOption {
	return func(e *Entry) {
		e.Description = description
	}
}

func WithEntryMetadata(md map[string]string) EntryOption {
	return func(e *Entry) {
		e.Metadata = md
	}
}

func WithEntryCreatedAt(ca time.Time) EntryOption {
	return func(e *Entry) {
		e.CreatedAt = ca
	}
}
