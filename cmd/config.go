package cmd

import "time"

type Config struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
	Port int `json:"port"`
	DB   struct {
		Url        string `json:"url"`
		MigrateUrl string `json:"migreateUrl" yaml:"migrateUrl"`
	} `json:"db" yaml:"db"`

	ShutdownGracePeriod time.Duration `json:"shutdownGracePeriod"`
}
