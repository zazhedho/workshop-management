package handler

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"workshop-management/internal/dto"
	"workshop-management/internal/services/users"
	"workshop-management/pkg/logger"
	"workshop-management/pkg/response"
	"workshop-management/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{Service: s}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body dto.UserRegister true "User registration details"
// @Success 201 {object} response.Success
// @Failure 400 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /user/register [post]
func (h *UserHandler) Register(ctx *gin.Context) {
	var (
		logId     uuid.UUID
		logPrefix string
		req       dto.UserRegister
	)

	logId = utils.GenerateLogId(ctx)
	logPrefix = fmt.Sprintf("[%s][UserHandler][Register]", logId)

	if err := ctx.BindJSON(&req); err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; BindJSON ERROR: %s;", logPrefix, err.Error()))

		res := response.Response(http.StatusBadRequest, utils.InvalidRequest, logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	user, err := h.Service.RegisterUser(req)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; userService.RegisterUser; Error: %+v", logPrefix, err))
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Error: email already exists", logPrefix))
			res := response.Response(http.StatusBadRequest, utils.MsgExists, logId, nil)
			res.Error = response.Errors{Code: http.StatusBadRequest, Message: "email already exists"}
			ctx.JSON(http.StatusBadRequest, res)
			return
		}

		res := response.Response(http.StatusInternalServerError, utils.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusCreated, "User registered successfully", logId, user)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Success: %+v;", logPrefix, utils.JsonEncode(user)))
	ctx.JSON(http.StatusCreated, res)
}

// Login godoc
// @Summary Login a user
// @Description Login a user
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body dto.Login true "User login details"
// @Success 200 {object} response.Success
// @Failure 400 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /user/login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var (
		logId     uuid.UUID
		logPrefix string
		req       dto.Login
	)

	logId = utils.GenerateLogId(ctx)
	logPrefix = fmt.Sprintf("[%s][UserController][Login]", logId)

	if err := ctx.BindJSON(&req); err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; BindJSON ERROR: %s;", logPrefix, err.Error()))

		res := response.Response(http.StatusBadRequest, utils.InvalidRequest, logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	token, err := h.Service.LoginUser(req, logId.String())
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.LoginUser; ERROR: %s;", logPrefix, err))
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == utils.ErrHashPassword {
			res := response.Response(http.StatusBadRequest, utils.InvalidCred, logId, nil)
			res.Error = response.Errors{Code: http.StatusBadRequest, Message: utils.MsgCredential}
			ctx.JSON(http.StatusBadRequest, res)
			return
		}

		res := response.Response(http.StatusInternalServerError, utils.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusOK, "success", logId, map[string]interface{}{"token": token})
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Success: %+v;", logPrefix, utils.JsonEncode(token)))
	ctx.JSON(http.StatusOK, res)
}

// Logout godoc
// @Summary Logout a user
// @Description Logout a user
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Success
// @Failure 500 {object} response.Error
// @Security ApiKeyAuth
// @Router /user/logout [post]
func (h *UserHandler) Logout(ctx *gin.Context) {
	var (
		logId     uuid.UUID
		logPrefix string
	)

	logId = utils.GenerateLogId(ctx)
	logPrefix = fmt.Sprintf("[%s][UserController][Logout]", logId)

	token, ok := ctx.Get("token")
	if !ok {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; token not found in context", logPrefix))
		res := response.Response(http.StatusInternalServerError, utils.MsgFail, logId, nil)
		res.Error = "token not found"
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	if err := h.Service.LogoutUser(token.(string)); err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; userService.LogoutUser; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, utils.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusOK, "User logged out successfully", logId, nil)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Success: User logged out successfully", logPrefix))
	ctx.JSON(http.StatusOK, res)
}
