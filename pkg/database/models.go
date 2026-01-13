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
	URL               string        `gorm:"index"`                 // we can't enforce uniqueness here (because of the soft deletes)
	RootDomain        *string       `gorm:"index"`                 // example.com
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

// ThinkResult we make omitempty to not serialize default e.g. 0.0 or ""
type ThinkResult struct {
	TitleOriginal               string  `json:"title_original,omitempty"`
	DescriptionOriginal         string  `json:"description_original,omitempty"`
	TitleCorrected              string  `json:"title_corrected,omitempty"`
	TitleCorrectionReason       string  `json:"title_correction_reason,omitempty"`
	DescriptionCorrected        string  `json:"description_corrected,omitempty"`
	DescriptionCorrectionReason string  `json:"description_correction_reason,omitempty"`
	Framing                     float64 `json:"framing,omitempty"`
	FramingReason               string  `json:"framing_reason,omitempty"`
	Clickbait                   float64 `json:"clickbait,omitempty"`
	ClickbaitReason             string  `json:"clickbait_reason,omitempty"`
	Persuasive                  float64 `json:"persuasive,omitempty"`
	PersuasiveReason            string  `json:"persuasive_reason,omitempty"`
	HyperStimulus               float64 `json:"hyper_stimulus,omitempty"`
	HyperStimulusReason         string  `json:"hyper_stimulus_reason,omitempty"`
	Speculative                 float64 `json:"speculative,omitempty"`
	SpeculativeReason           string  `json:"speculative_reason,omitempty"`
	OverallReason               string  `json:"overall_reason,omitempty"`
}

func (j ThinkResult) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *ThinkResult) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

type MediaThumbnail struct {
	URL    string `xml:"url,attr" json:"url,omitempty"`
	Height int    `xml:"height,attr" json:"height,omitempty"`
	Width  int    `xml:"width,attr" json:"width,omitempty"`
}

type MediaContent struct {
	// Technical attributes required for <img src> or <video src>
	URL    string `xml:"url,attr" json:"url,omitempty"`
	Type   string `xml:"type,attr" json:"type,omitempty"`     // e.g., "image/jpeg"
	Medium string `xml:"medium,attr" json:"medium,omitempty"` // e.g., "image" or "video"

	// Dimensions are essential for preventing HTML Layout Shift (CLS)
	Height int `xml:"height,attr" json:"height,omitempty"`
	Width  int `xml:"width,attr" json:"width,omitempty"`

	// Descriptive text used for HTML 'alt' attributes and captions
	Title       string `xml:"title" json:"title,omitempty"`
	Description string `xml:"description" json:"description,omitempty"`

	// Essential for Video: Used for the <video poster="..."> attribute
	// We use a pointer so it is nil if no thumbnail exists
	Thumbnail *MediaThumbnail `xml:"thumbnail" json:"thumbnail,omitempty"`

	// Optional Copyright for the content
	Credit string `xml:"credit" json:"credit,omitempty"`
}

func (j MediaContent) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *MediaContent) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

type Item struct {
	ID           uuid.UUID     `gorm:"primaryKey;type:uuid"`
	CreatedAt    time.Time     `gorm:"not null;default:now()"`
	UpdatedAt    time.Time     `gorm:"not null;default:now()"`
	Hash         string        `gorm:"type:char(64);uniqueIndex:idx_hash_feed_url;uniqueIndex:idx_hash_feed;not null"`
	FeedID       uuid.UUID     `gorm:"type:uuid;index;uniqueIndex:idx_feed_url;uniqueIndex:idx_hash_feed_url;uniqueIndex:idx_hash_feed;not null"`
	Feed         Feed          `gorm:"foreignKey:FeedID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	URL          string        `gorm:"index;uniqueIndex:idx_feed_url;uniqueIndex:idx_hash_feed_url;not null"`
	Content      string        `gorm:"type:text;not null"`
	PubDate      time.Time     `gorm:"not null;index;default:now()"`
	MediaContent *MediaContent `gorm:"type:jsonb"`
	ThinkResult  *ThinkResult  `gorm:"type:jsonb"`
	ThinkError   *string       `gorm:"type:text;null"`
	ThinkRating  float64       `gorm:"not null;default:0.0"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (item *Item) BeforeCreate(tx *gorm.DB) error {
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	return nil
}
