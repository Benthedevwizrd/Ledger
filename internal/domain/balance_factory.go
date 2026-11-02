package domain

import (
	"errors"
	"time"

	"github.com/47monad/xuid"
)

var (
	ErrEmptyBalanceAccountID = errors.New("account id is required")
	ErrEmptyBalanceLedgerID  = errors.New("ledger id is required")
)

type BalanceFactory struct{}

func NewBalanceFactory() BalanceFactory {
	return BalanceFactory{}
}

type BalanceOption func(*Balance)

func (b BalanceFactory) Make(accountID, ledgerID xuid.XUID, opts ...BalanceOption) (*Balance, error) {
	if xuid.IsEmpty(accountID) {
		return nil, ErrEmptyBalanceAccountID
	}
	if xuid.IsEmpty(ledgerID) {
		return nil, ErrEmptyBalanceLedgerID
	}

	balance := &Balance{
		AccountID: accountID,
		LedgerID:  ledgerID,
	}

	for _, opt := range opts {
		opt(balance)
	}
	return balance, nil
}

func WithBalanceDebitsPosted(dp MonetaryValue) BalanceOption {
	return func(b *Balance) {
		b.DebitsPosted = dp
	}
}

func WithBalanceDebitsPending(dp MonetaryValue) BalanceOption {
	return func(b *Balance) {
		b.DebitsPending = dp
	}
}

func WithBalanceCreditsPosted(cp MonetaryValue) BalanceOption {
	return func(b *Balance) {
		b.CreditsPosted = cp
	}
}

func WithBalanceCreditsPending(cp MonetaryValue) BalanceOption {
	return func(b *Balance) {
		b.CreditsPending = cp
	}
}

func WithBalanceLastEntryID(id xuid.XUID) BalanceOption {
	return func(b *Balance) {
		b.LastEntryID = &id
	}
}

func WithBalanceUpdatedAt(ua time.Time) BalanceOption {
	return func(b *Balance) {
		b.UpdatedAt = ua
	}
}

func WithBalanceVersion(version int64) BalanceOption {
	return func(b *Balance) {
		b.Version = version
	}
}
