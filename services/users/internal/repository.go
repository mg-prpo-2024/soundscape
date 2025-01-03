package internal

import "gorm.io/gorm"

type Repository interface {
	SetCustomerId(userId uint, customerId string)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) SetCustomerId(userId uint, customerId string) {
	r.db.Model(&User{}).Where("id = ?", userId).Update("stripe_customer_id", customerId)
}
