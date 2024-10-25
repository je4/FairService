package service

import "github.com/gin-gonic/gin"

func NewResultMessage(ctx *gin.Context, status int, err error) {
	er := HTTPResultMessage{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

type HTTPResultMessage struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
