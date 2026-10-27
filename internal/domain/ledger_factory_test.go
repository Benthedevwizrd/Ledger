package domain_test

import (
	"ledger/internal/domain"
	"testing"
	"time"

	"github.com/47monad/xuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Make(t *testing.T) {
	t.Run("main", func(t *testing.T) {
		id1 := xuid.MustNewSortable("ledger_")
		name1 := "some_ledger"
		description1 := "some_ledger_description"

		now := time.Now().UTC()
		customCreatedAt := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
		customUpdatedAt := time.Date(2024, 6, 7, 8, 9, 10, 0, time.UTC)
		metadata1 := map[string]string{"owner": "ali", "env": "test"}

		testCases := []struct {
			name             string
			givenID          xuid.XUID
			givenName        string
			givenDescription string
			opts             []domain.LedgerOptions
			wanted           *domain.Ledger
			wantedError      error
		}{
			// ---- success ----
			{
				name:             "valid_inputs/success",
				givenID:          id1,
				givenName:        name1,
				givenDescription: description1,
				wanted: &domain.Ledger{
					ID:          id1,
					Name:        name1,
					Description: description1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			{
				name:             "with_metadata/success",
				givenID:          id1,
				givenName:        name1,
				givenDescription: description1,
				opts:             []domain.LedgerOptions{domain.WithLedgerMetadata(metadata1)},
				wanted: &domain.Ledger{
					ID:          id1,
					Name:        name1,
					Description: description1,
					Metadata:    metadata1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			{
				name:             "with_created_at/success",
				givenID:          id1,
				givenName:        name1,
				givenDescription: description1,
				opts:             []domain.LedgerOptions{domain.WithLedgerCreatedAt(customCreatedAt)},
				wanted: &domain.Ledger{
					ID:          id1,
					Name:        name1,
					Description: description1,
					CreatedAt:   customCreatedAt,
					UpdatedAt:   now,
				},
			},
			{
				name:             "with_updated_at/success",
				givenID:          id1,
				givenName:        name1,
				givenDescription: description1,
				opts:             []domain.LedgerOptions{domain.WithLedgerUpdatedAt(customUpdatedAt)},
				wanted: &domain.Ledger{
					ID:          id1,
					Name:        name1,
					Description: description1,
					CreatedAt:   now,
					UpdatedAt:   customUpdatedAt,
				},
			},
			{
				name:             "with_all_options/success",
				givenID:          id1,
				givenName:        name1,
				givenDescription: description1,
				opts: []domain.LedgerOptions{
					domain.WithLedgerMetadata(metadata1),
					domain.WithLedgerCreatedAt(customCreatedAt),
					domain.WithLedgerUpdatedAt(customUpdatedAt),
				},
				wanted: &domain.Ledger{
					ID:          id1,
					Name:        name1,
					Description: description1,
					Metadata:    metadata1,
					CreatedAt:   customCreatedAt,
					UpdatedAt:   customUpdatedAt,
				},
			},
			{
				name:             "multiple_options/last_option_wins",
				givenID:          id1,
				givenName:        name1,
				givenDescription: description1,
				opts: []domain.LedgerOptions{
					domain.WithLedgerCreatedAt(customUpdatedAt),
					domain.WithLedgerCreatedAt(customCreatedAt),
				},
				wanted: &domain.Ledger{
					ID:          id1,
					Name:        name1,
					Description: description1,
					CreatedAt:   customCreatedAt,
					UpdatedAt:   now,
				},
			},

			// ---- validation errors ----
			{
				name:             "empty_id/ErrEmptyLedgerID",
				givenID:          xuid.XUID{},
				givenName:        name1,
				givenDescription: description1,
				wantedError:      domain.ErrEmptyLedgerID,
			},
			{
				name:             "empty_name/ErrEmptyLedgerName",
				givenID:          id1,
				givenName:        "",
				givenDescription: description1,
				wantedError:      domain.ErrEmptyLedgerName,
			},
			{
				name:             "empty_description/ErrEmptyLedgerDescription",
				givenID:          id1,
				givenName:        name1,
				givenDescription: "",
				wantedError:      domain.ErrEmptyLedgerDescription,
			},
			{
				name:             "empty_id_and_name/id_error_comes_first",
				givenID:          xuid.XUID{},
				givenName:        "",
				givenDescription: description1,
				wantedError:      domain.ErrEmptyLedgerID,
			},
			{
				name:             "empty_name_and_description/name_error_comes_first",
				givenID:          id1,
				givenName:        "",
				givenDescription: "",
				wantedError:      domain.ErrEmptyLedgerName,
			},
			{
				name:             "all_empty/id_error_comes_first",
				givenID:          xuid.XUID{},
				givenName:        "",
				givenDescription: "",
				wantedError:      domain.ErrEmptyLedgerID,
			},
			{
				name:             "empty_name/ErrEmptyLedgerName_even_with_valid_options",
				givenID:          id1,
				givenName:        "",
				givenDescription: description1,
				opts:             []domain.LedgerOptions{domain.WithLedgerMetadata(metadata1)},
				wantedError:      domain.ErrEmptyLedgerName,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// act
				result, err := domain.NewLedgerFactory().Make(tc.givenID, tc.givenName, tc.givenDescription, tc.opts...)

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
				assert.Equal(t, tc.wanted.Name, result.Name)
				assert.Equal(t, tc.wanted.Description, result.Description)
				assert.Equal(t, tc.wanted.Metadata, result.Metadata)
				assert.WithinDuration(t, tc.wanted.CreatedAt, result.CreatedAt, 1*time.Second)
				assert.WithinDuration(t, tc.wanted.UpdatedAt, result.UpdatedAt, 1*time.Second)
			})
		}
	})
}
