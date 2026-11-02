package domain

import (
	"errors"
	"time"

	"github.com/47monad/xuid"
)

type LedgerOption func(*Ledger)

type LedgerFactory struct{}

func NewLedgerFactory() LedgerFactory {
	return LedgerFactory{}
}

var (
	ErrEmptyLedgerID          = errors.New("ledger id is required")
	ErrEmptyLedgerName        = errors.New("ledger name is required")
	ErrEmptyLedgerDescription = errors.New("ledger description is required")
)

func (l LedgerFactory) Make(id xuid.XUID, name string, description string, opts ...LedgerOption) (*Ledger, error) {
	if xuid.IsEmpty(id) {
		return nil, ErrEmptyLedgerID
	}

	if name == "" {
		return nil, ErrEmptyLedgerName
	}

	if description == "" {
		return nil, ErrEmptyLedgerDescription
	}

	ledger := &Ledger{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	for _, opt := range opts {
		opt(ledger)
	}

	return ledger, nil
}

func WithLedgerMetadata(md map[string]string) LedgerOption {
	return func(l *Ledger) {
		l.Metadata = md
	}
}

func WithLedgerCreatedAt(ca time.Time) LedgerOption {
	return func(l *Ledger) {
		l.CreatedAt = ca
	}
}

func WithLedgerUpdatedAt(ua time.Time) LedgerOption {
	return func(l *Ledger) {
		l.UpdatedAt = ua
	}
}
