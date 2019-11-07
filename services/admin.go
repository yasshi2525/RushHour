package services

import (
	"fmt"
	"log"

	"github.com/yasshi2525/RushHour/entities"
)

// CreateIfAdmin creates named Admin User
func CreateIfAdmin() {
	if o, ok := Model.Logins[entities.Local][serviceConf.Auther.Digest(serviceConf.AppConf.Secret.Admin.UserName)]; ok {
		if o.Level != entities.Admin {
			panic(fmt.Errorf("admin %s exists, but lv is not admin but %d", serviceConf.AppConf.Secret.Admin.UserName, o.Level))
		}
		return
	}
	if o, err := PasswordSignUp(serviceConf.AppConf.Secret.Admin.UserName, "admin", serviceConf.AppConf.Secret.Admin.Password, 0, entities.Admin); err != nil {
		log.Println("skip create administrator because already exist")
	} else {
		o.Level = entities.Admin
		log.Println("administrator was successfully created")
	}
}
