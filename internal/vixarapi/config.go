package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/keenywheels/backend/internal/pkg/client/vk"
	"github.com/keenywheels/backend/internal/vixarapi/delivery/http/cookie"
	srvcSearch "github.com/keenywheels/backend/internal/vixarapi/service/search"
	userSvc "github.com/keenywheels/backend/internal/vixarapi/service/user"
	"github.com/keenywheels/backend/pkg/notifier/smtp"
	"github.com/keenywheels/backend/pkg/redis"
	"github.com/spf13/viper"
)

// HttpConfig config for http server
type HttpConfig struct {
	Port            string        `mapstructure:"port"`
	Host            string        `mapstructure:"host"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// LoggerConfig struct for logger config
type LoggerConfig struct {
	LogLevel      string `mapstructure:"loglvl"`
	Mode          string `mapstructure:"mode"`
	Encoding      string `mapstructure:"encoding"`
	LogPath       string `mapstructure:"log_path"`
	MaxLogSize    int    `mapstructure:"max_log_size"`
	MaxLogBackups int    `mapstructure:"max_log_backups"`
	MaxLogAge     int    `mapstructure:"max_log_age"`
}

// CORSConfig for CORS
type CORSConfig struct {
	AllowOrigins     []string      `mapstructure:"allow_origins"`
	AllowMethods     []string      `mapstructure:"allow_methods"`
	AllowHeaders     []string      `mapstructure:"allow_headers"`
	AllowCredentials bool          `mapstructure:"allow_credentials"`
	MaxAge           time.Duration `mapstructure:"max_age"`
}

// ServiceConfig contains all configs which connected to services
type ServiceConfig struct {
	UserSvc userSvc.Config `mapstructure:"user"`
}

// AppConfig contains all configs which connected to main app
type AppConfig struct {
	HttpCfg         HttpConfig                 `mapstructure:"http"`
	LoggerCfg       LoggerConfig               `mapstructure:"logger"`
	CORSConfig      CORSConfig                 `mapstructure:"cors"`
	SchedulerConfig srvcSearch.SchedulerConfig `mapstructure:"scheduler"`
	VKConfig        vk.Config                  `mapstructure:"vk"`
	Service         ServiceConfig              `mapstructure:"service"`
	CookieConfig    cookie.Config              `mapstructure:"cookie"`
	SMTPConfig      smtp.Config                `mapstructure:"smtp"`
}

// PostgresConfig config for postgres
type PostgresConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	User         string        `mapstructure:"user"`
	Password     string        `mapstructure:"password"`
	DBName       string        `mapstructure:"dbname"`
	SSLMode      string        `mapstructure:"sslmode"`
	MaxPoolSize  int           `mapstructure:"max_pool_size"`
	ConnAttempts int           `mapstructure:"conn_attempts"`
	ConnTimeout  time.Duration `mapstructure:"conn_timeout"`
}

// DSN return dsn using PostgresConfig
func (cfg *PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
}

// Config global config, contains all configs
type Config struct {
	AppCfg      AppConfig      `mapstructure:"app"`
	PostgresCfg PostgresConfig `mapstructure:"postgres"`
	RedisCfg    redis.Config   `mapstructure:"redis"`
}

// LoadConfig function which reads config file and return Config instance
func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config error: %w", err)
	}

	// env support
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
