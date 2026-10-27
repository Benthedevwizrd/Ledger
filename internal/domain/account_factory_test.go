package domain_test

import (
	"ledger/internal/domain"
	"testing"
	"time"

	"github.com/47monad/xuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AccountMake(t *testing.T) {
	t.Run("main", func(t *testing.T) {
		id1 := xuid.MustNewSortable("account_")
		ledgerID1 := xuid.MustNewSortable("ledger_")
		code1 := "1010"
		name1 := "cash"

		now := time.Now().UTC()
		customCreatedAt := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
		customUpdatedAt := time.Date(2024, 6, 7, 8, 9, 10, 0, time.UTC)

		testCases := []struct {
			name        string
			givenID     xuid.XUID
			givenLedger xuid.XUID
			givenCode   string
			givenName   string
			opts        []domain.AccountOption
			wanted      *domain.Account
			wantedError error
		}{
			// ---- success ----
			{
				name:        "valid_inputs/success",
				givenID:     id1,
				givenLedger: ledgerID1,
				givenCode:   code1,
				givenName:   name1,
				wanted: &domain.Account{
					ID:            id1,
					LedgerID:      ledgerID1,
					Code:          code1,
					Name:          name1,
					Type:          domain.AccountTypeAsset,
					NormalBalance: domain.DirectionCredit,
					CreatedAt:     now,
					UpdatedAt:     now,
				},
			},
			{
				name:        "multiple_opts/success",
				givenID:     id1,
				givenLedger: ledgerID1,
				givenCode:   code1,
				givenName:   name1,
				opts: []domain.AccountOption{
					domain.WithAccountType(domain.AccountTypeLiability),
					domain.WithAccountNormalBalance(domain.DirectionCredit),
				},
				wanted: &domain.Account{
					ID:            id1,
					LedgerID:      ledgerID1,
					Code:          code1,
					Name:          name1,
					Type:          domain.AccountTypeLiability,
					NormalBalance: domain.DirectionCredit,
					CreatedAt:     now,
					UpdatedAt:     now,
				},
			},
			{
				name:        "WithAccountCreatedAT/success",
				givenID:     id1,
				givenLedger: ledgerID1,
				givenCode:   code1,
				givenName:   name1,
				opts:        []domain.AccountOption{domain.WithAccountCreatedAt(customCreatedAt)},
				wanted: &domain.Account{
					ID:            id1,
					LedgerID:      ledgerID1,
					Code:          code1,
					Name:          name1,
					Type:          domain.AccountTypeAsset,
					NormalBalance: domain.DirectionCredit,
					CreatedAt:     customCreatedAt,
					UpdatedAt:     now,
				},
			},
			{
				name:        "withAccountUpdatedAt/success",
				givenID:     id1,
				givenLedger: ledgerID1,
				givenCode:   code1,
				givenName:   name1,
				opts:        []domain.AccountOption{domain.WithAccountUpdatedAt(customUpdatedAt)},
				wanted: &domain.Account{
					ID:            id1,
					LedgerID:      ledgerID1,
					Code:          code1,
					Name:          name1,
					Type:          domain.AccountTypeAsset,
					NormalBalance: domain.DirectionCredit,
					CreatedAt:     now,
					UpdatedAt:     customUpdatedAt,
				},
			},
			{
				name:        "all_opts_valid_inputs/success",
				givenID:     id1,
				givenLedger: ledgerID1,
				givenCode:   code1,
				givenName:   name1,
				opts: []domain.AccountOption{
					domain.WithAccountType(domain.AccountTypeRevenue),
					domain.WithAccountNormalBalance(domain.DirectionCredit),
					domain.WithAccountCreatedAt(customCreatedAt),
					domain.WithAccountUpdatedAt(customUpdatedAt),
				},
				wanted: &domain.Account{
					ID:            id1,
					LedgerID:      ledgerID1,
					Code:          code1,
					Name:          name1,
					Type:          domain.AccountTypeRevenue,
					NormalBalance: domain.DirectionCredit,
					CreatedAt:     customCreatedAt,
					UpdatedAt:     customUpdatedAt,
				},
			},
			{
				name:        "multiple_opts_for_the_type/last_one_wins",
				givenID:     id1,
				givenLedger: ledgerID1,
				givenCode:   code1,
				givenName:   name1,
				opts: []domain.AccountOption{
					domain.WithAccountType(domain.AccountTypeLiability),
					domain.WithAccountType(domain.AccountTypeEquity),
				},
				wanted: &domain.Account{
					ID:            id1,
					LedgerID:      ledgerID1,
					Code:          code1,
					Name:          name1,
					Type:          domain.AccountTypeEquity,
					NormalBalance: domain.DirectionCredit,
					CreatedAt:     now,
					UpdatedAt:     now,
				},
			},

			// ---- validation errors ----
			{
				name:        "empty_id/ErrEmptyAccountID",
				givenID:     xuid.XUID{},
				givenLedger: ledgerID1,
				givenCode:   code1,
				givenName:   name1,
				wantedError: domain.ErrEmptyAccountID,
			},
			{
				name:        "empty_ledgerID/ErrEmptyLedgerID",
				givenID:     id1,
				givenLedger: xuid.XUID{},
				givenCode:   code1,
				givenName:   name1,
				wantedError: domain.ErrEmptyAccountLedgerID,
			},
			{
				name:        "empty_code/ErrEmptyAccountCode",
				givenID:     id1,
				givenLedger: ledgerID1,
				givenCode:   "",
				givenName:   name1,
				wantedError: domain.ErrEmptyAccountCode,
			},
			{
				name:        "empty_name/ErrEmptyAccountName",
				givenID:     id1,
				givenLedger: ledgerID1,
				givenCode:   code1,
				givenName:   "",
				wantedError: domain.ErrEmptyAccountName,
			},
			{
				name:        "Empty_id_and_ledger_id/ErrEmptyAccountID",
				givenID:     xuid.XUID{},
				givenLedger: xuid.XUID{},
				givenCode:   code1,
				givenName:   name1,
				wantedError: domain.ErrEmptyAccountID,
			},
			{
				name:        "Empty_ledger_id_and_code/ErrEmptyLedgerID",
				givenID:     id1,
				givenLedger: xuid.XUID{},
				givenCode:   "",
				givenName:   name1,
				wantedError: domain.ErrEmptyAccountLedgerID,
			},
			{
				name:        "Empty_code_and_name/ErrEmptyAccountCode",
				givenID:     id1,
				givenLedger: ledgerID1,
				givenCode:   "",
				givenName:   "",
				wantedError: domain.ErrEmptyAccountCode,
			},
			{
				name:        "all_mandatory_fields_empty/ErrEmptyAccountID",
				givenID:     xuid.XUID{},
				givenLedger: xuid.XUID{},
				givenCode:   "",
				givenName:   "",
				wantedError: domain.ErrEmptyAccountID,
			},
			{
				name:        "multiple_valid_opts/empty_name/ErrEmptyAccountName",
				givenID:     id1,
				givenLedger: ledgerID1,
				givenCode:   code1,
				givenName:   "",
				opts: []domain.AccountOption{
					domain.WithAccountType(domain.AccountTypeLiability),
					domain.WithAccountNormalBalance(domain.DirectionCredit),
				},
				wantedError: domain.ErrEmptyAccountName,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// act
				result, err := domain.NewAccountFactory().Make(tc.givenID, tc.givenLedger, tc.givenCode, tc.givenName, tc.opts...)

				// assert
				if tc.wantedError != nil {
					require.ErrorIs(t, err, tc.wantedError)
					assert.Nil(t, result)
					return
				}

				require.NoError(t, err)
				require.NotNil(t, result)
				require.NotNil(t, tc.wanted)

				assert.Equal(t, tc.wanted.ID, result.ID)
				assert.Equal(t, tc.wanted.LedgerID, result.LedgerID)
				assert.Equal(t, tc.wanted.Code, result.Code)
				assert.Equal(t, tc.wanted.Name, result.Name)
				assert.Equal(t, tc.wanted.Type, result.Type)
				assert.Equal(t, tc.wanted.NormalBalance, result.NormalBalance)
				assert.WithinDuration(t, tc.wanted.CreatedAt, result.CreatedAt, 1*time.Second)
				assert.WithinDuration(t, tc.wanted.UpdatedAt, result.UpdatedAt, 1*time.Second)
			})
		}
	})
}
