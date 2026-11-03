package fakes

import (
	"ledger/internal/domain"

	"github.com/47monad/xuid"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/shopspring/decimal"
)

type EntryFakerContext struct{}
type EntryFaker struct{}

func NewEntryFaker() EntryFaker {
	return EntryFaker{}
}

func (e EntryFaker) make(opts ...domain.EntryOption) *domain.Entry {
	id, err := xuid.NewSortable(domain.EntryIDPrefix)
	if err != nil {
		panic(err)
	}

	transactionID, err := xuid.NewSortable(domain.TransactionIDPrefix)
	if err != nil {
		panic(err)
	}

	ledgerID, err := xuid.NewSortable(domain.LedgerIDPrefix)
	if err != nil {
		panic(err)
	}

	accountID, err := xuid.NewSortable(domain.AccountIDPrefix)
	if err != nil {
		panic(err)
	}

	direction := domain.Direction(gofakeit.RandomString([]string{
		string(domain.DirectionCredit),
		string(domain.DirectionDebit),
	}))

	money := domain.MonetaryValue{
		Amount:   decimal.NewFromInt(int64(gofakeit.IntRange(1, 1_000_000))),
		Currency: gofakeit.CurrencyShort(),
	}

	line := int32(gofakeit.IntRange(1, 10))
	description := gofakeit.Sentence(5)

	numberOfItems := gofakeit.IntRange(1, 10)
	metadata := make(map[string]string, numberOfItems)
	for len(metadata) < numberOfItems {
		metadata[gofakeit.Word()] = gofakeit.Sentence(5)
	}

	allOpts := append([]domain.EntryOption{
		domain.WithEntryLine(line),
		domain.WithEntryDescription(description),
		domain.WithEntryMetadata(metadata),
	}, opts...)

	made, err := domain.NewEntryFactory().Make(id, transactionID, ledgerID, accountID, direction, money, allOpts...)
	if err != nil {
		panic(err)
	}

	return made
}

func (e EntryFaker) Make(opts ...domain.EntryOption) *domain.Entry {
	return e.make(opts...)
}

func (e EntryFaker) MakeMany(count int, opts ...domain.EntryOption) []*domain.Entry {
	if count < 1 {
		count = 1
	}
	entries := make([]*domain.Entry, 0, count)
	for range count {
		made := e.make(opts...)
		entries = append(entries, made)
	}
	return entries
}

func WithEntryID(id xuid.XUID) domain.EntryOption {
	return func(e *domain.Entry) {
		e.ID = id
	}
}

func WithEntryTransactionID(id xuid.XUID) domain.EntryOption {
	return func(e *domain.Entry) {
		e.TransactionID = id
	}
}

func WithEntryLedgerID(id xuid.XUID) domain.EntryOption {
	return func(e *domain.Entry) {
		e.LedgerID = id
	}
}

func WithEntryAccountID(id xuid.XUID) domain.EntryOption {
	return func(e *domain.Entry) {
		e.AccountID = id
	}
}

func WithEntryDirection(direction domain.Direction) domain.EntryOption {
	return func(e *domain.Entry) {
		e.Direction = direction
	}
}

func WithEntryMonetaryValue(money domain.MonetaryValue) domain.EntryOption {
	return func(e *domain.Entry) {
		e.MonetaryValue = money
	}
}
