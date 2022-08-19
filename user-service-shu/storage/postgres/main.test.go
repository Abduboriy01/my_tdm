package postgres

import (
	"log"
	"os"
	"testing"

	"github.com/my_tdm/user-service-shu/config"
	"github.com/my_tdm/user-service-shu/pkg/db"
	"github.com/my_tdm/user-service-shu/pkg/logger"
	//	"github.com/my_tdm/user-service-shu/storage/repo"
)

var repo *userRepo

func TestMain(m *testing.M) {
	cfg := config.Load()

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	repo = NewUserRepo(connDB)

	os.Exit(m.Run())
}
