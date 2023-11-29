package db

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	return db, mock
}

func TestCreateNetwork(t *testing.T) {
	db, mock := newMockDB()
	mockPostgres := &Postgres{
		db,
	}

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
	db, mock := newMockDB()
	mockPostgres := &Postgres{
		db,
	}

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
