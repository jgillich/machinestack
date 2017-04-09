package config

// Config is what defines the behaviour of the api
type Config struct {

	// Address is the socket to bind to
	Address string `mapstructure:"address"`

	LogLevel string `mapstructure:"log_level"`

	AllowOrigins []string `mapstructure:"allow_origins"`

	// TLSConfig holds various TLS related configurations
	TLSConfig *TLSConfig `mapstructure:"-"`

	JwtConfig *JwtConfig `mapstructure:"-"`

	SchedulerConfig *SchedulerConfig `mapstructure:"-"`

	DriverConfig *DriverConfig `mapstructure:"-"`

	PostgresConfig *PostgresConfig `mapstructure:"-"`
}

// TLSConfig holds various TLS related configurations
type TLSConfig struct {
	Enable bool
	Auto   bool
	Cert   string
	Key    string
}

type DriverConfig struct {
	// Enable specifies the name of drivers to enable
	Enable []string `mapstructure:"enable"`

	// Options provides arbitrary key-value configuration for internals,
	// like authentication and drivers. The format is:
	//
	//	namespace.option = value
	Options DriverOptions `mapstructure:"options"`
}

type DriverOptions map[string]string

type SchedulerConfig struct {
	Name string `mapstructure:"name"`
}

type PostgresConfig struct {
	Address  string `mapstructure:"address"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type JwtConfig struct {
	Secret string `mapstructure:"secret"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Address:      ":7842",
		TLSConfig:    &TLSConfig{},
		DriverConfig: DefaultDriverConfig(),
		LogLevel:     "DEBUG",
	}
}

func DefaultDriverConfig() *DriverConfig {
	return &DriverConfig{
		Enable:  []string{},
		Options: map[string]string{},
	}
}
