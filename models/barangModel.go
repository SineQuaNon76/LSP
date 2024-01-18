package models

import (
	"time"
)

type Barang struct {
	ID         int       `gorm:"primary_key" json:"id" `
	NamaBarang string    `gorm:"type:varchar(144); not null" json:"nama_barang"`
	Stock      int       `gorm:"type:int(8); not null" json:"stock"`
	Harga      int       `gorm:"type:int(8); not null" json:"harga"`
	Created_At time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"-"`
	Update_At  time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"-"`
}

type User struct {
	ID        int       `gorm:"primary_key" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username" validate:"required"`
	Password  string    `gorm:"not null" json:"password" validate:"required,min=8"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"-"`
}
