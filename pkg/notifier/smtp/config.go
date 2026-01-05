package smtp

// Config contains config for SMTP
type Config struct {
	host     string `mapstructure:"host"`
	port     string `mapstructure:"port"`
	username string `mapstructure:"username"`
	password string `mapstructure:"password"`
}
