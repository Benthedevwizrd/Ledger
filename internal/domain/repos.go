package domain

import (
	"context"

	"github.com/47monad/xuid"
)

type LedgerRepo interface {
	Insert(ctx context.Context, ledger *Ledger) error
	GetByID(ctx context.Context, id xuid.XUID) (*Ledger, error)
}

type AccountRepo interface {
	Insert(ctx context.Context, account *Account) error
	GetByID(ctx context.Context, id xuid.XUID) (*Account, error)
}

type EntryRepo interface {
	Insert(ctx context.Context, entry *Entry) error
	GetByID(ctx context.Context, id xuid.XUID) (*Entry, error)
}

type TransactionRepo interface {
	Insert(ctx context.Context, transaction *Transaction) error
	GetByID(ctx context.Context, id xuid.XUID) (*Transaction, error)
}

type BalanceRepo interface {
	Insert(ctx context.Context, balance *Balance) error
	GetByAccountID(ctx context.Context, accountID xuid.XUID) (*Balance, error)
	UpdateByAccountID(ctx context.Context, balance *Balance) (*Balance, error)
}
