package db

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
