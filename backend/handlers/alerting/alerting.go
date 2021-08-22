package alerting

import (
	"github.com/gin-gonic/gin"
	"motherbear/backend/constants"
	"motherbear/backend/db"
	"motherbear/backend/isaacerror"
	"motherbear/backend/logger"
	"motherbear/backend/utility"
	"net/http"
	"strconv"
)

type Request struct {
	Data  []AlertingData `json:"data"`
	Total int            `json:"total" example:"1" format:"int32"`
}
type Response struct {
	Data AlertingData `json:"data"`
}

type ResponseList struct {
	Data  []AlertingData `json:"data"`
	Total int            `json:"total" example:"1" format:"int32"`
}

type AlertingData struct {
	ID                       string `json:"id" example:"PKCH_0000000000000001"`
	Name                     string `json:"name" example:"channel1"`
	UnsyncBlockToleranceTime int    `json:"unsyncBlockToleranceTime" example:"30" format:"int32"`
	SlowResponseTime         int    `json:"slowResponseTime" example:"5" format:"uint32"`
}

var validUnsyncBlockToleranceTime = []int{
	360, 420, 480, 540, 600,
}

var validSlowResponseTime = []int{
	2, 5, 8, 10, 20,
}

// GetHandlerList godoc
// @Tags Alerting
// @Summary GET handler of Alerting
// @Description Get alerting configuration value of channels.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Success 200 {object} alerting.ResponseList "Result for get alerting configuration value of channels"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /alerting [get]
func GetHandler(c *gin.Context) {
	// Get alerting configuration data and convert alerting configuration data to response data.
	var response ResponseList
	response = GetAlertingDataToResponse()

	bytes := "bytes 0-" + strconv.Itoa(response.Total) + "/" + strconv.Itoa(response.Total)
	c.Header(constants.HTTPHeaderContentRange, bytes)
	c.Header(constants.HTTPHeaderXTotalCount, strconv.Itoa(response.Total))

	c.JSON(http.StatusOK, response)
}

// PutHandler godoc
// @Tags Alerting
// @Summary PUT handler of Alerting
// @Description Update alerting configuration value of channels.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param body body alerting.Request true "alerting configuration value of channels to be update"
// @Success 200 {object} channels.Response "Result for get the updated alerting configuration value of channels"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /alerting [put]
func PutHandler(c *gin.Context) {
	// Get request body.
	var data Request
	err := c.ShouldBindJSON(&data)
	if err != nil {
		// No request body.
		internalError := isaacerror.SysErrNoRequestBody.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorInvalidParameter, internalError)
		c.JSON(http.StatusBadRequest, message)
		return
	}

	// Check valid value.
	for _, value := range data.Data {
		if utility.IsExistValueInList(value.UnsyncBlockToleranceTime, validUnsyncBlockToleranceTime) == false {
			internalError := isaacerror.SysErrNoRequestBody.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorInvalidParameter, internalError)
			c.JSON(http.StatusBadRequest, message)
			return
		}

		if utility.IsExistValueInList(value.SlowResponseTime, validSlowResponseTime) == false {
			internalError := isaacerror.SysErrNoRequestBody.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorInvalidParameter, internalError)
			c.JSON(http.StatusBadRequest, message)
			return
		}
	}

	// Get all alerting configuration data.
	var alertingTB []db.CONFIGURATION_DATA_ALERT_TB
	alertingTB = db.GetAlertConfigTable()

	// Check if the channel exists.
	for _, requestValue := range data.Data {
		check := false
		for _, DBValue := range alertingTB {
			if requestValue.ID == DBValue.CHANNEL_PK {
				check = true
				break
			}
		}

		if check == false {
			internalError := isaacerror.SysErrNoRequestBody.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorInvalidParameter, internalError)
			c.JSON(http.StatusBadRequest, message)
			return
		}
	}

	// Update alerting table in database.
	for _, requestValue := range data.Data {
		for _, DBValue := range alertingTB {
			if requestValue.ID == DBValue.CHANNEL_PK {
				unsyncBlockToleranceTime := requestValue.UnsyncBlockToleranceTime
				slowBlockResponseTime := requestValue.SlowResponseTime
				if unsyncBlockToleranceTime != DBValue.MAX_TIME_SEC_FOR_UNSYNC || slowBlockResponseTime != DBValue.MAX_TIME_SEC_FOR_RESPONSE {
					db.UpdateConfigurationAlertInfo(requestValue.ID, DBValue.ALERT_METHOD, unsyncBlockToleranceTime, slowBlockResponseTime)
					break
				}
			}
		}
	}

	// Get alerting configuration data and convert alerting configuration data to response data.
	var response ResponseList
	response = GetAlertingDataToResponse()

	bytes := "bytes 0-" + strconv.Itoa(response.Total) + "/" + strconv.Itoa(response.Total)
	c.Header(constants.HTTPHeaderContentRange, bytes)
	c.Header(constants.HTTPHeaderXTotalCount, strconv.Itoa(response.Total))

	c.JSON(http.StatusOK, response)
	logger.Infof("Successfully updated the setting of channel alert.")
}

func GetAlertingDataToResponse() ResponseList {
	// Get all alerting configuration data.
	var alertingTB []db.CONFIGURATION_DATA_ALERT_TB
	alertingTB = db.GetAlertConfigTable()
	alertingTBLen := len(alertingTB)

	// Get channels map that can be convert channel pk to channel name.
	channelPKToName := db.GetChannelPKToNameMap()

	// Convert the alerting configuration data to response data.
	var response ResponseList
	response.Data = make([]AlertingData, alertingTBLen)
	response.Total = alertingTBLen

	for i, value := range alertingTB {
		response.Data[i].ID = value.CHANNEL_PK
		response.Data[i].Name = channelPKToName[value.CHANNEL_PK]
		response.Data[i].UnsyncBlockToleranceTime = value.MAX_TIME_SEC_FOR_UNSYNC
		response.Data[i].SlowResponseTime = value.MAX_TIME_SEC_FOR_RESPONSE
	}

	return response
}
