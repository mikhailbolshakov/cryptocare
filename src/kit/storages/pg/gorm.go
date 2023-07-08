package pg

import (
	"gorm.io/gorm"
	"time"
)

// GormDto specifies base attrs for GORM dto
type GormDto struct {
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *gorm.DeletedAt
}

// StringToNull transforms empty string to nil string, so that gorm stores it as NULL
func StringToNull(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// NullToString transforms NULL to empty string
func NullToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
