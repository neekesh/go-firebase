package services

import (
	"letschat/api/repository"
	"letschat/infrastructure"
	"letschat/models"
	"letschat/utils"
)

type MessageService struct {
	messagerepository repository.MessageRepository
	logger            infrastructure.Logger
}

func NewMessageService(
	messagerepository repository.MessageRepository,
	logger infrastructure.Logger,
) MessageService {
	return MessageService{
		messagerepository: messagerepository,
		logger:            logger,
	}
}

func (c MessageService) Create(message models.MessageM) error {
	return c.messagerepository.Create(message)
}

func (c MessageService) Delete(id string) error {
	return c.messagerepository.Delete(id)
}
func (c MessageService) FindAll(pagination utils.Pagination, roomId string) (*[]models.MessageM,string, error) {
	return c.messagerepository.FindAll(pagination, roomId)
}
