package processor

import (
	"fmt"
	"strings"
	"time"

	"github.com/keenywheels/backend/internal/pkg/client/llm"
	"github.com/keenywheels/backend/pkg/mailer/smtp"
	"github.com/spf13/viper"
)

// ClientsConfig struct for external clients config
type ClientsConfig struct {
	LLM llm.Config `mapstructure:"llm"`
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

// ProcessorConfig struct for processor config
type ProcessorConfig struct {
	WorkersCount int           `mapstructure:"workers_count"`
	MaxRetries   int           `mapstructure:"max_retries"`
	RetryDelay   time.Duration `mapstructure:"retry_delay"`
}

// PostgresConfig struct for postgres config
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

// DSN returns the dsn for postgres connection
func (pc PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		pc.User,
		pc.Password,
		pc.Host,
		pc.Port,
		pc.DBName,
		pc.SSLMode,
	)
}

// AppConfig contains application configuration
type AppConfig struct {
	Clients   ClientsConfig   `mapstructure:"clients"`
	Processor ProcessorConfig `mapstructure:"processor"`
	Postgres  PostgresConfig  `mapstructure:"postgres"`
	LoggerCfg LoggerConfig    `mapstructure:"logger"`
	SMTPCfg   smtp.Config     `mapstructure:"smtp"`
}

// KafkaTopics contains all kafka topics
type KafkaTopics struct {
	ScraperData   string `mapstructure:"scraper_data"`
	Notifications string `mapstructure:"notifications"`
}

// KafkaConfig contains Kafka configuration
type KafkaConfig struct {
	GroupID string      `mapstructure:"group_id"`
	Brokers []string    `mapstructure:"brokers"`
	Topics  KafkaTopics `mapstructure:"topics"`
}

// Config is the main configuration struct
type Config struct {
	App      AppConfig   `mapstructure:"app"`
	KafkaCfg KafkaConfig `mapstructure:"kafka"`
}

// LoadConfig parse yaml config into Config struct
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
