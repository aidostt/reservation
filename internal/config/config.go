package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	defaultGRPCPort = "9090"
	authority       = "qrcode-generation-service"
	EnvLocal        = "local"
)

type (
	Config struct {
		Environment string
		Authority   string
		GRPC        GRPCConfig `mapstructure:"grpc"`
		Postgres    PostgresConfig
	}

	PostgresConfig struct {
		User     string
		Host     string
		Password string
		Port     string
		DBName   string
	}

	GRPCConfig struct {
		Host    string        `mapstructure:"host"`
		Port    string        `mapstructure:"port"`
		Timeout time.Duration `mapstructure:"timeout"`
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
	return viper.UnmarshalKey("grpc", &cfg.GRPC)
}

func setFromEnv(cfg *Config) {
	cfg.Postgres.User = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
	cfg.Postgres.Port = os.Getenv("POSTGRES_PORT")
	cfg.Postgres.DBName = os.Getenv("POSTGRES_DB")

	cfg.GRPC.Host = os.Getenv("GRPC_HOST")

	cfg.Environment = EnvLocal
	cfg.Authority = authority
	cfg.Environment = EnvLocal
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
	viper.SetDefault("grpc.port", defaultGRPCPort)
}
