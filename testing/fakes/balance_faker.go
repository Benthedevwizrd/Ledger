package fakes

import (
	"time"

	"ledger/internal/domain"

	"github.com/47monad/xuid"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/shopspring/decimal"
)

type BalanceFakerContext struct{}
type BalanceFaker struct{}

func NewBalanceFaker() BalanceFaker {
	return BalanceFaker{}
}

func (b BalanceFaker) make(opts ...domain.BalanceOption) *domain.Balance {
	accountID, err := xuid.NewSortable(domain.AccountIDPrefix)
	if err != nil {
		panic(err)
	}

	ledgerID, err := xuid.NewSortable(domain.LedgerIDPrefix)
	if err != nil {
		panic(err)
	}

	lastEntryID, err := xuid.NewSortable(domain.EntryIDPrefix)
	if err != nil {
		panic(err)
	}

	currency := gofakeit.CurrencyShort()
	now := time.Now().UTC()

	allOpts := append([]domain.BalanceOption{
		domain.WithBalanceDebitsPosted(randomMoney(currency)),
		domain.WithBalanceCreditsPosted(randomMoney(currency)),
		domain.WithBalanceDebitsPending(randomMoney(currency)),
		domain.WithBalanceCreditsPending(randomMoney(currency)),
		domain.WithBalanceLastEntryID(lastEntryID),
		domain.WithBalanceVersion(int64(gofakeit.IntRange(0, 1000))),
		domain.WithBalanceUpdatedAt(now),
	}, opts...)

	made, err := domain.NewBalanceFactory().Make(accountID, ledgerID, allOpts...)
	if err != nil {
		panic(err)
	}

	return made
}

func (b BalanceFaker) Make(opts ...domain.BalanceOption) *domain.Balance {
	return b.make(opts...)
}

func (b BalanceFaker) MakeMany(count int, opts ...domain.BalanceOption) []*domain.Balance {
	if count < 1 {
		count = 1
	}
	balances := make([]*domain.Balance, 0, count)
	for range count {
		made := b.make(opts...)
		balances = append(balances, made)
	}
	return balances
}

func randomMoney(currency string) domain.MonetaryValue {
	return domain.MonetaryValue{
		Amount:   decimal.NewFromInt(int64(gofakeit.IntRange(1, 1_000_000))),
		Currency: currency,
	}
}

func WithBalanceAccountID(id xuid.XUID) domain.BalanceOption {
	return func(b *domain.Balance) {
		b.AccountID = id
	}
}

func WithBalanceLedgerID(id xuid.XUID) domain.BalanceOption {
	return func(b *domain.Balance) {
		b.LedgerID = id
	}
}
