package db

import (
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Postgres struct {
	gDB *gorm.DB
}

func NewPostgresql(cfg *Config) (*Postgres, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	gdB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	gdB.AutoMigrate(&model.Network{}, &model.Component{})

	return &Postgres{gdB}, err
}

func NewMockPostgres() (*Postgres, sqlmock.Sqlmock) {
	mockDB, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	return &Postgres{db}, mock
}
