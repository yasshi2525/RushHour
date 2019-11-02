package services

import (
	"fmt"
	"os"

	"github.com/yasshi2525/RushHour/app/services/auth"

	"github.com/BurntSushi/toml"
	"github.com/revel/revel"
	validator "gopkg.in/go-playground/validator.v9"
)

type admin struct {
	UserName string
	Password string
}

type secret struct {
	Admin admin
	Auth  auth.Config
}

// Secret defines secret constant variable
var Secret secret

// LoadSecret load and validate secret.conf
func LoadSecret() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if _, err := toml.DecodeFile(fmt.Sprintf("%s/%s", dir, "conf/secret.conf"), &Secret); err != nil {
		panic(fmt.Errorf("failed to load secret: %v", err))
	}
	if err := validator.New().Struct(Secret); err != nil {
		panic(fmt.Errorf("%+v, %v", Secret, err))
	}

	defer revel.AppLog.Info("secret file was successfully loaded.")
	return
}
