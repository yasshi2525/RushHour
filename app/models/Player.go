package models

import "github.com/jinzhu/gorm"

// Player represents user information
type Player struct {
	gorm.Model

	DisplayName string
	Password    string
}
