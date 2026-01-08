package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID      `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time      // not null default now
	UpdatedAt time.Time      // not null default false
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return nil
}

type Feed struct {
	// gorm.Model
	Base
	URL     string `gorm:"index"`
	Enabled bool   `gorm:"not null;default:false"` // not null default false
}

// StringArray aliases []string to implement sql.Scanner and driver.Valuer for PostgreSQL text[]
type StringArray []string

func (a *StringArray) Scan(src interface{}) error {
	return (*pq.StringArray)(a).Scan(src)
}

func (a StringArray) Value() (driver.Value, error) {
	return pq.StringArray(a).Value()
}

type CachedFeed struct {
	Base
	FeedID    uuid.UUID   `gorm:"type:uuid;uniqueIndex;not null"` // not null
	Feed      Feed        `gorm:"foreignKey:FeedID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	XMLHeader string      `gorm:"type:text;not null"`                // not null
	ItemRefs  StringArray `gorm:"type:text[];not null;default:'{}'"` // not null default []
}

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

type Item struct {
	Base
	Hash         string    `gorm:"type:char(64);uniqueIndex:idx_hash_feed_url;uniqueIndex:idx_hash_feed;not null"`                      // not null
	FeedID       uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_feed_url;uniqueIndex:idx_hash_feed_url;uniqueIndex:idx_hash_feed;not null"` // not null
	Feed         Feed      `gorm:"foreignKey:FeedID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	URL          string    `gorm:"uniqueIndex:idx_feed_url;uniqueIndex:idx_hash_feed_url;not null"` // not null
	AIResult     JSONB     `gorm:"type:jsonb;not null"`                                             // not null
	DebugContent string    `gorm:"type:text;default:null"`                                          // null
	MinHash      string    `gorm:"default:null"`                                                    // null
}
