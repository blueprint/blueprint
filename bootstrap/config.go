package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

// Config bootstraps configuration and environment values.
func Config() error {
	env := viper.GetString("env")
	if len(env) > 0 {
		environments := viper.GetStringMapString("environments." + env)
		if len(environments) < 1 && env != "dev" {
			log.Printf("environment '%s' not found in config file", env)
		}

		for key, value := range environments {
			viper.Set(key, value)
		}
	}
	return nil
}
