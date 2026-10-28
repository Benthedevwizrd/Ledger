package domain

import (
	"errors"
	"time"

	"github.com/47monad/xuid"
)

var (
	ErrEmptyTransactonID              = errors.New("transaction id is required")
	ErrEmptyTransactionLedgerID       = errors.New("transaction ledger id is required")
	ErrEmptyTransactionIdempotencyKey = errors.New("idempotency key is required")
	ErrEmptyTransactionReferece       = errors.New("reference is required")
)

type TransactionFactory struct{}

func NewTransactionFactory() TransactionFactory {
	return TransactionFactory{}
}

type TransactionOption func(*Transaction)

func (t TransactionFactory) Make(id, ledgerID xuid.XUID, idempotencyKey, reference string, opts ...TransactionOption) (*Transaction, error) {
	if xuid.IsEmpty(id) {
		return nil, ErrEmptyTransactonID
	}

	if xuid.IsEmpty(ledgerID) {
		return nil, ErrEmptyTransactionLedgerID
	}

	if idempotencyKey == "" {
		return nil, ErrEmptyTransactionIdempotencyKey
	}

	if reference == "" {
		return nil, ErrEmptyTransactionReferece
	}

	now := time.Now().UTC()

	transaction := &Transaction{
		ID:             id,
		LedgerID:       ledgerID,
		IdempotencyKey: idempotencyKey,
		Reference:      reference,
		EffectiveAt:    now,
		Status:         TransactionStatusPending,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	for _, opt := range opts {
		opt(transaction)
	}

	return transaction, nil
}

func WithTransactionDescription(description string) TransactionOption {
	return func(t *Transaction) {
		t.Description = description
	}
}

func WithTransactionStatus(status TransactionStatus) TransactionOption {
	return func(t *Transaction) {
		t.Status = status
	}
}

func WithTransactionEffectiveAt(ea time.Time) TransactionOption {
	return func(t *Transaction) {
		t.EffectiveAt = ea
	}
}

func WithTransactionPostedAt(pa time.Time) TransactionOption {
	return func(t *Transaction) {
		t.PostedAt = &pa
	}
}

func WithTransactionReversesID(id xuid.XUID) TransactionOption {
	return func(t *Transaction) {
		t.ReversesID = &id
	}
}

func WithTransactionMetadata(md map[string]string) TransactionOption {
	return func(t *Transaction) {
		t.Metadata = md
	}
}

func WithTransactionEntries(entries []Entry) TransactionOption {
	return func(t *Transaction) {
		t.Entries = entries
	}
}

func WithTransactionCreatedAt(ca time.Time) TransactionOption {
	return func(t *Transaction) {
		t.CreatedAt = ca
	}
}

func WithTransactionUpdatedAt(ua time.Time) TransactionOption {
	return func(t *Transaction) {
		t.UpdatedAt = ua
	}
}
