package services

import (
	"fmt"

	"github.com/revel/revel"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services/auth"
)

// CreateIfAdmin creates named Admin User
func CreateIfAdmin() {
	if o, ok := Model.Logins[entities.Local][auth.Digest(Secret.Admin.UserName)]; ok {
		if o.Level != entities.Admin {
			panic(fmt.Errorf("admin %s exists, but lv is not admin but %d", Secret.Admin.UserName, o.Level))
		}
		return
	} else {
		if o, err := PasswordSignUp(Secret.Admin.UserName, "admin", Secret.Admin.Password, 0); err != nil {
			panic(err)
		} else {
			o.Level = entities.Admin
			revel.AppLog.Info("administrator was successfully created")
		}
	}
}
