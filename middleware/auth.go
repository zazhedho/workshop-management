package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"workshop-management/internal/repository/auth"
	"workshop-management/pkg/logger"
	"workshop-management/pkg/response"
	"workshop-management/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Middleware struct to hold dependencies
type Middleware struct {
	BlacklistRepo auth.Blacklist
}

// NewMiddleware creates a new middleware with its dependencies
func NewMiddleware(blacklistRepo auth.Blacklist) *Middleware {
	return &Middleware{
		BlacklistRepo: blacklistRepo,
	}
}

// AuthMiddleware is the authentication middleware
func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			err       error
			logId     uuid.UUID
			logPrefix string
		)

		logId = utils.GenerateLogId(ctx)
		logPrefix = fmt.Sprintf("[%s][AuthMiddleware]", logId)

		tokenString, dataJWT, err := utils.JwtClaims(ctx)
		if err != nil {
			logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Invalid Token: %s; Error: %s;", logPrefix, tokenString, err.Error()))
			res := response.Response(http.StatusUnauthorized, utils.MsgFail, logId, nil)
			res.Error = err.Error()
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		logPrefix += fmt.Sprintf("[%s][%s]", utils.InterfaceString(dataJWT["jti"]), utils.InterfaceString(dataJWT["user_id"]))

		// Check if token is blacklisted
		_, err = m.BlacklistRepo.GetByToken(tokenString)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; blacklistRepo.GetByToken; Error: %+v", logPrefix, err))
			res := response.Response(http.StatusInternalServerError, utils.MsgFail, logId, nil)
			res.Error = err.Error()
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}

		//the token is valid but has been logged out
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Invalid Token: %s; Error: token is blacklisted;", logPrefix, tokenString))
			res := response.Response(http.StatusUnauthorized, utils.MsgFail, logId, nil)
			res.Error = "Please login and try again"
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		ctx.Set(utils.CtxKeyAuthData, dataJWT)
		ctx.Set("token", tokenString)

		ctx.Next()
	}
}
