package domain

import (
	"errors"
	"time"

	"github.com/47monad/xuid"
)

type AccountOption func(*Account)

type AccountFactory struct{}

func NewAccountFactory() AccountFactory {
	return AccountFactory{}
}

var (
	ErrEmptyAccountID       = errors.New("account id is required")
	ErrEmptyAccountCode     = errors.New("account code is required")
	ErrEmptyAccountLedgerID = errors.New("ledger id is required")
	ErrEmptyAccountName     = errors.New("account name is required")
)

func (a AccountFactory) Make(id, ledgerID xuid.XUID, code, name string, opts ...AccountOption) (*Account, error) {
	if xuid.IsEmpty(id) {
		return nil, ErrEmptyAccountID
	}

	if xuid.IsEmpty(ledgerID) {
		return nil, ErrEmptyAccountLedgerID
	}

	if code == "" {
		return nil, ErrEmptyAccountCode
	}

	if name == "" {
		return nil, ErrEmptyAccountName
	}

	account := &Account{
		ID:            id,
		LedgerID:      ledgerID,
		Code:          code,
		Type:          AccountTypeAsset,
		Name:          name,
		NormalBalance: DirectionCredit,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	for _, opt := range opts {
		opt(account)
	}

	return account, nil
}

func WithAccountType(accountType AccountType) AccountOption {
	return func(a *Account) {
		a.Type = accountType
	}
}

func WithAccountNormalBalance(dir Direction) AccountOption {
	return func(a *Account) {
		a.NormalBalance = dir
	}
}

func WithAccountCreatedAt(ca time.Time) AccountOption {
	return func(a *Account) {
		a.CreatedAt = ca
	}
}

func WithAccountUpdatedAt(ua time.Time) AccountOption {
	return func(a *Account) {
		a.UpdatedAt = ua
	}
}
