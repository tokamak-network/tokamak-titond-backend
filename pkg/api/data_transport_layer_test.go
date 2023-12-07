package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/db"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

func TestCreateDTL(t *testing.T) {
	fakek8s := kubernetes.NewFakeKubernetes()
	mockPostgres, mock := db.NewMockPostgres()

	mockAPI := &TitondAPI{
		k8s: fakek8s,
		db:  mockPostgres,
	}

	dtl := &model.Component{
		Name:      "Test-DTL",
		Type:      "data-transport-layer",
		NetworkID: 1,
	}

	rows := sqlmock.NewRows([]string{"contract_address_url"}).AddRow("test.url")
	mock.ExpectQuery(`SELECT(.*)`).
		WithArgs(1, 0).
		WillReturnRows(rows)

	rows2 := sqlmock.NewRows([]string{"id", "type", "network_id"}).AddRow(1, dtl.Type, 1)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT(.*)`).
		WithArgs(time.Now().Unix(), time.Now().Unix(), 0, dtl.Name, dtl.Type, false, "", 1).
		WillReturnRows(rows2)
	mock.ExpectCommit()

	result, err := mockAPI.CreateDTL(dtl)
	if err != nil {
		fmt.Printf("error : %s\n", err)
	} else {
		fmt.Printf("result: %v\n", result)
	}

	assert.Nil(t, mock.ExpectationsWereMet())
}
