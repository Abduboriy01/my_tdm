package main

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/my_tdm/api-gateway/api"
	"github.com/my_tdm/api-gateway/config"
	"github.com/my_tdm/api-gateway/pkg/logger"
	"github.com/my_tdm/api-gateway/services"
	rds "github.com/my_tdm/api-gateway/storage/redis"
)

func main() {
	//	var redisRepo repo.RedisRepositoryStorage
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "api_gateway")

	psqlString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",

		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabse)

	_, err := gormadapter.NewAdapter("postgres", psqlString, true)
	if err != nil {
		log.Error("new adapter error", logger.Error(err))

	}

	// casbinEnforcer, err := casbin.NewEnforcer(cfg.CasbinConfigPath, db)
	// if err != nil {
	// 	log.Error("new enforcer error", logger.Error(err))
	// 	return
	// }

	// file casbin
	casbinEnforcer, err := casbin.NewEnforcer(cfg.CasbinConfigPath, "./config/policy_defenition.csv")
	if err != nil {
		log.Error("new enforcer error", logger.Error(err))
		return
	}

	err = casbinEnforcer.LoadPolicy()
	if err != nil {
		log.Error("new load policy error", logger.Error(err))
		return
	}

	pool := redis.Pool{
		MaxIdle: 80,

		MaxActive: 12000,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPost))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	redisRepo := rds.NewRedisRepo(&pool)
	casbinEnforcer.GetRoleManager().(*defaultrolemanager.RoleManager).AddMatchingFunc("keyMatch", util.KeyMatch)
	casbinEnforcer.GetRoleManager().(*defaultrolemanager.RoleManager).AddMatchingFunc("keyMatch3", util.KeyMatch3)

	serviceManager, err := services.NewServiceManager(&cfg)
	if err != nil {
		log.Error("gRPC dial error", logger.Error(err))
	}

	server := api.New(api.Option{
		Conf:           cfg,
		Logger:         log,
		Casbin:         casbinEnforcer,
		ServiceManager: serviceManager,
		RedisRepo:      redisRepo,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run http server", logger.Error(err))
		panic(err)
	}

}
