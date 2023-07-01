package core

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type Image struct {
	ID       int    `json:"id" gorm:"column:id;" db:"id"`
	Name     string `json:"name" gorm:"column:name;" db:"name"`
	Path     string `json:"path" gorm:"column:path;" db:"path"`
	URL      string `json:"url" gorm:"-" db:"-"`
	Width    int    `json:"width" gorm:"column:width;" db:"width"`
	Height   int    `json:"height" gorm:"column:height;" db:"height"`
	Provider string `json:"provider,omitempty" gorm:"column:provider;" db:"provider"`
}

func (Image) TableName() string { return "images" }

func (img *Image) Fulfill(domain string) {
	img.URL = fmt.Sprintf("%s/%s", domain, img.Path)
}

func (img *Image) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.WithStack(errors.New(fmt.Sprintf("Failed to unmarshal data from DB: %s", value)))
	}

	var i Image
	if err := json.Unmarshal(bytes, &img); err != nil {
		return errors.WithStack(err)
	}

	*img = i
	return nil
}

func (img *Image) Value() (driver.Value, error) {
	if img == nil {
		return nil, nil
	}
	return json.Marshal(img)
}

type Images []Image

func (i *Images) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.WithStack(errors.New(fmt.Sprintf("Failed to unmarshal JSONB value: %s", value)))
	}

	var data []Image
	if err := json.Unmarshal(bytes, &data); err != nil {
		return errors.WithStack(err)
	}

	*i = data
	return nil
}

func (i *Images) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
	}

	return json.Marshal(i)
}
