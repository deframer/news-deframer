package service

import (
	"fmt"

	"github.com/deframer/news-deframer/pkg/config"
)

func basicAuthFromConfig(cfg *config.Config, user, pass string) error {
	if cfg == nil {
		return nil
	}

	userSet := cfg.BasicAuthUser != ""
	passSet := cfg.BasicAuthPassword != ""
	if !userSet && !passSet {
		return nil
	}
	if userSet && user != cfg.BasicAuthUser {
		return fmt.Errorf("unauthorized")
	}
	if passSet && pass != cfg.BasicAuthPassword {
		return fmt.Errorf("unauthorized")
	}
	return nil
}
