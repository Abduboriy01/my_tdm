package events

import (
	"github.com/my_tdm/post-service/config"
	pb "github.com/my_tdm/post-service/genproto"
	"github.com/my_tdm/post-service/pkg/logger"
	"github.com/my_tdm/post-service/storage"
)

type EventHandler struct {
	config  config.Config
	storage storage.IStorage
	log     logger.Logger
}

func NewEventHandlerFunc(config config.Config, storage storage.IStorage, log logger.Logger) *EventHandler {
	return &EventHandler{
		config:  config,
		storage: storage,
		log:     log,
	}
}

func (h *EventHandler) Handle(value []byte) error {
	var user pb.User
	err := user.Unmarshal(value)
	if err != nil {
		return err
	}
	err = h.storage.Post().CreatePostUser(&user)
	if err != nil {
		return err
	}

	return nil
}
