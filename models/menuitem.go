package models

import (
	"time"
)

//MenuItem Struct containing the login data
type MenuItem struct {
	ID          int     `gorm:"primary_key"`
	Description string  `gorm:"size:250;not null;unique_index"`
	Price       float32 // price in dollars
	PrepTime    float32 // prep time in minutes
	CreatedAt   time.Time
	UpdatedAt   time.Time
	//DeletedAt *time.Time `sql:"index"`
}
