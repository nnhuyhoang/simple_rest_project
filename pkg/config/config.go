package config

import (
	"strings"
	"time"

	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/app"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceName     string
	BaseURL         string
	Port            string
	Env             string
	AllowedOrigins  string
	DBHost          string
	DBPort          string
	DBUser          string
	DBName          string
	DBPass          string
	DBSSLMode       string
	JWTSecret       string
	Timezone        string
	AccessTokenTTL  int64
	RefreshTokenTTL int64
	Version         string
	RabbitURL       string
	RabbitPort      string
	RabbitUser      string
	RabbitPassword  string
}

// JWTConfig save jwt config key and secret
type JWTConfig struct {
	Key    string
	Secret []byte
	Source string
}

// GetCORS in config
func (c *Config) GetCORS() []string {
	cors := strings.Split(c.AllowedOrigins, ";")
	rs := []string{}
	for idx := range cors {
		itm := cors[idx]
		if strings.TrimSpace(itm) != "" {
			rs = append(rs, itm)
		}
	}
	return rs
}

// Loader load config from reader into Viper
type Loader interface {
	Load(viper.Viper) (*viper.Viper, error)
}

// generateConfigFromViper generate config from viper data
func generateConfigFromViper(v *viper.Viper) Config {

	return Config{
		Port:        v.GetString("PORT"),
		BaseURL:     v.GetString("BASE_URL"),
		ServiceName: v.GetString("SERVICE_NAME"),
		Env:         v.GetString("ENV"),

		AllowedOrigins: v.GetString("ALLOWED_ORIGINS"),

		DBHost:          v.GetString("DB_HOST"),
		DBPort:          v.GetString("DB_PORT"),
		DBUser:          v.GetString("DB_USER"),
		DBName:          v.GetString("DB_NAME"),
		DBPass:          v.GetString("DB_PASS"),
		DBSSLMode:       v.GetString("DB_SSL_MODE"),
		JWTSecret:       v.GetString("JWT_SECRET"),
		Timezone:        v.GetString("TIMEZONE"),
		AccessTokenTTL:  v.GetInt64("ACCESS_TOKEN_TTL"),
		RefreshTokenTTL: v.GetInt64("REFRESH_TOKEN_TTL"),
		RabbitURL:       v.GetString("RABBIT_URL"),
		RabbitPort:      v.GetString("RABBIT_PORT"),
		RabbitUser:      v.GetString("RABBIT_USER"),
		RabbitPassword:  v.GetString("RABBIT_PASSWORD"),
		Version:         app.Version,
	}
}

// DefaultConfigLoaders is default loader list
func DefaultConfigLoaders() []Loader {
	loaders := []Loader{}
	fileLoader := NewFileLoader(".env", ".")
	loaders = append(loaders, fileLoader)
	loaders = append(loaders, NewENVLoader())

	return loaders
}

// LoadConfig load config from loader list
func LoadConfig(loaders []Loader) Config {
	v := viper.New()
	v.SetDefault("PORT", "8080")
	v.SetDefault("ENV", "local")

	for idx := range loaders {
		newV, err := loaders[idx].Load(*v)

		if err == nil {
			v = newV
		}
	}
	return generateConfigFromViper(v)
}

// GetTokenTTL generate next timeout time
func (c *Config) GetTokenTTL() time.Time {
	return time.Now().Add(time.Minute * time.Duration(c.AccessTokenTTL))
}

// GetRefreshTokenTTL generate next timeout time
func (c *Config) GetRefreshTokenTTL() time.Time {
	return time.Now().Add(time.Minute * time.Duration(c.RefreshTokenTTL))
}
