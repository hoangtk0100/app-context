package core

import (
	"time"

	"github.com/hoangtk0100/app-context/util"
)

type SQLModel struct {
	Id        int        `json:"-" gorm:"column:id;" db:"id"`
	FakeId    *util.UID  `json:"id" gorm:"-"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at;" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at;" db:"updated_at"`
}

func NewSQLModel() SQLModel {
	now := time.Now().UTC()

	return SQLModel{
		Id:        0,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

func (sqlModel *SQLModel) Mask(objectId int) {
	uid := util.NewUID(uint32(sqlModel.Id), objectId, 1)
	sqlModel.FakeId = &uid
}
