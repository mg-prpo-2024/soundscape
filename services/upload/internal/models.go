package internal

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID               uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Auth0Id          string    `gorm:"unique"`
	Email            string
	StripeCustomerId sql.NullString
	CreatedAt        time.Time // Automatically managed by GORM for creation time
	UpdatedAt        time.Time // Automatically managed by GORM for update time
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}
}
