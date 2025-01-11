package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var DEFAULT_SUCCESS_MSG = "OK"
var DEFAULT_ERROR_MSG = "ERROR"
var DEFAULT_ERROR_CODE = 500
var DEFAULT_SUCCESS_CODE = 200

func ResponseOK(c *gin.Context) {
	response := gin.H{
		"code":    DEFAULT_SUCCESS_CODE,
		"message": DEFAULT_SUCCESS_MSG,
	}
	c.JSON(http.StatusOK, response)
}

func ResponseOKWithData(c *gin.Context, data interface{}) {
	response := gin.H{
		"code":    DEFAULT_SUCCESS_CODE,
		"message": DEFAULT_SUCCESS_MSG,
	}
	if data != nil {
		response["data"] = data
	}

	c.JSON(http.StatusOK, response)
}

func ResponseFail(c *gin.Context) {
	response := gin.H{
		"code":    DEFAULT_ERROR_CODE,
		"message": DEFAULT_ERROR_MSG,
	}
	c.JSON(http.StatusOK, response)
}

func ResponseFailWithData(c *gin.Context, code int, msg string) {
	response := gin.H{
		"code":    DEFAULT_ERROR_CODE,
		"message": DEFAULT_ERROR_MSG,
	}
	if msg != "" {
		response["message"] = msg
	}
	if code != 0 {
		response["code"] = code
	}
	c.JSON(http.StatusOK, response)
}
