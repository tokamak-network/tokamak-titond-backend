package model

import (
	"gorm.io/plugin/soft_delete"
)

type Network struct {
	ID                  uint                  `json:"id" gorm:"primarykey"`
	CreatedAt           int64                 `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           int64                 `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt           soft_delete.DeletedAt `json:"-"`
	ContractAddressURL  string                `json:"contract_address_url"`
	StateDumpURL        string                `json:"state_dump_url"`
	LatestComponentType string                `json:"latest_component" gorm:"-"`
	Status              bool                  `json:"status" gorm:"-"`
}

type Component struct {
	ID        uint                  `json:"id" gorm:"primarykey"`
	CreatedAt int64                 `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64                 `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt soft_delete.DeletedAt `json:"-"`
	Name      string                `json:"name"`
	Type      string                `json:"type"`
	Status    bool                  `json:"status" gorm:"-"`
	PublicURL string                `json:"public_url" gorm:"-"`
	NetworkID uint                  `json:"network_id"`
	Network   Network               `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
