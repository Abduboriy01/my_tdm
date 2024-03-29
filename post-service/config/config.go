package config

import (
	"os"

    "github.com/spf13/cast"
)

// Config ...
type Config struct {
    Environment       string // develop, staging, production
    PostgresHost      string
    PostgresPort      int
    PostgresDatabase  string
    PostgresUser      string
    PostgresPassword  string
    LogLevel          string
    RPCPort           string
    UserServiceHost   string
    UserServicePort   int

    KafkaHost string
    KafkaPort int
}

// Load loads environment  inflates Config
func Load() Config {
    c := Config{}

    c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

    c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "database-1.cqgpi523wr5p.ap-south-1.rds.amazonaws.com"))
    c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
    c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "postdb"))
    c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
    c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "newpassword"))

    
    c.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVISE_HOST", "user_service"))
    c.UserServicePort = cast.ToInt(getOrReturnDefault("USER_SERVICE_PORT", 8098))

    c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
    c.RPCPort = cast.ToString(getOrReturnDefault("RPC_PORT", ":8908"))

    c.KafkaHost = cast.ToString(getOrReturnDefault("KAFKA_HOST", "kafka"))
    c.KafkaPort = cast.ToInt(getOrReturnDefault("kAFKA_PORT", 9092))

    return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
    _, exists := os.LookupEnv(key)
    if exists {
        return os.Getenv(key)
    }

    return defaultValue
}
