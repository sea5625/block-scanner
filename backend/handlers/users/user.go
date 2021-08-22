package users

import (
	"motherbear/backend/constants"
	"motherbear/backend/db"
	"motherbear/backend/isaacerror"
	"motherbear/backend/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Data UserInputData `json:"data"`
}
type ResponseList struct {
	Data  []UserOutputData `json:"data"`
	Total int              `json:"total" example:"1"`
}
type Response struct {
	Data UserOutputData `json:"data"`
}

type UserInputData struct {
	UserID             string   `json:"userId" example:"admin"`
	FirstName          string   `json:"firstName" example:"DooSik"`
	LastName           string   `json:"lastName" example:"Choi:"`
	Email              string   `json:"email" example:"amdin@admin.com"`
	PhoneNumber        string   `json:"phoneNumber" example:"010-1234-1234"`
	Channels           []string `json:"channels" example:"PKCH_0000000000000001,PKCH_0000000000000002"`
	PermissionToAccess []string `json:"permissionToAccess" example:"Node,MonitoringLog"`
	Password           string   `json:"password,omitempty" example:"password123"`
	NewPassword        string   `json:"newPassword,omitempty" example:"password!@#"`
	UserType           string   `json:"userType,omitempty" example:"THIRD_PARTY"`
}

type UserOutputData struct {
	ID                 string        `json:"id" example:"PKID_0000000000000001"`
	UserID             string        `json:"userId" example:"admin"`
	FirstName          string        `json:"firstName" example:"DooSik"`
	LastName           string        `json:"lastName" example:"Choi:"`
	Email              string        `json:"email" example:"amdin@admin.com"`
	PhoneNumber        string        `json:"phoneNumber" example:"010-1234-1234"`
	Channels           []ChannelData `json:"channels"`
	PermissionToAccess []string      `json:"permissionToAccess" example:"Node,MonitoringLog"`
	UserType           string        `json:"userType" example:"USER_ADMIN"`
}

type ChannelData struct {
	ID   string `json:"id" example:"PKCH_0000000000000001"`
	Name string `json:"name" example:"channel1"`
}

// GetHandler godoc
// @Tags Users
// @Summary GET handler of users
// @Description Get many users resource.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Success 200 {object} users.ResponseList "Result for get many users resource."
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /users [get]
func GetHandlerList(c *gin.Context) {
}

// GetHandler godoc
// @Tags Users
// @Summary GET handler of users
// @Description Get many user resources.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param id path string true "User ID(not userID) to be get"
// @Success 200 {object} users.ResponseList "Result for get one user resource."
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /users/{id} [get]
func GetHandler(c *gin.Context) {
	id := c.Param(constants.RequestResourceID)

	var response interface{}
	if id != "" {
		// Get one user resource in DB.
		userOutputData, err := getUserOutputData(id)
		if err != nil {
			if err == isaacerror.SysErrNoUserInDB {
				internalError := isaacerror.SysErrNoUserInDB.Error()
				logger.Error(internalError)
				message := isaacerror.GetAPIError(isaacerror.ErrorNoUserInDB, internalError)
				c.JSON(http.StatusInternalServerError, message)
				return
			} else {
				internalError := isaacerror.SysUnknownError.Error()
				logger.Error(internalError)
				message := isaacerror.GetAPIError(isaacerror.ErrorUnknown, internalError)
				c.JSON(http.StatusInternalServerError, message)
				return
			}
		}

		// Set response data for get user API.
		var responseStruct Response
		responseStruct.Data = userOutputData

		response = responseStruct
	} else {
		// Get many user resources in db.
		userOutputData := getUserOutputDataList()

		// Set response data for get user API.
		var responseDataList ResponseList
		responseDataList.Data = userOutputData
		responseDataList.Total = len(responseDataList.Data)

		response = responseDataList
		bytes := "bytes 0-" + strconv.Itoa(responseDataList.Total) + "/" + strconv.Itoa(responseDataList.Total)
		c.Header(constants.HTTPHeaderContentRange, bytes)
		c.Header(constants.HTTPHeaderXTotalCount, strconv.Itoa(responseDataList.Total))
	}

	c.JSON(http.StatusOK, response)

}

// PostHandler godoc
// @Tags Users
// @Summary POST handler of users
// @Description Create the user resource.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param body body users.Request true "The user data to be create"
// @Success 200 {object} users.Response "Result for get the created user resource"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /users [post]
func PostHandler(c *gin.Context) {
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

	// Check parameter.
	if data.Data.UserID == "" || data.Data.FirstName == "" || data.Data.LastName == "" || data.Data.Password == "" || data.Data.Email == "" || data.Data.PhoneNumber == "" {
		// Invalid Parameter.
		internalError := isaacerror.SysErrInvalidParameter.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorInvalidParameter, internalError)
		c.JSON(http.StatusBadRequest, message)
		return
	}

	// Check duplicate user name.
	duplicationUser := db.GetUserInfoByUserID(data.Data.UserID)
	if duplicationUser.USER_ID != "" {
		// Duplicated user name.
		internalError := isaacerror.SysErrDuplicatedUserID.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorDuplicatedUserID, internalError)
		c.JSON(http.StatusConflict, message)
		return
	}

	// Check if valid channel.
	if !db.IsExistChannelListInDB(data.Data.Channels) {
		// Selected channel that not exist in DB.
		internalError := isaacerror.SysErrSelectedChannelThatNotExistInDB.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorSelectedChannelThatNotExistInDB, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	// Check if valid permission.
	if !db.IsExistPermissionListInDB(data.Data.PermissionToAccess) {
		// Selected permission that not exist in DB.
		internalError := isaacerror.SysErrSelectedPermissionThatNotExistInDB.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorSelectedPermissionThatNotExistInDB, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	// Insert user in database.
	id := db.InsertUserInfo(data.Data.UserID, data.Data.FirstName, data.Data.LastName, data.Data.Password, data.Data.Email, data.Data.PhoneNumber, data.Data.UserType, data.Data.Channels, data.Data.PermissionToAccess)
	if id == "" {
		// Failed insert user.
		internalError := isaacerror.SysErrFailToInsertUser.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToInsertUser, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	// Get one user resource in DB.
	userOutputData, err := getUserOutputData(id)
	if err != nil {
		if err == isaacerror.SysErrNoUserInDB {
			internalError := isaacerror.SysErrNoUserInDB.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorNoUserInDB, internalError)
			c.JSON(http.StatusInternalServerError, message)
			return
		} else {
			internalError := isaacerror.SysUnknownError.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorUnknown, internalError)
			c.JSON(http.StatusInternalServerError, message)
			return
		}
	}

	// Convert node data to response data.
	var response Response
	response.Data = userOutputData
	c.JSON(http.StatusOK, response)
	logger.Infof("Successfully created the user, %s", response.Data.UserID)
}

// PutHandler godoc
// @Tags Users
// @Summary PUT handler of users
// @Description Update the user resource.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param id path string true "User ID(not userID) to be update"
// @Param body body users.Request true "The user data to be update"
// @Success 200 {object} users.ResponseList "Result for get one user resource."
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /users/{id} [put]
func PutHandler(c *gin.Context) {
	// Get parameter.
	id := c.Param(constants.RequestResourceID)
	if id == "" {
		// Invalid parameter.
		internalError := isaacerror.SysErrInvalidParameter.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorInvalidParameter, internalError)
		c.JSON(http.StatusBadRequest, message)
		return
	}

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

	// Check parameter.
	if data.Data.UserID == "" || data.Data.FirstName == "" || data.Data.LastName == "" || data.Data.Email == "" || data.Data.PhoneNumber == "" {
		// Invalid Parameter.
		internalError := isaacerror.SysErrInvalidParameter.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorInvalidParameter, internalError)
		c.JSON(http.StatusBadRequest, message)
		return
	}

	// Get user data.
	var userTB db.USER_INFO_TB
	userTB = *db.GetUserInfoByPK(id)
	if userTB.PK == "" {
		internalError := isaacerror.SysErrNoUserInDB.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorNoUserInDB, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	// Check duplicate user name.
	if userTB.USER_ID != data.Data.UserID {
		duplicationUser := db.GetUserInfoByUserID(data.Data.UserID)
		if duplicationUser.USER_ID != "" {
			// Duplicated user name.
			internalError := isaacerror.SysErrDuplicatedUserID.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorDuplicatedUserID, internalError)
			c.JSON(http.StatusConflict, message)
			return
		}
	}

	// Check if valid channel.
	if !db.IsExistChannelListInDB(data.Data.Channels) {
		// Selected channel that not exist in DB.
		internalError := isaacerror.SysErrSelectedChannelThatNotExistInDB.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorSelectedChannelThatNotExistInDB, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	// Check if valid permission.
	if !db.IsExistPermissionListInDB(data.Data.PermissionToAccess) {
		// Selected permission that not exist in DB.
		internalError := isaacerror.SysErrSelectedPermissionThatNotExistInDB.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorSelectedPermissionThatNotExistInDB, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	// Update user in database.
	err = db.UpdateUser(id, data.Data.UserID, data.Data.FirstName, data.Data.LastName, data.Data.NewPassword, data.Data.Email, data.Data.PhoneNumber, constants.DBUserTypeCodeCommon, data.Data.Channels, data.Data.PermissionToAccess)
	if err != nil {
		// Failed update user.
		internalError := isaacerror.SysErrFailToUpdateUser.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToUpdateUser, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	// Get one user resource in DB.
	userOutputData, err := getUserOutputData(id)
	if err != nil {
		if err == isaacerror.SysErrNoUserInDB {
			internalError := isaacerror.SysErrNoUserInDB.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorNoUserInDB, internalError)
			c.JSON(http.StatusInternalServerError, message)
			return
		} else {
			internalError := isaacerror.SysUnknownError.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorUnknown, internalError)
			c.JSON(http.StatusInternalServerError, message)
			return
		}
	}

	// Convert node data to response data.
	var response Response
	response.Data = userOutputData
	c.JSON(http.StatusOK, response)
	logger.Infof("Successfully updated the user, %s, %s", userTB.USER_ID, response.Data.UserID)
}

// DeleteHandler godoc
// @Tags Users
// @Summary DELETE handler of users
// @Description Delete the user resource.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param id path string true "User ID(not userID) to be delete"
// @Failure 400 {string} string "error"
// @Failure 401 {string} string "error"
// @Failure 404 {string} string "error"
// @Failure 500 {string} string "error"
// @Router /users/{id} [delete]
func DeleteHandler(c *gin.Context) {
	// Get parameter.
	id := c.Param(constants.RequestResourceID)
	if id == "" {
		// Invalid parameter.
		internalError := isaacerror.SysErrInvalidParameter.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorInvalidParameter, internalError)
		c.JSON(http.StatusBadRequest, message)
		return
	}

	// Get one user resource in DB.
	userOutputData, err := getUserOutputData(id)
	if err != nil {
		if err == isaacerror.SysErrNoUserInDB {
			internalError := isaacerror.SysErrNoUserInDB.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorNoUserInDB, internalError)
			c.JSON(http.StatusInternalServerError, message)
			return
		} else {
			internalError := isaacerror.SysUnknownError.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorUnknown, internalError)
			c.JSON(http.StatusInternalServerError, message)
			return
		}
	}

	err = db.DeleteUserInfo(id)
	if err != nil {
		internalError := isaacerror.SysErrFailToDeleteUser.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToDeleteUser, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	// Set response data for get user API.
	var response Response
	response.Data = userOutputData

	c.JSON(http.StatusOK, response)
	logger.Infof("Successfully deleted the user, %s", response.Data.UserID)
}

func getUserOutputDataList() []UserOutputData {
	var userOutputData []UserOutputData

	// Get many user resources in db.
	var userListTB []db.USER_INFO_TB
	userListTB = db.GetUserInfo()

	// Get channelUserMapping list
	var channelUserMappingTB []db.CHANNEL_USER_MAPPING_TB
	channelUserMappingTB = db.GetUserPermissionChannelsTable()

	// Get channels map that can be convert channelUserMapping pk to channelUserMapping name.
	channelIDToName := db.GetChannelPKToNameMap()

	// Get permission page in user.
	var permissionPageTB []db.PERMISSION_USER_MAPPING_TB
	permissionPageTB = db.GetUserPermissionTable()

	userOutputData = make([]UserOutputData, len(userListTB))
	for i, userTB := range userListTB {
		userOutputData[i] = makeUserOutputData(userTB, channelUserMappingTB, permissionPageTB, channelIDToName)
	}

	return userOutputData
}

func getUserOutputData(id string) (UserOutputData, error) {
	var userOutputData UserOutputData

	// Get one user resource in DB.
	var userTB db.USER_INFO_TB
	userTB = *db.GetUserInfoByPK(id)

	if userTB.PK == "" {
		// No user in DB.
		return userOutputData, isaacerror.SysErrNoUserInDB
	}

	// Get channelUserMapping list
	var channelUserMappingTB []db.CHANNEL_USER_MAPPING_TB
	channelUserMappingTB = db.GetUserPermissionChannels(id)

	// Get channels map that can be convert channelUserMapping pk to channelUserMapping name.
	channelIDToName := db.GetChannelPKToNameMap()

	// Get permission page in user.
	var permissionPageTB []db.PERMISSION_USER_MAPPING_TB
	permissionPageTB = db.GetUserPermissions(id)

	userOutputData = makeUserOutputData(userTB, channelUserMappingTB, permissionPageTB, channelIDToName)

	return userOutputData, nil
}

func makeUserOutputData(userTB db.USER_INFO_TB, channelMappingTB []db.CHANNEL_USER_MAPPING_TB, permissionMappingTB []db.PERMISSION_USER_MAPPING_TB, channelIDToName map[string]string) UserOutputData {
	var userOutputData UserOutputData

	userOutputData.ID = userTB.PK
	userOutputData.UserID = userTB.USER_ID
	userOutputData.FirstName = userTB.FIRST_NAME
	userOutputData.LastName = userTB.LAST_NAME
	userOutputData.Email = userTB.EMAIL_ADRES
	userOutputData.PhoneNumber = userTB.MOBILE_PHONE_NO
	userOutputData.Channels = make([]ChannelData, 0)
	userOutputData.UserType = userTB.TYPE_CODE
	for _, channelMappingData := range channelMappingTB {
		if userTB.PK == channelMappingData.PK {
			channelDataTemp := ChannelData{
				ID:   channelMappingData.CHANNEL_PK,
				Name: channelIDToName[channelMappingData.CHANNEL_PK],
			}

			userOutputData.Channels = append(userOutputData.Channels, channelDataTemp)
		}
	}
	var permissionPageList []string
	permissionPageList = make([]string, 0)
	for _, permissionPage := range permissionMappingTB {
		if userTB.PK == permissionPage.PK {
			if permissionPage.PERMISSION_CHECK {
				permissionPageList = append(permissionPageList, permissionPage.PERMISSION_ALIAS)
			}
		}
	}
	userOutputData.PermissionToAccess = permissionPageList

	return userOutputData
}
