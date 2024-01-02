package db

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	"gorm.io/gorm"
)

func TestCreateNetwork(t *testing.T) {
	mockPostgres, mock := NewMockPostgres()

	rows := sqlmock.NewRows([]string{"id", "contract_address_url", "state_dump_url"}).AddRow(1, "test_url", "test_url")
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT(.*)`).
		WithArgs(time.Now().Unix(), time.Now().Unix(), 0, "test_url", "test_url", false).
		WillReturnRows(rows)
	mock.ExpectCommit()

	mockPostgres.CreateNetwork(&model.Network{
		ContractAddressURL: "test_url",
		StateDumpURL:       "test_url",
	})

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestReadNetwork(t *testing.T) {
	mockPostgres, mock := NewMockPostgres()

	tests := []struct {
		name         string
		expectedRows *sqlmock.Rows
		expectedArgs uint
		actualArgs   uint
		expectedErr  error
	}{
		{
			name:         "should find rows",
			expectedRows: sqlmock.NewRows([]string{"id"}).AddRow(1),
			expectedArgs: 1,
			actualArgs:   1,
			expectedErr:  nil,
		},
		{
			name:         "should not find rows",
			expectedRows: sqlmock.NewRows([]string{"id"}),
			expectedArgs: 2,
			actualArgs:   2,
			expectedErr:  gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`SELECT(.*)`).
				WithArgs(tt.expectedArgs, 0).
				WillReturnRows(tt.expectedRows)

			_, err := mockPostgres.ReadNetwork(tt.actualArgs)
			if err != nil {
				assert.EqualError(t, tt.expectedErr, "record not found")
			}

			assert.Nil(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReadNetworkByPage(t *testing.T) {
	postgres, mock := NewMockPostgres()

	rows := sqlmock.NewRows([]string{"id", "ContractAddressURL"}).
		AddRow(1, "Network1").
		AddRow(2, "Network2").
		AddRow(3, "Network3")
	offset := 0
	limit := 2
	mock.ExpectQuery(`SELECT(.*)`).WillReturnRows(rows)

	result, err := postgres.ReadNetworkByRange(offset, limit)

	assert.NoError(t, err)
	assert.Len(t, result, 3)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateComponent(t *testing.T) {
	mockPostgres, mock := NewMockPostgres()

	rows := sqlmock.NewRows([]string{"id", "type", "network_id"}).AddRow(1, "l2geth", 1)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT(.*)`).
		WithArgs(time.Now().Unix(), time.Now().Unix(), 0, "test", "l2geth", false, "", 1).
		WillReturnRows(rows)
	mock.ExpectCommit()

	mockPostgres.CreateComponent(&model.Component{
		Name:      "test",
		Type:      "l2geth",
		NetworkID: 1,
	})

	assert.Nil(t, mock.ExpectationsWereMet())
}
