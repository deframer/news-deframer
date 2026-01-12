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

// Base is a gorm.Model but with a UUID id
type Base struct {
	ID        uuid.UUID      `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time      `gorm:"not null;default:now()"`
	UpdatedAt time.Time      `gorm:"not null;default:now()"`
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
	Base
	URL               string        `gorm:"index"`
	EnforceFeedDomain bool          `gorm:"not null;default:true"` // item url must be from our URL
	Enabled           bool          `gorm:"not null;default:false;index"`
	Polling           bool          `gorm:"not null;default:false"`
	FeedSchedule      *FeedSchedule `gorm:"foreignKey:ID;references:ID"`
}

type CachedFeed struct {
	ID         uuid.UUID   `gorm:"primaryKey;type:uuid"` // this is a FK to Feed.ID
	CreatedAt  time.Time   `gorm:"not null;default:now()"`
	UpdatedAt  time.Time   `gorm:"not null;default:now()"`
	XMLContent *string     `gorm:"type:text"`
	ItemRefs   StringArray `gorm:"type:text[];not null;default:'{}'"`
}

// StringArray aliases []string to implement sql.Scanner and driver.Valuer for PostgreSQL text[]
type StringArray []string

func (a *StringArray) Scan(src interface{}) error {
	return (*pq.StringArray)(a).Scan(src)
}

func (a StringArray) Value() (driver.Value, error) {
	return pq.StringArray(a).Value()
}

type FeedSchedule struct {
	ID          uuid.UUID  `gorm:"primaryKey;type:uuid"` // this is a FK to Feed.ID
	CreatedAt   time.Time  `gorm:"not null;default:now()"`
	UpdatedAt   time.Time  `gorm:"not null;default:now()"`
	NextRunAt   *time.Time `gorm:"index"`
	LockedUntil *time.Time
	LastError   *string `gorm:"type:text"`
}

type Item struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid"`
	CreatedAt      time.Time `gorm:"not null;default:now()"`
	UpdatedAt      time.Time `gorm:"not null;default:now()"`
	Hash           string    `gorm:"type:char(64);uniqueIndex:idx_hash_feed_url;uniqueIndex:idx_hash_feed;not null"`
	FeedID         uuid.UUID `gorm:"type:uuid;index;uniqueIndex:idx_feed_url;uniqueIndex:idx_hash_feed_url;uniqueIndex:idx_hash_feed;not null"`
	Feed           Feed      `gorm:"foreignKey:FeedID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	URL            string    `gorm:"index;uniqueIndex:idx_feed_url;uniqueIndex:idx_hash_feed_url;not null"`
	AnalyzerResult *JSONB    `gorm:"type:jsonb"`
	Content        string    `gorm:"type:text;not null"`
	PubDate        time.Time `gorm:"not null;index;default:now()"`
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

// BeforeCreate will set a UUID rather than numeric ID.
func (item *Item) BeforeCreate(tx *gorm.DB) error {
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	return nil
}
