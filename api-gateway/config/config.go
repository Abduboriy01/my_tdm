package config

import (
	"os"

	"github.com/spf13/cast"
)

// Config ...
type Config struct {
	Environment string // develop, staging, production

	UserServiceHost string
	UserServicePort int

	PostgresHost string
	PostgresPort int
	PostgresUser string
	PostgresPassword string
	PostgresDatabse string

	// context timeout in seconds
	CtxTimeout int
	RedisHost  string
	RedisPost  int

	LogLevel string
	HTTPPort string
	CasbinConfigPath string

	SigninKey string
}

// Load loads environment vars and inflates Config
func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "database-1.cqgpi523wr5p.ap-south-1.rds.amazonaws.com"))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresDatabse = cast.ToString(getOrReturnDefault("POSGTRES_DB", "postdb"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASWORD", "newpassword"))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":8080"))
	c.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVICE_HOST", "user_service"))
	c.UserServicePort = cast.ToInt(getOrReturnDefault("USER_SERVICE_PORT", 8098))

	c.RedisHost = cast.ToString(getOrReturnDefault("REDIS_HOST", "redisdb"))
	c.RedisPost = cast.ToInt(getOrReturnDefault("REDIS_PORT", 6379))
	c.CasbinConfigPath=cast.ToString(getOrReturnDefault("CASBIN_CONFIG_PATH","./config/rbac_model.conf"))
	c.SigninKey = cast.ToString(getOrReturnDefault("SIGNIN_KEY", "uplarzyvhdtxqulrutblcyoiyyhoidmvdhntrxadcldejovctqwkizkeiyhuzzjhknwabjvwfcdouffzvkjiearbhxchyoalqqlmfbbgwfmydzybqszykmdeqbvcsgqjcggkfaxhsxikoujnivvxagaiabvkcfoo"))
	c.CtxTimeout = cast.ToInt(getOrReturnDefault("CTX_TIMEOUT", 7))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
