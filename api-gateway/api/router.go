package api

import (
	casbinN "github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/my_tdm/api-gateway/api/auth"
	"github.com/my_tdm/api-gateway/api/casbin"
	_ "github.com/my_tdm/api-gateway/api/docs" // swag
	v1 "github.com/my_tdm/api-gateway/api/handlers/v1"
	"github.com/my_tdm/api-gateway/config"
	"github.com/my_tdm/api-gateway/pkg/logger"
	"github.com/my_tdm/api-gateway/services"
	"github.com/my_tdm/api-gateway/storage/repo"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// New @BasePath /v1
// New ...
// @SecurityDefinitions.apikey BearerAuth
// @Description GetMyProfile
// @in header
// @name Authorization
type Option struct {
	Conf           config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	RedisRepo      repo.RedisRepositoryStorage
	Casbin         *casbinN.Enforcer
}

// New ...
func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	jwtHandler := &auth.JwtHandler{
		SigninKey: option.Conf.SigninKey,
		Log:       option.Logger,
	}

	router.Use(casbin.NewJwtRoleStruct(option.Casbin, option.Conf, *jwtHandler))

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
		Redis:          option.RedisRepo,
	})

	api := router.Group("/v1")
	api.POST("/users", handlerV1.CreateUser)
	api.GET("/users/:id", handlerV1.GetUser)
	api.GET("/users", handlerV1.ListUsers)
	api.POST("/users/register", handlerV1.RegisterUser)
	api.POST("/users/verfication", handlerV1.VerifyUser)
	api.GET("/users/login/:email/:password", handlerV1.Login)
	// api.PUT("/users/:id", handlerV1.UpdateUser)
	// api.DELETE("/users/:id", handlerV1.DeleteUser)

	url := ginSwagger.URL("swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
