// Package config defines the global config instance.
package config

import (
	"log"
	"sync"

	"github.com/che-kwas/iam-kit/config"
)

var (
	cfg  *config.Config
	once sync.Once
)

// Cfg returns the global config instance.
func Cfg() *config.Config {
	if cfg == nil {
		log.Fatal("Config not initialized")
	}

	return cfg
}

func InitConfig(cfgPath, appName string) error {
	if cfg != nil {
		return nil
	}

	var err error
	once.Do(func() {
		cfg, err = config.NewConfig(cfgPath, appName)
	})

	return err
}
