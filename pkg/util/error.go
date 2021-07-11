package util

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
)

func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	switch arg := err.(type) {
	case model.Error:
		c.JSON(int(arg.Code), arg)

	case *model.Error:
		c.JSON(int(arg.Code), arg)

	case error:
		c.JSON(http.StatusInternalServerError, model.Error{
			Code:    http.StatusInternalServerError,
			Message: model.Message(arg.Error()),
		})
	}
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
