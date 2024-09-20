package types

import (
	"time"

	"github.com/google/uuid"
)

type EventType struct {
	EventID     uuid.UUID   `bson:"_id,omitempty" json:"event_id,omitempty" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string      `bson:"name" json:"name" validate:"required" gorm:"not null"`
	Description string      `bson:"description,omitempty" json:"description,omitempty" gorm:"type:text"`
	Date        time.Time   `bson:"date" json:"date" validate:"required" gorm:"not null"`
	Location    string      `bson:"location,omitempty" json:"location,omitempty" gorm:"type:varchar(255)"`
	Capacity    int         `bson:"capacity" json:"capacity" validate:"required,min=1" gorm:"not null"`
	CreatedAt   time.Time   `bson:"created_at" json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `bson:"updated_at" json:"updated_at" gorm:"autoUpdateTime"`
	OrganizerID uuid.UUID   `bson:"organizer_id" json:"organizer_id" gorm:"type:uuid;not null"`
	Attendees   []uuid.UUID `bson:"attendees" json:"attendees" gorm:"type:uuid[]"`
}
