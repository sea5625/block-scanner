package users

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
	"motherbear/backend/constants"
	"motherbear/backend/db"
	"motherbear/backend/handlers/auth"
	"motherbear/backend/isaacerror"
	"motherbear/backend/utility"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

const dbPath string = ":memory:"

// User data to be add.
var userTestData = []db.USER_INFO_TB{
	{ // Add user1.
		PK:                               "PKID_0000000000000001",
		USER_ID:                          "user1",
		FIRST_NAME:                       "user1First",
		LAST_NAME:                        "user1Last",
		PASSWORD:                         "user1pass",
		EMAIL_ADRES:                      "user1@user.com",
		MOBILE_PHONE_NO:                  "010-1234-1234",
		TYPE_CODE:                        constants.DBUserTypeCodeCommon,
		TIME_WHEN_USER_MODIFIED_PASSWORD: time.Now(),
		CREATE_DATE:                      time.Now(),
		LATEST_LOGIN_DATE:                time.Now(),
	},
	{ // Add user2.
		PK:                               "PKID_0000000000000002",
		USER_ID:                          "user2",
		FIRST_NAME:                       "user2First",
		LAST_NAME:                        "user2Last",
		PASSWORD:                         "user2pass",
		EMAIL_ADRES:                      "user2@user.com",
		MOBILE_PHONE_NO:                  "010-1234-1234",
		TYPE_CODE:                        constants.DBUserTypeCodeCommon,
		TIME_WHEN_USER_MODIFIED_PASSWORD: time.Now(),
		CREATE_DATE:                      time.Now(),
		LATEST_LOGIN_DATE:                time.Now(),
	},
}

// Should be set channel mapping data when add user.
var channelUserMappingTestData = []db.CHANNEL_USER_MAPPING_TB{
	{
		PK:         "PKID_0000000000000001",
		CHANNEL_PK: "PKCH_0000000000000001",
	},
	{
		PK:         "PKID_0000000000000001",
		CHANNEL_PK: "PKCH_0000000000000002",
	},
	{
		PK:         "PKID_0000000000000002",
		CHANNEL_PK: "PKCH_0000000000000001",
	},
	{
		PK:         "PKID_0000000000000002",
		CHANNEL_PK: "PKCH_0000000000000002",
	},
}

// Should be set permission mapping data when add user.
var permissionUserMappingTestData = []db.PERMISSION_USER_MAPPING_TB{
	{
		PK:               "PKID_0000000000000001",
		PERMISSION_ALIAS: constants.DBUserPermissionNode,
		PERMISSION_CHECK: true,
	},
	{
		PK:               "PKID_0000000000000001",
		PERMISSION_ALIAS: constants.DBUserPermissionMonitoringLog,
		PERMISSION_CHECK: true,
	},
	{
		PK:               "PKID_0000000000000002",
		PERMISSION_ALIAS: constants.DBUserPermissionNode,
		PERMISSION_CHECK: true,
	},
	{
		PK:               "PKID_0000000000000002",
		PERMISSION_ALIAS: constants.DBUserPermissionMonitoringLog,
		PERMISSION_CHECK: false,
	},
}

var userResponseList ResponseList // Correct users list.
var channelList []string          // Using channel list in test.
var permissionList []string       // Permission list in test.

// Test GET-List user API.
func TestGetHandler(t *testing.T) {
	Setup()

	// Test get user list.
	router := gin.Default()
	router.GET(constants.UsersGetListAPIURL, GetHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodGET, constants.UsersGetListAPIURL, nil)
	router.ServeHTTP(w, request)

	var resultList ResponseList
	_ = json.Unmarshal(w.Body.Bytes(), &resultList)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, strconv.Itoa(userResponseList.Total), w.Header().Get("X-Total-Count"))
	assert.Equal(t, true, checkRepsonseList(userResponseList, resultList))
}

// Test GET user API.
func TestGetHandler2(t *testing.T) {
	Setup()

	// Test get user.
	user := userResponseList.Data[0].ID
	router := gin.Default()
	router.GET(constants.UsersGetAPIURL, GetHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodGET, constants.UsersAPIBaseURL+"/"+user, nil)
	router.ServeHTTP(w, request)

	var result Response
	_ = json.Unmarshal(w.Body.Bytes(), &result)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, true, isIdenticalUserData(userResponseList.Data[0], result.Data))
}

// Test POST user API.
func TestPostHandler(t *testing.T) {
	Setup()

	// Set create user data.
	var addData UserInputData
	addData.UserID = "user10"
	addData.FirstName = "user10First"
	addData.LastName = "user10Last"
	addData.Password = "user10Pass"
	addData.Email = "user10@user.com"
	addData.PhoneNumber = "010-1234-1234"
	addData.Channels = channelList
	addData.PermissionToAccess = permissionList

	// Add user Data in correct response list.
	var userOutputData UserOutputData
	userOutputData.UserID = addData.UserID
	userOutputData.FirstName = addData.FirstName
	userOutputData.LastName = addData.LastName
	userOutputData.Email = addData.Email
	userOutputData.PhoneNumber = addData.PhoneNumber
	userOutputData.UserType = constants.DBUserTypeCodeCommon
	userOutputData.Channels = getChannelDataList(addData.Channels)
	userOutputData.PermissionToAccess = addData.PermissionToAccess

	userResponseList.Data = append(userResponseList.Data, userOutputData)
	userResponseList.Total = len(userResponseList.Data)

	// Set request data.
	var requestData Request
	requestData.Data = addData
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.POST(constants.UsersPostAPIURL, PostHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPOST, constants.UsersPostAPIURL, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	var result Response
	_ = json.Unmarshal(w.Body.Bytes(), &result)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, true, isIdenticalUserData(userResponseList.Data[userResponseList.Total-1], result.Data))
	assert.Equal(t, true, checkUserListInDB(userResponseList))
}

// Test to return error when duplicate user ID.
func TestPostHandler2(t *testing.T) {
	Setup()

	// Set create user data.
	var addData UserInputData
	addData.UserID = "user1"
	addData.FirstName = "user10First"
	addData.LastName = "user10Last"
	addData.Password = "user10Pass"
	addData.Email = "user10@user.com"
	addData.PhoneNumber = "010-1234-1234"
	addData.Channels = channelList
	addData.PermissionToAccess = permissionList

	// Set request data.
	var requestData Request
	requestData.Data = addData
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.POST(constants.UsersPostAPIURL, PostHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPOST, constants.UsersPostAPIURL, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	// Check duplicate user ID.
	var err isaacerror.APIError
	_ = json.Unmarshal(w.Body.Bytes(), &err)

	assert.Equal(t, 409, w.Code)
	assert.Equal(t, isaacerror.ErrorDuplicatedUserID, err.Errors[0].UserMessage)
}

// Test to return error when include invalid channel name.
func TestPostHandler3(t *testing.T) {
	Setup()

	// Set create user data.
	var addData UserInputData
	addData.UserID = "user10"
	addData.FirstName = "user10First"
	addData.LastName = "user10Last"
	addData.Password = "user10Pass"
	addData.Email = "user10@user.com"
	addData.PhoneNumber = "010-1234-1234"
	addData.Channels = []string{"test"}
	addData.PermissionToAccess = permissionList

	// Set request data.
	var requestData Request
	requestData.Data = addData
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.POST(constants.UsersPostAPIURL, PostHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPOST, constants.UsersPostAPIURL, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	// Check duplicate user ID.
	var err isaacerror.APIError
	_ = json.Unmarshal(w.Body.Bytes(), &err)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, isaacerror.ErrorSelectedChannelThatNotExistInDB, err.Errors[0].UserMessage)
}

// Test to return error when include invalid permission.
func TestPostHandler4(t *testing.T) {
	Setup()

	// Set create user data.
	var addData UserInputData
	addData.UserID = "user10"
	addData.FirstName = "user10First"
	addData.LastName = "user10Last"
	addData.Password = "user10Pass"
	addData.Email = "user10@user.com"
	addData.PhoneNumber = "010-1234-1234"
	addData.Channels = channelList
	addData.PermissionToAccess = []string{"test"}

	// Set request data.
	var requestData Request
	requestData.Data = addData
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.POST(constants.UsersPostAPIURL, PostHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPOST, constants.UsersPostAPIURL, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	// Check duplicate user ID.
	var err isaacerror.APIError
	_ = json.Unmarshal(w.Body.Bytes(), &err)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, isaacerror.ErrorSelectedPermissionThatNotExistInDB, err.Errors[0].UserMessage)
}

// Test to create third party user and login.
func TestPostHandler5(t *testing.T) {
	Setup()

	// Set create user data.
	var addData UserInputData
	addData.UserID = "user10"
	addData.FirstName = "user10First"
	addData.LastName = "user10Last"
	addData.Password = "user10Pass"
	addData.Email = "user10@user.com"
	addData.PhoneNumber = "010-1234-1234"
	addData.UserType = constants.DBUserTypeCodeThirdParty

	// Add user Data in correct response list.
	var userOutputData UserOutputData
	userOutputData.UserID = addData.UserID
	userOutputData.FirstName = addData.FirstName
	userOutputData.LastName = addData.LastName
	userOutputData.Email = addData.Email
	userOutputData.PhoneNumber = addData.PhoneNumber
	userOutputData.UserType = constants.DBUserTypeCodeThirdParty

	userResponseList.Data = append(userResponseList.Data, userOutputData)
	userResponseList.Total = len(userResponseList.Data)

	// Set request data.
	var requestData Request
	requestData.Data = addData
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.POST(constants.UsersPostAPIURL, PostHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPOST, constants.UsersPostAPIURL, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	var result Response
	_ = json.Unmarshal(w.Body.Bytes(), &result)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, true, isIdenticalUserData(userResponseList.Data[userResponseList.Total-1], result.Data))
	assert.Equal(t, true, checkUserListInDB(userResponseList))

	// Test to return passwordStatus is 3001 when login.
	var loginData auth.LOGIN

	loginData.ID = addData.UserID
	loginData.PASSWORD = addData.Password

	loginDataJson, _ := json.Marshal(loginData)

	router = gin.Default()
	router.POST(constants.AuthLoginAPIURL, auth.Login)

	w = httptest.NewRecorder()
	request, _ = http.NewRequest(constants.HTTPMethodPOST, constants.AuthLoginAPIURL, bytes.NewBuffer(loginDataJson))
	router.ServeHTTP(w, request)

	var token auth.JWTtoken
	json.Unmarshal(w.Body.Bytes(), &token)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 3001, token.PASSWORDSTATUS)

}

// Test Put user API.
func TestPutHandler(t *testing.T) {
	Setup()

	// Set user data to update.
	var updateIndex = 1
	var updateID = userResponseList.Data[updateIndex].ID

	var updateData UserInputData
	updateData.UserID = userResponseList.Data[updateIndex].UserID
	updateData.FirstName = "user10First"
	updateData.LastName = "user10Last"
	updateData.Email = "user10@user.com"
	updateData.PhoneNumber = "010-1234-1234"
	updateData.Channels = []string{channelList[0]}
	updateData.PermissionToAccess = []string{permissionList[0]}

	// update user Data in correct response list.
	userResponseList.Data[updateIndex].UserID = updateData.UserID
	userResponseList.Data[updateIndex].FirstName = updateData.FirstName
	userResponseList.Data[updateIndex].LastName = updateData.LastName
	userResponseList.Data[updateIndex].Email = updateData.Email
	userResponseList.Data[updateIndex].PhoneNumber = updateData.PhoneNumber
	userResponseList.Data[updateIndex].Channels = getChannelDataList(updateData.Channels)
	userResponseList.Data[updateIndex].PermissionToAccess = updateData.PermissionToAccess

	// Set request data.
	var requestData Request
	requestData.Data = updateData
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.PUT(constants.UsersPutAPIURL, PutHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPUT, constants.UsersAPIBaseURL+"/"+updateID, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	var result Response
	_ = json.Unmarshal(w.Body.Bytes(), &result)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, true, isIdenticalUserData(userResponseList.Data[updateIndex], result.Data))
	assert.Equal(t, true, checkUserListInDB(userResponseList))
}

// Test to change password using PUT user API.
func TestPutHandler2(t *testing.T) {
	Setup()

	// Set user data to update.
	var updateIndex = 1
	var updateID = userResponseList.Data[updateIndex].ID

	var updateData UserInputData
	updateData.UserID = "user10"
	updateData.FirstName = "user10First"
	updateData.LastName = "user10Last"
	updateData.Email = "user10@user.com"
	updateData.PhoneNumber = "010-1234-1234"
	updateData.Channels = []string{channelList[0]}
	updateData.PermissionToAccess = []string{permissionList[0]}

	// Set password.
	for _, value := range userTestData {
		if updateID == value.PK {
			updateData.Password = value.PASSWORD
			break
		}
	}
	updateData.NewPassword = "user10Pass"
	changedPasswordDate := time.Now()

	// Update user Data in correct response list.
	userResponseList.Data[updateIndex].UserID = updateData.UserID
	userResponseList.Data[updateIndex].FirstName = updateData.FirstName
	userResponseList.Data[updateIndex].LastName = updateData.LastName
	userResponseList.Data[updateIndex].Email = updateData.Email
	userResponseList.Data[updateIndex].PhoneNumber = updateData.PhoneNumber
	userResponseList.Data[updateIndex].Channels = getChannelDataList(updateData.Channels)
	userResponseList.Data[updateIndex].PermissionToAccess = updateData.PermissionToAccess

	// Set request data.
	var requestData Request
	requestData.Data = updateData
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.PUT(constants.UsersPutAPIURL, PutHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPUT, constants.UsersAPIBaseURL+"/"+updateID, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	var result Response
	_ = json.Unmarshal(w.Body.Bytes(), &result)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, true, isIdenticalUserData(userResponseList.Data[updateIndex], result.Data))
	assert.Equal(t, true, checkUserListInDB(userResponseList))

	// Check password.
	userTB := db.GetUserInfoByPK(updateID)
	assert.Equal(t, true, db.CheckPasswordHash(updateData.NewPassword, userTB.PASSWORD))
	assert.Equal(t, true, changedPasswordDate.Before(userTB.TIME_WHEN_USER_MODIFIED_PASSWORD))
}

// Test to return error when duplicate user ID.
func TestPutHandler3(t *testing.T) {
	Setup()

	// Set user data to update.
	var updateIndex = 1
	var updateID = userResponseList.Data[updateIndex].ID

	var updateData UserInputData
	updateData.UserID = "admin"
	updateData.FirstName = "user10First"
	updateData.LastName = "user10Last"
	updateData.Email = "user10@user.com"
	updateData.PhoneNumber = "010-1234-1234"
	updateData.Channels = []string{channelList[0]}
	updateData.PermissionToAccess = []string{permissionList[0]}

	// Set password.
	updateData.Password = "1234"
	updateData.NewPassword = "user10Pass"

	// Set request data.
	var requestData Request
	requestData.Data = updateData
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.PUT(constants.UsersPutAPIURL, PutHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPUT, constants.UsersAPIBaseURL+"/"+updateID, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	var err isaacerror.APIError
	_ = json.Unmarshal(w.Body.Bytes(), &err)

	assert.Equal(t, 409, w.Code)
	assert.Equal(t, isaacerror.ErrorDuplicatedUserID, err.Errors[0].UserMessage)
}

// Test to return error when include invalid channel name.
func TestPutHandler4(t *testing.T) {
	Setup()

	// Set user data to update.
	var updateIndex = 1
	var updateID = userResponseList.Data[updateIndex].ID

	var updateData UserInputData
	updateData.UserID = "user10"
	updateData.FirstName = "user10First"
	updateData.LastName = "user10Last"
	updateData.Email = "user10@user.com"
	updateData.PhoneNumber = "010-1234-1234"
	updateData.Channels = []string{"test"}
	updateData.PermissionToAccess = []string{permissionList[0]}

	// Set password.
	updateData.Password = "1234"
	updateData.NewPassword = "user10Pass"

	// Set request data.
	var requestData Request
	requestData.Data = updateData
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.PUT(constants.UsersPutAPIURL, PutHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPUT, constants.UsersAPIBaseURL+"/"+updateID, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	var err isaacerror.APIError
	_ = json.Unmarshal(w.Body.Bytes(), &err)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, isaacerror.ErrorSelectedChannelThatNotExistInDB, err.Errors[0].UserMessage)
}

// Test to return error when include invalid permission.
func TestPutHandler5(t *testing.T) {
	Setup()

	// Set user data to update.
	var updateIndex = 1
	var updateID = userResponseList.Data[updateIndex].ID

	var updateData UserInputData
	updateData.UserID = "user10"
	updateData.FirstName = "user10First"
	updateData.LastName = "user10Last"
	updateData.Email = "user10@user.com"
	updateData.PhoneNumber = "010-1234-1234"
	updateData.Channels = []string{channelList[0]}
	updateData.PermissionToAccess = []string{"test"}

	// Set password.
	updateData.Password = "1234"
	updateData.NewPassword = "user10Pass"

	// Set request data.
	var requestData Request
	requestData.Data = updateData
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.PUT(constants.UsersPutAPIURL, PutHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPUT, constants.UsersAPIBaseURL+"/"+updateID, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	var err isaacerror.APIError
	_ = json.Unmarshal(w.Body.Bytes(), &err)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, isaacerror.ErrorSelectedPermissionThatNotExistInDB, err.Errors[0].UserMessage)
}

func TestDeleteHandler(t *testing.T) {
	Setup()

	// Set user data to update.
	var deleteIndex = 1
	var deleteID = userResponseList.Data[deleteIndex].ID

	router := gin.Default()
	router.DELETE(constants.UsersDeleteAPIURL, DeleteHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodDELETE, constants.UsersAPIBaseURL+"/"+deleteID, nil)
	router.ServeHTTP(w, request)

	var result Response
	_ = json.Unmarshal(w.Body.Bytes(), &result)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, true, isIdenticalUserData(userResponseList.Data[deleteIndex], result.Data))
	assert.Equal(t, "", db.GetUserInfoByPK(deleteID).USER_ID)
	assert.Equal(t, 0, len(db.GetUserPermissionChannels(deleteID)))
	assert.Equal(t, 0, len(db.GetUserPermissions(deleteID)))
}

func checkUserListInDB(correctResponseList ResponseList) bool {
	var userOutputDataInDB []UserOutputData
	userOutputDataInDB = getUserOutputDataList()

	var newResponseList ResponseList
	newResponseList.Data = userOutputDataInDB
	newResponseList.Total = len(newResponseList.Data)

	if !checkRepsonseList(correctResponseList, newResponseList) {
		return false
	}

	return true
}

func checkRepsonseList(var1, var2 ResponseList) bool {
	if var1.Total != var2.Total {
		return false
	}
	if len(var1.Data) != len(var2.Data) {
		return false
	}
	for _, value1 := range var1.Data {
		check := false
		for _, value2 := range var2.Data {
			if value1.ID == value2.ID {
				if isIdenticalUserData(value1, value2) {
					check = true
					break
				}
			} else {
				if value1.UserID == value2.UserID {
					if isIdenticalUserData(value1, value2) {
						check = true
						break
					}
				}
			}
		}
		if !check {
			return false
		}
	}

	return true
}

func isIdenticalUserData(var1 UserOutputData, var2 UserOutputData) bool {
	if var1.UserID != var2.UserID {
		return false
	}
	if var1.FirstName != var2.FirstName {
		return false
	}
	if var1.LastName != var2.LastName {
		return false
	}
	if var1.Email != var2.Email {
		return false
	}
	if var1.PhoneNumber != var2.PhoneNumber {
		return false
	}
	if var1.UserType != var2.UserType {
		return false
	}
	if !isIdenticalChannelMappingData(var1.Channels, var2.Channels) {
		return false
	}
	if !isIdenticalPermissionMappingData(var1.PermissionToAccess, var2.PermissionToAccess) {
		return false
	}

	return true
}

func isIdenticalChannelMappingData(var1, var2 []ChannelData) bool {
	if len(var1) != len(var2) {
		return false
	}

	for _, value1 := range var1 {
		check := false
		for _, value2 := range var2 {
			if value1.ID == value2.ID {
				if value1.Name == value2.Name {
					check = true
					break
				}
			}
		}
		if !check {
			return false
		}
	}

	return true
}

func isIdenticalPermissionMappingData(var1, var2 []string) bool {
	return utility.IsIdenticalSlice(var1, var2)
}

func getChannelDataList(channelIDList []string) []ChannelData {
	var channelData []ChannelData

	// Get channels map that can be convert channelUserMapping pk to channelUserMapping name.
	channelIDToName := db.GetChannelPKToNameMap()

	channelData = make([]ChannelData, len(channelIDList))
	for i, channelID := range channelIDList {
		channelData[i].ID = channelID
		channelData[i].Name = channelIDToName[channelID]
	}

	return channelData
}

func Setup() {
	// Create database.
	db.InitDB("sqlite3", dbPath)
	db.InitCreateTable()

	// Set test data.
	for _, value := range userTestData {
		hashedPassword, _ := db.HashPassword(value.PASSWORD)
		value.PASSWORD = hashedPassword
		db.DBgorm().Create(&value)
	}

	for _, value := range channelUserMappingTestData {
		db.DBgorm().Create(&value)

		if !utility.IsExistValueInList(value.CHANNEL_PK, channelList) {
			channelList = append(channelList, value.CHANNEL_PK)
		}
	}

	for i, value := range channelList {
		channelTB := &db.CONFIGURATION_CHANNEL_TB{
			CHANNEL_PK:   value,
			CHANNEL_NAME: "channel" + strconv.Itoa(i+1),
		}
		db.DBgorm().Create(&channelTB)
	}

	for _, value := range permissionUserMappingTestData {
		db.DBgorm().Create(&value)

		if !utility.IsExistValueInList(value.PERMISSION_ALIAS, permissionList) {
			permissionList = append(permissionList, value.PERMISSION_ALIAS)
		}
	}

	// Get many user resources in db.
	userResponseList.Data = getUserOutputDataList()
	userResponseList.Total = len(userResponseList.Data)
}
