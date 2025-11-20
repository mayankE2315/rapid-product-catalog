package config

import (
	"fmt"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

//go:generate mockgen -destination mocks/config_mock.go -package config -source=config.go Config

type Config interface {
	Get() *Values
	IsProductionEnv() bool
}

func (c *Values) Get() *Values {
	return c
}

func (c *Values) IsProductionEnv() bool {
	return c.Environment == "production"
}

func NewConfig() (Config, error) {
	viper := viper.NewWithOptions(viper.KeyDelimiter("::"))
	viper.AutomaticEnv()
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("rapido")
	viper.AddConfigPath("config")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath("../../config/")
	viper.AddConfigPath("../../../config/")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	var values Values
	decoderConfig := &mapstructure.DecoderConfig{
		Result:           &values,
		ErrorUnused:      true,
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return nil, fmt.Errorf("fatal error creating decoder: %w", err)
	}

	if err := decoder.Decode(viper.AllSettings()); err != nil {
		return nil, fmt.Errorf("fatal error unable to decode config file: %w", err)
	}

	return &values, nil
}

func (c *Values) LogLevel() string {
	if "" == c.Log.Level {
		return "info"
	}

	return c.Log.Level
}

func (c *Values) ListenAddress() string {
	return ":" + strconv.Itoa(c.Server.Port)
}

var appConfig Config

func GetConfig() Config {
	return appConfig
}

func SetConfig(c Config) {
	appConfig = c
}
