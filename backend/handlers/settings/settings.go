package settings

import (
	"github.com/gin-gonic/gin"
	"motherbear/backend/configuration"
	"motherbear/backend/isaacerror"
	"motherbear/backend/logger"
	"net/http"
)

type Request struct {
	Data SettingsData `json:"data"`
}
type Response struct {
	Data SettingsData `json:"data"`
}

type SettingsData struct {
	SessionTimeout int `json:"sessionTimeout" example:"30" format:"int32"`
}

// GetHandlerList godoc
// @Tags Settings
// @Summary GET handler of configuration values in ISAAC is using.
// @Description Get configuration values in ISAAC is using.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Success 200 {object} settings.Response "Result for get configuration value"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /settings [get]
func GetHandler(c *gin.Context) {
	// Get sessionTimeout value.
	sessionTimeout := configuration.Conf().ETC.SessionTimeout

	// Set response data.
	var response Response
	response.Data.SessionTimeout = sessionTimeout

	c.JSON(http.StatusOK, response)
}

// PutHandler godoc
// @Tags Settings
// @Summary PUT handler of configuration values in ISAAC is using.
// @Description Update configuration values in ISAAC is using.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param body body settings.Request true "configuration value to be update"
// @Success 200 {object} settings.Response "Result for get the updated configuration value"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /settings [put]
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

	// Get sessionTimeout in request body.
	settingTimeOut := data.Data.SessionTimeout

	// Set configuration value.
	configuration.Conf().ETC.SessionTimeout = settingTimeOut

	// Update  the YAML node configuration file.
	configuration.ChangeConfigFile(configuration.GetFilePath(), configuration.Conf())

	// Set response data.
	var response Response
	response.Data.SessionTimeout = configuration.Conf().ETC.SessionTimeout

	c.JSON(http.StatusOK, response)
	logger.Infof("Successfully updated SessionTimeout setting, %d", response.Data.SessionTimeout)
}
