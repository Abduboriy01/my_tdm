package main

import (
	"net"
	//	"string"

	"github.com/my_tdm/user-service-shu/config"
	"github.com/my_tdm/user-service-shu/events"
	pb "github.com/my_tdm/user-service-shu/genproto"
	"github.com/my_tdm/user-service-shu/pkg/db"
	"github.com/my_tdm/user-service-shu/pkg/logger"
	messegebroker "github.com/my_tdm/user-service-shu/pkg/messagebroker"
	"github.com/my_tdm/user-service-shu/service"
	grpcClient "github.com/my_tdm/user-service-shu/service/grpc_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "user-service")
	defer logger.Cleanup(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	grpcC, err := grpcClient.New(cfg)
	if err != nil {
		log.Error("error establishing grpc connection", logger.Error(err))
		return
	}

	// KAFKA
	publishersMap := make(map[string]messegebroker.Publisher)
	userTopicPublisher := events.NewKafkaProducerBroker(cfg, log, "user.user")
	defer func() {
		err := userTopicPublisher.Stop()
		if err != nil {
			log.Error("filed to stop kafka producer", logger.Error(err))
		}
	}()

	publishersMap["user"] = userTopicPublisher
	//KAFKA END

	userService := service.NewUserService(connDB, log, grpcC, publishersMap)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, userService)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
