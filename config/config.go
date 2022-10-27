package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Env             string `mapstruct:"ENV"`
	Port            string `mapstructure:"PORT"`
	DBURL           string `mapstructure:"DATABASE_URL"`
	DBName          string `mapstructure:"DATABASE_NAME"`
	MigrationPath 	string `mapstructure:"MIGRATION_PATH"`
	ShouldMigrate 	bool `mapstructure:"SHOULD_MIGRATE"`
}

func LoadConfig() (Config, error) {
	viper.SetConfigFile(".env")

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, nil
		}
		return config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}
	return config, nil
}

func (c *Config) IsLocal() bool {
	return strings.Contains(c.Env, "local")
}