package internal

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

/*
export const usersTable = pgTable("users", {
  id: uuid().primaryKey().defaultRandom(),
  name: varchar({ length: 255 }).notNull(),
  email: varchar({ length: 255 }).notNull().unique(),
  age: integer().notNull(),
  imageUrl: varchar({ length: 255 }),
  stripeCustomerId: varchar({ length: 255 }).notNull().unique(),
  createdAt: timestamp().defaultNow(),
});

export const plansTable = pgTable("plans", {
  id: uuid().primaryKey().defaultRandom(),
  name: varchar().notNull().unique(),
  description: varchar().notNull(),
  // priceCents: integer().notNull(), // TODO: this would be different depending on country
  stripeProductId: varchar().notNull().unique(),
  createdAt: timestamp().defaultNow(),
});

*/

type User struct {
	ID               uint // Standard field for the primary key
	StripeCustomerId sql.NullString
	CreatedAt        time.Time // Automatically managed by GORM for creation time
	UpdatedAt        time.Time // Automatically managed by GORM for update time
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}
}
