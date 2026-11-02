package fakes

import (
	"ledger/internal/domain"
	"time"

	"github.com/47monad/xuid"
	"github.com/brianvoe/gofakeit/v7"
)

type LedgerFakerContext struct {
}

type LedgerFaker struct{}

func NewLedgerFaker() LedgerFaker {
	return LedgerFaker{}
}

func (l LedgerFaker) make(opts ...domain.LedgerOption) *domain.Ledger {
	id, err := xuid.NewSortable(domain.LedgerIDPrefix)
	if err != nil {
		panic(err)
	}

	name := gofakeit.Company()
	description := gofakeit.Sentence(8)
	now := time.Now().UTC()

	numberOfItems := gofakeit.IntRange(1, 10)
	metadata := make(map[string]string, numberOfItems)
	for len(metadata) < numberOfItems {
		metadata[gofakeit.Word()] = gofakeit.Sentence(5)
	}

	// defaults first, caller opts last so they can override
	allOpts := append([]domain.LedgerOption{
		domain.WithLedgerCreatedAt(now),
		domain.WithLedgerUpdatedAt(now),
		domain.WithLedgerMetadata(metadata),
	}, opts...)

	ledger, err := domain.NewLedgerFactory().Make(id, name, description, allOpts...)
	if err != nil {
		panic(err)
	}
	return ledger
}

func (l LedgerFaker) Make(opts ...domain.LedgerOption) *domain.Ledger {
	return l.make(opts...)
}

func (l LedgerFaker) MakeMany(count int, opts ...domain.LedgerOption) []*domain.Ledger {
	if count < 1 {
		count = 1
	}
	ledgers := make([]*domain.Ledger, 0, count)
	for range count {
		made := l.make(opts...)
		ledgers = append(ledgers, made)
	}

	return ledgers
}

func WithLedgerID(id xuid.XUID) domain.LedgerOption {
	return func(l *domain.Ledger) {
		l.ID = id
	}
}

func WithLedgerDescription(description string) domain.LedgerOption {
	return func(l *domain.Ledger) {
		l.Description = description
	}
}

func WithLedgerName(name string) domain.LedgerOption {
	return func(l *domain.Ledger) {
		l.Name = name
	}
}
