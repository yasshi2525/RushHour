package config

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"gopkg.in/go-playground/validator.v9"
)

// Config stores configurable variables
type Config struct {
	Game   CnfGame
	Secret CnfSecret
}

// Load load and validate game.conf/secret.conf
func Load(confDir string) (Config, error) {
	config := Config{}
	route := map[string]interface{}{
		"game.conf":   &config.Game,
		"secret.conf": &config.Secret,
	}
	for file, v := range route {
		if _, err := toml.DecodeFile(fmt.Sprintf("%s/%s", confDir, file), v); err != nil {
			return config, fmt.Errorf("failed to load conf: %v", err)
		}
		if err := validator.New().Struct(v); err != nil {
			return config, fmt.Errorf("%+v, %v", v, err)
		}
		log.Println("config file was successfully loaded.")
	}

	return config, nil
}
