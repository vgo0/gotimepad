package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OTOptions struct {
	DB *string // otp db location
	Message *string // message file location
	B64 *bool // is base64
	MessageDirect *string // is passed via command line
	ToConsole *bool // output to console
	NoValidateExists bool // default false, used to skip file existence check in testing if needed
}

type OTPage struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`
	Page []byte `gorm:"type:blob"`
	Used bool
}

// Generates our UUID as the primary key
func (page *OTPage) BeforeCreate(tx *gorm.DB) (err error) {
	// this will panic if it fails
	page.ID = uuid.New()

	return nil
}