package model

import (
	"gorm.io/plugin/soft_delete"
)

type Network struct {
	ID                 uint                  `json:"id" gorm:"primarykey"`
	CreatedAt          int64                 `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          int64                 `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt          soft_delete.DeletedAt `json:"-"`
	ContractAddressURL string                `json:"contract_address_url"`
	StateDumpURL       string                `json:"state_dump_url"`
	Status             bool                  `json:"status"`
}

type Component struct {
	ID        uint                  `json:"id" gorm:"primarykey"`
	CreatedAt int64                 `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64                 `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt soft_delete.DeletedAt `json:"-"`
	Name      string                `json:"name"`
	Type      string                `json:"type" form:"type" binding:"required"`
	Status    bool                  `json:"status"`
	PublicURL string                `json:"public_url"`
	NetworkID uint                  `json:"network_id" form:"network_id" binding:"required"`
	Network   Network               `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
