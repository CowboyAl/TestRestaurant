package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
)

//User Struct containing the login data
type User struct {
	ID        string `gorm:"size:20;primary_key"`
	Username  string `gorm:"size:250;not null;unique_index"`
	Password  string `gorm:"size:80;not null"`
	Address   string
	Distance  float32
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt *time.Time `sql:"index"`
}

// BeforeCreate - make sure to generate a guid for the ID column
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", xid.New().String())
	return nil
}
