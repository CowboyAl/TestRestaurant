package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
)

// OrderID, UserID, ItemID, PickedUp, Delivered

//Order Struct containing the login data
type Order struct {
	ID        string `gorm:"size:20;primary_key"`
	UserID    string
	ItemID    int
	PickedUp  bool
	Delivered bool
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt *time.Time `sql:"index"`
}

// BeforeCreate - make sure to generate a guid for the ID column
func (u *Order) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", xid.New().String())
	return nil
}
