package controllers

import (
	"letschat/api/services"
	"letschat/errors"
	"letschat/infrastructure"
	"letschat/models"
	"letschat/responses"
	"letschat/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MessageController -> struct
type MessageController struct {
	logger         infrastructure.Logger
	env            infrastructure.Env
	messageService services.MessageService
}

// NewMessageController-> constructor
func NewMessageController(
	logger infrastructure.Logger,
	env infrastructure.Env,
	messageService services.MessageService,

) MessageController {
	return MessageController{
		logger:         logger,
		env:            env,
		messageService: messageService,
	}
}

func (cc MessageController) Create(c *gin.Context) {
	var message models.MessageM
	if err := c.ShouldBindJSON(&message); err != nil {
		cc.logger.Zap.Error("Error creating message: [ShouldBingJSON]:", err.Error())
		err = errors.BadRequest.Wrapf(err, "JSON Binding error")
		responses.HandleError(c, err)
		return
	}

	err := cc.messageService.Create(message)
	if err != nil {
		cc.logger.Zap.Error("Error creating message: ", err.Error())
		responses.HandleError(c, err)
		return
	}
	responses.JSON(c, http.StatusOK, "Message Created Succesfully")
}

func (cc MessageController) Delete(c *gin.Context) {
	id := c.Param("id")
	if len(id) == 0 {
		msg := "Please send roomId id "
		cc.logger.Zap.Error(msg)
		err := errors.BadRequest.New(msg)
		responses.HandleError(c, err)
		return
	}
	err := cc.messageService.Delete(id)
	if err != nil {
		cc.logger.Zap.Error("Error deleting message: ", err.Error())
		responses.HandleError(c, err)
		return
	}
	responses.JSON(c, http.StatusOK, "Message Deleted Succesfully")
}

func (mc MessageController) FindAll(c *gin.Context) {
	pagination := utils.BuildPagination(c)
	roomId := c.Param("roomid")
	if len(roomId) == 0 {
		msg := "Please send roomId id "
		mc.logger.Zap.Error(msg)
		err := errors.BadRequest.New(msg)
		responses.HandleError(c, err)
		return
	}
	sd, nextCursor, err := mc.messageService.FindAll(pagination, roomId)

	if err != nil {
		mc.logger.Zap.Error("Error finding all Message: ", err.Error())
		responses.HandleError(c, err)
		return
	}
	responses.JSONCursor(c, http.StatusOK, sd, nextCursor)

}
