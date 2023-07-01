package core

import (
	"time"

	"github.com/hoangtk0100/app-context/util"
)

type SQLModel struct {
	ID        int        `json:"-" gorm:"column:id;" db:"id"`
	FakeID    *util.UID  `json:"id" gorm:"-"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at;" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at;" db:"updated_at"`
}

func NewSQLModel() SQLModel {
	now := time.Now().UTC()

	return SQLModel{
		ID:        0,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

func (sqlModel *SQLModel) Mask(objectId int) {
	uid := util.NewUID(uint32(sqlModel.ID), objectId, 1)
	sqlModel.FakeID = &uid
}
