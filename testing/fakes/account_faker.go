package fakes

import (
	"ledger/internal/domain"

	"github.com/47monad/xuid"
	"github.com/brianvoe/gofakeit/v7"
)

type AccountFakerContext struct{}
type AccountFaker struct{}

func NewAccountFaker() AccountFaker {
	return AccountFaker{}
}

func (a AccountFaker) make(opts ...domain.AccountOption) *domain.Account {
	id, err := xuid.NewSortable(domain.AccountIDPrefix)
	if err != nil {
		panic(err)
	}

	ledgerID, err := xuid.NewSortable(domain.LedgerIDPrefix)
	if err != nil {
		panic(err)
	}

	code := gofakeit.AchAccount()

	name := gofakeit.CelebrityBusiness()

	accountType := domain.AccountType(gofakeit.RandomString([]string{
		string(domain.AccountTypeAsset),
		string(domain.AccountTypeEquity),
		string(domain.AccountTypeExpense),
		string(domain.AccountTypeLiability),
		string(domain.AccountTypeRevenue),
	}))

	normalBalance := domain.Direction(gofakeit.RandomString([]string{
		string(domain.DirectionCredit),
		string(domain.DirectionDebit),
	}))

	accountStatus := domain.AccountStatus(gofakeit.RandomString([]string{
		string(domain.AccountStatusActive),
		string(domain.AccountStatusClosed),
		string(domain.AccountStatusFrozen),
	}))

	numberOfItems := gofakeit.IntRange(1, 10)
	metadata := make(map[string]string, numberOfItems)
	for len(metadata) < numberOfItems {
		metadata[gofakeit.Word()] = gofakeit.Sentence(5)
	}

	allOpts := append([]domain.AccountOption{
		domain.WithAccountType(accountType),
		domain.WithAccountNormalBalance(normalBalance),
		domain.WithAccountStatus(accountStatus),
		domain.WithAccountMetaData(metadata),
	}, opts...)

	made, err := domain.NewAccountFactory().Make(id, ledgerID, code, name, allOpts...)
	if err != nil {
		panic(err)
	}

	return made
}
