package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	defaultHTTPPort               = "8000"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1

	EnvLocal = "local"
)

type (
	Config struct {
		Environment string
		HTTP        HTTPConfig
		Postgres    PostgresConfig
	}

	PostgresConfig struct {
		URI string
	}

	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}
)

func Init(configsDir, envDir string) (*Config, error) {
	populateDefaults()
	loadEnvVariables(envDir)
	if err := parseConfigFile(configsDir, ""); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	return viper.UnmarshalKey("http", &cfg.HTTP)
}

func setFromEnv(cfg *Config) {
	cfg.Postgres.URI = os.Getenv("POSTGRES_URI")

	cfg.HTTP.Host = os.Getenv("HTTP_HOST")

	cfg.Environment = "development"
}

func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.SetConfigName(env)

	return viper.MergeInConfig()
}

func loadEnvVariables(envPath string) {
	err := godotenv.Load(envPath)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

}

func populateDefaults() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.max_header_megabytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.timeouts.read", defaultHTTPRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHTTPRWTimeout)
}
