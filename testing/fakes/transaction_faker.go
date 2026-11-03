package fakes

import (
	"ledger/internal/domain"

	"github.com/47monad/xuid"
	"github.com/brianvoe/gofakeit/v7"
)

type TransactionFakerContext struct{}
type TransactionFaker struct{}

func NewTransactionFaker() TransactionFaker {
	return TransactionFaker{}
}

func (t TransactionFaker) make(opts ...domain.TransactionOption) *domain.Transaction {
	id, err := xuid.NewSortable(domain.TransactionIDPrefix)
	if err != nil {
		panic(err)
	}

	ledgerID, err := xuid.NewSortable(domain.LedgerIDPrefix)
	if err != nil {
		panic(err)
	}

	idempotencyKey := gofakeit.UUID()
	reference := gofakeit.UUID()

	status := domain.TransactionStatus(gofakeit.RandomString([]string{
		string(domain.TransactionStatusPending),
		string(domain.TransactionStatusPosted),
		string(domain.TransactionStatusVoided),
	}))

	description := gofakeit.Sentence(5)

	numberOfItems := gofakeit.IntRange(1, 10)
	metadata := make(map[string]string, numberOfItems)
	for len(metadata) < numberOfItems {
		metadata[gofakeit.Word()] = gofakeit.Sentence(5)
	}

	allOpts := append([]domain.TransactionOption{
		domain.WithTransactionStatus(status),
		domain.WithTransactionDescription(description),
		domain.WithTransactionMetadata(metadata),
	}, opts...)

	made, err := domain.NewTransactionFactory().Make(id, ledgerID, idempotencyKey, reference, allOpts...)
	if err != nil {
		panic(err)
	}

	return made
}

func (t TransactionFaker) Make(opts ...domain.TransactionOption) *domain.Transaction {
	return t.make(opts...)
}

func (t TransactionFaker) MakeMany(count int, opts ...domain.TransactionOption) []*domain.Transaction {
	if count < 1 {
		count = 1
	}
	transactions := make([]*domain.Transaction, 0, count)
	for range count {
		made := t.make(opts...)
		transactions = append(transactions, made)
	}
	return transactions
}

func WithTransactionID(id xuid.XUID) domain.TransactionOption {
	return func(t *domain.Transaction) {
		t.ID = id
	}
}

func WithTransactionLedgerID(id xuid.XUID) domain.TransactionOption {
	return func(t *domain.Transaction) {
		t.LedgerID = id
	}
}

func WithTransactionIdempotencyKey(key string) domain.TransactionOption {
	return func(t *domain.Transaction) {
		t.IdempotencyKey = key
	}
}

func WithTransactionReference(reference string) domain.TransactionOption {
	return func(t *domain.Transaction) {
		t.Reference = reference
	}
}
