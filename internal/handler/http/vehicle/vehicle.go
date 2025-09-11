package vehicle

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"workshop-management/internal/dto"
	"workshop-management/internal/services/vehicle"
	"workshop-management/pkg/logger"
	"workshop-management/pkg/response"
	"workshop-management/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HandlerVehicle struct {
	Service *vehicle.ServiceVehicle
}

func NewVehicleHandler(s *vehicle.ServiceVehicle) *HandlerVehicle {
	return &HandlerVehicle{Service: s}
}

func (h *HandlerVehicle) CreateVehicle(ctx *gin.Context) {
	var req dto.AddVehicle
	authData := utils.GetAuthData(ctx)
	userId := utils.InterfaceString(authData["user_id"])
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][HandlerVehicle][CreateVehicle]", logId)

	if err := ctx.BindJSON(&req); err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; BindJSON ERROR: %s;", logPrefix, err.Error()))

		res := response.Response(http.StatusBadRequest, utils.InvalidRequest, logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data, err := h.Service.CreateVehicle(userId, req)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.CreateVehicle; Error: %+v", logPrefix, err))
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Error: License plate already exists", logPrefix))
			res := response.Response(http.StatusBadRequest, utils.MsgExists, logId, nil)
			res.Error = response.Errors{Code: http.StatusBadRequest, Message: "license plate already exists"}
			ctx.JSON(http.StatusBadRequest, res)
			return
		}

		res := response.Response(http.StatusInternalServerError, utils.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusCreated, "Add vehicle successfully", logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Success: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusCreated, res)
}

func (h *HandlerVehicle) GetVehicle(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][HandlerVehicle][GetVehicle]", logId)

	vehicleId := ctx.Param("id")
	data, err := h.Service.GetVehicle(vehicleId)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; GetVehicle; Error: %+v", logPrefix, err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res := response.Response(http.StatusNotFound, utils.NotFound, logId, nil)
			res.Error = "vehicle not found"
			ctx.JSON(http.StatusNotFound, res)
			return
		}

		res := response.Response(http.StatusInternalServerError, utils.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusOK, "success", logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusOK, res)
}

func (h *HandlerVehicle) FetchVehicles(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][HandlerVehicle][FetchVehicles]", logId)
	authData := utils.GetAuthData(ctx)

	var userId string
	isCustomer := utils.InterfaceString(authData["role"]) == utils.RoleCustomer
	if isCustomer {
		userId = utils.InterfaceString(authData["user_id"])
	}

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

	vehicles, totalData, err := h.Service.FetchVehicles(page, limit, orderBy, orderDir, search, userId)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; FetchVehicles; Error: %+v", logPrefix, err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res := response.Response(http.StatusNotFound, utils.NotFound, logId, nil)
			res.Error = "List vehicle not found"
			ctx.JSON(http.StatusNotFound, res)
			return
		}

		res := response.Response(http.StatusInternalServerError, utils.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.PaginationResponse(http.StatusOK, int(totalData), page, limit, logId, vehicles)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(vehicles)))
	ctx.JSON(http.StatusOK, res)
}
