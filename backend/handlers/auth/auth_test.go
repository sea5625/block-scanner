package auth

import (
	"bytes"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"motherbear/backend/constants"
	"motherbear/backend/db"
	"motherbear/backend/handlers/channels"
	"motherbear/backend/handlers/nodes"
	"motherbear/backend/handlers/users"
	"motherbear/backend/isaacerror"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const testID = "admin"
const testPASSWORD = "admin123"

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

func TestLogin(t *testing.T) {
	Setup()
	var loginData LOGIN

	loginData.ID = testID
	loginData.PASSWORD = testPASSWORD

	loginDataJson, _ := json.Marshal(loginData)

	router := gin.Default()
	router.POST(constants.AuthLoginAPIURL, Login)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPOST, constants.AuthLoginAPIURL, bytes.NewBuffer(loginDataJson))
	router.ServeHTTP(w, request)

	var token JWTtoken
	json.Unmarshal(w.Body.Bytes(), &token)

	// Test first login
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 3002, token.PASSWORDSTATUS)

	// Test login
	// Set latest change password
	var user db.USER_INFO_TB
	user.TIME_WHEN_USER_MODIFIED_PASSWORD = time.Now()
	db.DBgorm().Model(&user).Update(&user)

	w = httptest.NewRecorder()
	request, _ = http.NewRequest(constants.HTTPMethodPOST, constants.AuthLoginAPIURL, bytes.NewBuffer(loginDataJson))
	router.ServeHTTP(w, request)

	json.Unmarshal(w.Body.Bytes(), &token)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 3001, token.PASSWORDSTATUS)

	// Check user's logging date time is write or not.
	db.DBgorm().Where("USER_ID=?", testID).First(&user)
	assert.Equal(t, false, user.LATEST_LOGIN_DATE.IsZero())
}

func TestToken(t *testing.T) {
	Setup()

	iat := time.Now().UTC()
	exp := iat.Add(time.Minute * time.Duration(30))
	iatString := iat.Format(time.RFC3339)
	expString := exp.Format(time.RFC3339)

	// Get admin.
	userTB := db.GetUserInfoByUserID("admin")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		constants.JWTPayloadKeyIat:        iatString,
		constants.JWTPayloadKeyExp:        expString,
		constants.JWTPayloadKeyUser:       userTB.PK,
		constants.JWTPayloadKeyUserType:   constants.DBUserTypeCodeAdmin,
		constants.JWTPayloadKeyPermission: []string{constants.DBUserPermissionNode},
	})

	// Sign and get the complete encoded token as a string using the secret.
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Fatal("Fail to generate signed string.")
	}

	router := gin.Default()
	router.Use(AuthentificateMiddleware())

	router.PUT(constants.APIVersionURL+constants.AuthReissueTokenAPIURL, Token)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPUT, constants.APIVersionURL+constants.AuthReissueTokenAPIURL, nil)
	request.Header.Add(constants.HTTPHeaderAuthorization, constants.HTTPAuthorizationJWTType+tokenString)
	request.Header.Add(constants.HTTPHeaderContentType, constants.HTTPContentTypeApplicationJson)
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Test to permission about admin.
func TestAuthentificateMiddleware(t *testing.T) {
	Setup()

	iat := time.Now().UTC()
	exp := iat.Add(time.Minute * time.Duration(30))
	iatString := iat.Format(time.RFC3339)
	expString := exp.Format(time.RFC3339)

	// Get admin.
	userTB := db.GetUserInfoByUserID("admin")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		constants.JWTPayloadKeyIat:        iatString,
		constants.JWTPayloadKeyExp:        expString,
		constants.JWTPayloadKeyUser:       userTB.PK,
		constants.JWTPayloadKeyUserType:   constants.DBUserTypeCodeAdmin,
		constants.JWTPayloadKeyPermission: []string{constants.DBUserPermissionNode},
	})

	// Sign and get the complete encoded token as a string using the secret.
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Fatal("Fail to generate signed string.")
	}

	router := gin.Default()
	router.Use(AuthentificateMiddleware())

	router.GET(constants.APIVersionURL+constants.UsersGetListAPIURL, users.GetHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodGET, constants.APIVersionURL+constants.UsersGetListAPIURL, nil)
	request.Header.Add(constants.HTTPHeaderAuthorization, constants.HTTPAuthorizationJWTType+tokenString)
	request.Header.Add(constants.HTTPHeaderContentType, constants.HTTPContentTypeApplicationJson)
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Test to permission about user.
func TestAuthentificateMiddleware2(t *testing.T) {
	Setup()

	iat := time.Now().UTC()
	exp := iat.Add(time.Minute * time.Duration(30))
	iatString := iat.Format(time.RFC3339)
	expString := exp.Format(time.RFC3339)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		constants.JWTPayloadKeyIat:        iatString,
		constants.JWTPayloadKeyExp:        expString,
		constants.JWTPayloadKeyUser:       "PKID_0000000000000001",
		constants.JWTPayloadKeyUserType:   constants.DBUserTypeCodeCommon,
		constants.JWTPayloadKeyPermission: []string{constants.DBUserPermissionNode},
	})

	// Sign and get the complete encoded token as a string using the secret.
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Fatal("Fail to generate signed string.")
	}

	router := gin.Default()
	router.Use(AuthentificateMiddleware())

	router.GET(constants.APIVersionURL+constants.ChannelsGetListAPIURL, channels.GetHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodGET, constants.APIVersionURL+constants.ChannelsGetListAPIURL, nil)
	request.Header.Add(constants.HTTPHeaderAuthorization, constants.HTTPAuthorizationJWTType+tokenString)
	request.Header.Add(constants.HTTPHeaderContentType, constants.HTTPContentTypeApplicationJson)
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Test to return error when Used unauthorized API.
func TestAuthentificateMiddleware3(t *testing.T) {
	Setup()

	iat := time.Now().UTC()
	exp := iat.Add(time.Minute * time.Duration(30))
	iatString := iat.Format(time.RFC3339)
	expString := exp.Format(time.RFC3339)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		constants.JWTPayloadKeyIat:        iatString,
		constants.JWTPayloadKeyExp:        expString,
		constants.JWTPayloadKeyUser:       "PKID_0000000000000001",
		constants.JWTPayloadKeyUserType:   constants.DBUserTypeCodeCommon,
		constants.JWTPayloadKeyPermission: []string{constants.DBUserPermissionMonitoringLog},
	})

	// Sign and get the complete encoded token as a string using the secret.
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Fatal("Fail to generate signed string.")
	}

	router := gin.Default()
	router.Use(AuthentificateMiddleware())

	router.GET(constants.APIVersionURL+constants.NodesGetListAPIURL, nodes.GetHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodGET, constants.APIVersionURL+constants.NodesGetListAPIURL, nil)
	request.Header.Add(constants.HTTPHeaderAuthorization, constants.HTTPAuthorizationJWTType+tokenString)
	request.Header.Add(constants.HTTPHeaderContentType, constants.HTTPContentTypeApplicationJson)

	router.ServeHTTP(w, request)

	var errResponse isaacerror.APIError
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, isaacerror.ErrorUnauthorizedUser, errResponse.Errors[0].UserMessage)
}

// Test to return error when used expired token.
func TestAuthentificateMiddleware4(t *testing.T) {
	Setup()

	iat := time.Now().UTC()
	exp := iat
	iatString := iat.Format(time.RFC3339)
	expString := exp.Format(time.RFC3339)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		constants.JWTPayloadKeyIat:        iatString,
		constants.JWTPayloadKeyExp:        expString,
		constants.JWTPayloadKeyUser:       "PKID_0000000000000001",
		constants.JWTPayloadKeyUserType:   constants.DBUserTypeCodeCommon,
		constants.JWTPayloadKeyPermission: []string{constants.DBUserPermissionNode},
	})

	// Sign and get the complete encoded token as a string using the secret.
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Fatal("Fail to generate signed string.")
	}

	router := gin.Default()
	router.Use(AuthentificateMiddleware())

	router.GET(constants.APIVersionURL+constants.NodesGetListAPIURL, nodes.GetHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodGET, constants.APIVersionURL+constants.NodesGetListAPIURL, nil)
	request.Header.Add(constants.HTTPHeaderAuthorization, constants.HTTPAuthorizationJWTType+tokenString)
	request.Header.Add(constants.HTTPHeaderContentType, constants.HTTPContentTypeApplicationJson)

	router.ServeHTTP(w, request)

	var errResponse isaacerror.APIError
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, isaacerror.ErrorExpiredToken, errResponse.Errors[0].UserMessage)
}

// Test to return error when no authorization value.
func TestAuthentificateMiddleware5(t *testing.T) {
	Setup()

	router := gin.Default()
	router.Use(AuthentificateMiddleware())

	router.GET(constants.APIVersionURL+constants.NodesGetListAPIURL, nodes.GetHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodGET, constants.APIVersionURL+constants.NodesGetListAPIURL, nil)
	request.Header.Add(constants.HTTPHeaderContentType, constants.HTTPContentTypeApplicationJson)

	router.ServeHTTP(w, request)

	var errResponse isaacerror.APIError
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, isaacerror.ErrorInvalidToken, errResponse.Errors[0].UserMessage)
}

// Test to return error when used unauthorized channel.
func TestAuthentificateMiddleware6(t *testing.T) {
	Setup()

	iat := time.Now().UTC()
	exp := iat.Add(time.Minute * time.Duration(30))
	iatString := iat.Format(time.RFC3339)
	expString := exp.Format(time.RFC3339)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		constants.JWTPayloadKeyIat:        iatString,
		constants.JWTPayloadKeyExp:        expString,
		constants.JWTPayloadKeyUser:       "PKID_0000000000000002",
		constants.JWTPayloadKeyUserType:   constants.DBUserTypeCodeCommon,
		constants.JWTPayloadKeyPermission: []string{constants.DBUserPermissionNode},
	})

	// Sign and get the complete encoded token as a string using the secret.
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Fatal("Fail to generate signed string.")
	}

	channelID := "PKCH_0000000000000002"
	router := gin.Default()
	router.Use(AuthentificateMiddleware())

	router.GET(constants.APIVersionURL+constants.NodesGetListAPIURL, nodes.GetHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodGET, constants.APIVersionURL+constants.NodesGetListAPIURL+"?"+constants.RequestParamChannel+"="+channelID, nil)
	request.Header.Add(constants.HTTPHeaderAuthorization, constants.HTTPAuthorizationJWTType+tokenString)
	request.Header.Add(constants.HTTPHeaderContentType, constants.HTTPContentTypeApplicationJson)

	router.ServeHTTP(w, request)

	var errResponse isaacerror.APIError
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Equal(t, isaacerror.ErrorFailToGetThatUnauthorizedChannel, errResponse.Errors[0].UserMessage)
}

// Test to get permission channel list when call 'Get-List' API by common user.
func TestAuthentificateMiddleware7(t *testing.T) {
	Setup()

	iat := time.Now().UTC()
	exp := iat.Add(time.Minute * time.Duration(30))
	iatString := iat.Format(time.RFC3339)
	expString := exp.Format(time.RFC3339)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		constants.JWTPayloadKeyIat:        iatString,
		constants.JWTPayloadKeyExp:        expString,
		constants.JWTPayloadKeyUser:       "PKID_0000000000000002",
		constants.JWTPayloadKeyUserType:   constants.DBUserTypeCodeCommon,
		constants.JWTPayloadKeyPermission: []string{constants.DBUserPermissionNode},
	})

	// Sign and get the complete encoded token as a string using the secret.
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Fatal("Fail to generate signed string.")
	}

	// Set correct permission channel list.
	var correctChannelPermissionList []string
	for _, value := range channelUserMappingTestData {
		if value.PK == "PKID_0000000000000002" {
			correctChannelPermissionList = append(correctChannelPermissionList, value.CHANNEL_PK)
		}
	}
	router := gin.Default()
	router.Use(AuthentificateMiddleware())

	router.GET(constants.APIVersionURL+constants.ChannelsGetListAPIURL, func(c *gin.Context) {
		permissionChannelList, _ := c.Get(constants.ContextKeyPermissionChannelList)

		c.JSON(http.StatusOK, gin.H{
			"channelPermissionList": permissionChannelList,
		})
	})

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodGET, constants.APIVersionURL+constants.ChannelsGetListAPIURL, nil)
	request.Header.Add(constants.HTTPHeaderAuthorization, constants.HTTPAuthorizationJWTType+tokenString)
	request.Header.Add(constants.HTTPHeaderContentType, constants.HTTPContentTypeApplicationJson)

	router.ServeHTTP(w, request)

	response := struct {
		CHANNELPERMISSIONLIST []string `json:"channelPermissionList"`
	}{
		CHANNELPERMISSIONLIST: []string{},
	}
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, correctChannelPermissionList, response.CHANNELPERMISSIONLIST)
}

// Test to return error when Content-Type in HTTP header is not application/json.
func TestAuthentificateMiddleware8(t *testing.T) {
	Setup()

	router := gin.Default()
	router.Use(AuthentificateMiddleware())

	router.POST(constants.APIVersionURL+constants.NodesPostAPIURL, nodes.GetHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPOST, constants.APIVersionURL+constants.NodesPostAPIURL, nil)

	router.ServeHTTP(w, request)

	var errResponse isaacerror.APIError
	_ = json.Unmarshal(w.Body.Bytes(), &errResponse)

	assert.Equal(t, http.StatusNotAcceptable, w.Code)
	assert.Equal(t, isaacerror.ErrorUsedUnsupportedContentType, errResponse.Errors[0].UserMessage)
}

func Setup() {
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
	}

	for _, value := range permissionUserMappingTestData {
		db.DBgorm().Create(&value)
	}
}
