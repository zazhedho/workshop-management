package middleware

import (
	"fmt"
	"net/http"
	"workshop-management/pkg/logger"
	"workshop-management/pkg/response"
	"workshop-management/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ErrorHandler(c *gin.Context, err any) {
	logId, _ := c.Value(utils.CtxKeyId).(uuid.UUID)
	logger.WriteLog(logger.LogLevelPanic, fmt.Sprintf("RECOVERY: %s; Error: %+v;", logId.String(), err))

	res := response.Response(http.StatusInternalServerError, fmt.Sprintf("%s (%s)", utils.MsgFail, logId.String()), logId, nil)
	c.AbortWithStatusJSON(http.StatusInternalServerError, res)
	return
}
