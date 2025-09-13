package service

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"workshop-management/internal/dto"
	"workshop-management/internal/services/service"
	"workshop-management/pkg/logger"
	"workshop-management/pkg/messages"
	"workshop-management/pkg/response"
	"workshop-management/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HandlerService (Service = jasa), not service layer
type HandlerService struct {
	Service *service.SrvService
}

func NewServiceHandler(s *service.SrvService) *HandlerService {
	return &HandlerService{Service: s}
}

func (h *HandlerService) Create(ctx *gin.Context) {
	authData := utils.GetAuthData(ctx)
	userId := utils.InterfaceString(authData["user_id"])
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][HandlerService][Create]", logId)

	var req dto.AddService
	if err := ctx.BindJSON(&req); err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; BindJSON ERROR: %s;", logPrefix, err.Error()))

		res := response.Response(http.StatusBadRequest, messages.InvalidRequest, logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Request: %+v;", logPrefix, utils.JsonEncode(req)))

	data, err := h.Service.Create(userId, req)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.Create; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusCreated, "Add service successfully", logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Success: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusCreated, res)
}

func (h *HandlerService) Fetch(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][HandlerService][Fetch]", logId)

	//query parameters
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	orderBy := ctx.DefaultQuery("order_by", "updated_at")
	orderDir := ctx.DefaultQuery("order_direction", "desc")
	search := ctx.Query("search")

	data, totalData, err := h.Service.Fetch(page, limit, orderBy, orderDir, search)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Fetch; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.PaginationResponse(http.StatusOK, int(totalData), page, limit, logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusOK, res)
}

func (h *HandlerService) GetById(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][HandlerService][GetById]", logId)

	serviceId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	data, err := h.Service.GetById(serviceId)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; GetById; Error: %+v", logPrefix, err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res := response.Response(http.StatusNotFound, messages.NotFound, logId, nil)
			res.Error = "service not found"
			ctx.JSON(http.StatusNotFound, res)
			return
		}

		res := response.Response(http.StatusInternalServerError, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusOK, "success", logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusOK, res)
}

func (h *HandlerService) Update(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][HandlerService][Update]", logId)
	authData := utils.GetAuthData(ctx)
	userId := utils.InterfaceString(authData["user_id"])

	serviceId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	var req dto.UpdateService
	if err = ctx.BindJSON(&req); err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; BindJSON ERROR: %s;", logPrefix, err.Error()))
		res := response.Response(http.StatusBadRequest, messages.InvalidRequest, logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Request: %+v;", logPrefix, utils.JsonEncode(req)))

	rows, err := h.Service.Update(userId, serviceId, req)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.Update; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}
	if rows == 0 {
		res := response.Response(http.StatusNotFound, messages.MsgNotFound, logId, nil)
		res.Error = response.Errors{Code: http.StatusNotFound, Message: messages.NotFound}
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	res := response.Response(http.StatusOK, fmt.Sprintf("Service with ID: '%s' updated successfully", serviceId), logId, nil)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Service with ID: '%s' updated successfully; Data: %v", logPrefix, serviceId, utils.JsonEncode(req)))
	ctx.JSON(http.StatusOK, res)
	return
}

func (h *HandlerService) Delete(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][HandlerService][Delete]", logId)
	authData := utils.GetAuthData(ctx)
	userId := utils.InterfaceString(authData["user_id"])

	serviceId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	if err = h.Service.Delete(userId, serviceId); err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.Delete; Error: %+v", logPrefix, err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res := response.Response(http.StatusNotFound, messages.MsgNotFound, logId, nil)
			res.Error = response.Errors{Code: http.StatusNotFound, Message: messages.NotFound}
			ctx.JSON(http.StatusNotFound, res)
			return
		}

		res := response.Response(http.StatusInternalServerError, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusOK, fmt.Sprintf("Service with ID: '%s' deleted successfully", serviceId), logId, nil)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Service with ID: '%s' deleted successfully", logPrefix, serviceId))
	ctx.JSON(http.StatusOK, res)
	return
}
