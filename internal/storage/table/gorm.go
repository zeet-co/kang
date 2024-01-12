package table

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type GormBase struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (base *GormBase) BeforeCreate(tx *gorm.DB) error {
	if base.ID != uuid.Nil {
		return nil
	}
	uuid, err := uuid.NewRandom()
	if err != nil {
		return errors.WithStack(err)
	}
	base.ID = uuid
	return nil
}

func (b *GormBase) GetID() uuid.UUID {
	return b.ID
}

type UUIDSlice []uuid.UUID

func (u UUIDSlice) Value() (driver.Value, error) {
	return json.Marshal(u)
}

func (u *UUIDSlice) Scan(value interface{}) error {
	if value == nil {
		*u = []uuid.UUID{}
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return schema.ErrUnsupportedDataType
	}

	return json.Unmarshal(data, u)
}
