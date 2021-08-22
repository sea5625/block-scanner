package auth

import (
	"crypto/rand"
	"motherbear/backend/configuration"
	"motherbear/backend/constants"
	"motherbear/backend/isaacerror"
	"motherbear/backend/utility"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"motherbear/backend/db"
	"motherbear/backend/logger"
)

type LOGIN struct {
	ID       string `json:"id"`
	PASSWORD string `json:"password"`
}

type USER struct {
	USER string `json:"user"`
}

type JWTtoken struct {
	TOKEN          string `json:"token"`
	PASSWORDSTATUS int    `json:"passwordStatus,omitempty"`
}

type JWTTokenOnly struct {
	TOKEN string `json:"token"`
}

type JWTpayload struct {
	IAT        string
	EXP        string
	USER       string
	USERID     string
	USERTYPE   string
	PERMISSION []string
}

var userCommonAuthorizedAPIList = map[string][]string{
	constants.UsersAPIBaseURL:    {constants.HTTPMethodPUT, constants.HTTPMethodGET},
	constants.ChannelsAPIBaseURL: {constants.HTTPMethodGET},
	constants.NodesAPIBaseURL:    {constants.HTTPMethodGET},
	constants.BlockAPIBaseURL:    {constants.HTTPMethodGET},
	constants.TxAPIBaseURL:       {constants.HTTPMethodGET},
	constants.SymptomAPIBaseURL:  {constants.HTTPMethodGET},
	constants.PrometheusGETAPIURL:  {constants.HTTPMethodGET},
}

var thirdPartyAllowableAPIList = []string{
	constants.ChannelsGetListAPIURL,
	constants.NodesAPIBaseURL,
	constants.BlockAPIBaseURL,
	constants.TxAPIBaseURL,
}

var userTypeList = map[string]map[string][]string{
	constants.DBUserTypeCodeCommon:     userCommonAuthorizedAPIList,
	constants.DBUserTypeCodeThirdParty: map[string][]string{},
}

var APIListWithoutAuthorization = []string{
	constants.APIVersionURL + constants.AuthLoginAPIURL,
	constants.APIVersionURL + constants.ResourcesAPIBaseURL + "/" + constants.ResourcesIDLoginLogoImage}

var channelPermissionAPIList = []string{constants.ChannelsAPIBaseURL, constants.NodesAPIBaseURL, constants.BlockAPIBaseURL, constants.TxAPIBaseURL, constants.SymptomAPIBaseURL}

var jwtSecret []byte
var once sync.Once

func init() {
	once.Do(func() {
		jwtSecret = make([]byte, 32)
		rand.Read(jwtSecret)
	})
}

// login godoc
// @Tags auth
// @Summary user login
// @Description user login
// @Accept  json
// @Produce  json
// @Param body body auth.LOGIN true "login ID / PASSWORD"
// @Success 200 {object} auth.JWTtoken
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /auth/login [POST]
func Login(c *gin.Context) {
	var login LOGIN
	c.BindJSON(&login)
	const statusNormal = 3001 // Password status is normal.
	const statusChange = 3002 // Should change password.

	database := db.DBgorm()

	user := db.GetUserInfoByUserID(login.ID)
	if user.USER_ID == "" {
		// No user in DB.
		internalError := isaacerror.SysErrNoUserInDB.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorFailedUserLogin, internalError)
		c.JSON(http.StatusBadRequest, message)
		return
	}

	// Password is wrong.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PASSWORD), []byte(login.PASSWORD)); err != nil {
		internalError := isaacerror.SysErrFailToValidPassword.Error()
		logger.Error(internalError, user.USER_ID)
		message := isaacerror.GetAPIError(isaacerror.ErrorFailedUserLogin, internalError)
		c.JSON(http.StatusBadRequest, message)
		return
	}

	// Generate JWT token.
	token, issueTime, err := generateToken(*user)
	if err != nil {
		message := isaacerror.GetAPIError(isaacerror.ErrorFailedUserLogin, err.Error())
		c.JSON(http.StatusBadRequest, message)
		return
	}

	// Put the latest date of current user.
	database.Model(&user).Update("LATEST_LOGIN_DATE", issueTime)

	var jwtToken JWTtoken
	// If user didn't set password before, should change their password.
	if user.TIME_WHEN_USER_MODIFIED_PASSWORD.IsZero() {
		jwtToken.PASSWORDSTATUS = statusChange
	} else {
		jwtToken.PASSWORDSTATUS = statusNormal
	}

	jwtToken.TOKEN = token

	c.JSON(http.StatusOK, jwtToken)
	logger.Info("Login successful.", user.USER_ID)
}

// PutHandler godoc
// @Tags auth
// @Summary PUT handler of reissue token.
// @Description Reissue token.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param userID body auth.USER true "User ID"
// @Success 200 {object} auth.JWTTokenOnly
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /auth/token [put]
func Token(c *gin.Context) {

	// Get user ID in gin context.
	userID, exists := c.Get(constants.ContextKeyUserID)
	if !exists {
		// No user ID.
		internalError := isaacerror.SysErrNoUserID.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToReissueToken, internalError)
		c.JSON(http.StatusBadRequest, message)
		return
	}

	userIDString := userID.(string)

	user := db.GetUserInfoByUserID(userIDString)
	if user.USER_ID == "" {
		// No user in DB.
		internalError := isaacerror.SysErrNoUserInDB.Error()
		logger.Error(internalError, userIDString)
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToReissueToken, internalError)
		c.JSON(http.StatusBadRequest, message)
		return
	}

	// Generate JWT token.
	token, _, err := generateToken(*user)
	if err != nil {
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToReissueToken, err.Error())
		c.JSON(http.StatusBadRequest, message)
		return
	}

	var jwtToken JWTTokenOnly

	jwtToken.TOKEN = token

	c.JSON(http.StatusOK, jwtToken)
	logger.Info("Token reissued.", user.USER_ID)
}

// AuthMiddleware intercepts the requests, and check for a valid jwt token.
func AuthentificateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set application/json on response Content-Type.
		c.Header(constants.HTTPHeaderContentType, constants.HTTPContentTypeApplicationJson)

		// Content-Type in HTTP header can use only application/json when HTTP method is POST, PUT.
		if c.Request.Method == constants.HTTPMethodPOST || c.Request.Method == constants.HTTPMethodPUT {
			if c.Request.Header.Get(constants.HTTPHeaderContentType) != constants.HTTPContentTypeApplicationJson {
				internalError := isaacerror.SysErrUsedUnsupportedContentType.Error()
				logger.Error(internalError)
				message := isaacerror.GetAPIError(isaacerror.ErrorUsedUnsupportedContentType, internalError)
				c.Abort()
				c.JSON(http.StatusNotAcceptable, message)
				return
			}
		}

		path := c.Request.URL.Path
		// Not check jwt token when login API.
		if utility.IsExistValueInList(path, APIListWithoutAuthorization) {
			return
		}

		// Verify JWT token and get JWT payload.
		payload, err := parseToken(c)
		if err != nil {
			internalError := err.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorInvalidToken, internalError)
			c.Abort()
			c.JSON(http.StatusUnauthorized, message)
			return
		}

		// Logging Request
		logFormat := makeLoggerFormat(c)
		logger.Info(logFormat, payload.USERID)

		// Check expire token.
		err = checkExpiredToken(payload)
		if err != nil {
			internalError := err.Error()
			logger.Error(internalError, payload.USERID)
			message := isaacerror.GetAPIError(isaacerror.ErrorExpiredToken, internalError)
			c.Abort()
			c.JSON(http.StatusUnauthorized, message)
			return
		}

		// If API is '/api/v1/auth/token'(reissue token), can used all user and user ID will transfer to token handler.
		if constants.APIVersionURL+constants.AuthReissueTokenAPIURL == path {
			c.Set(constants.ContextKeyUserID, payload.USERID)
			return
		}

		// Check permission when userType is no "USER_ADMIN".
		if payload.USERTYPE != constants.DBUserTypeCodeAdmin {
			// Check user if can use API.
			err = checkAPIPermission(c, payload)
			if err != nil {
				internalError := err.Error()
				logger.Error(internalError, payload.USERID)
				message := isaacerror.GetAPIError(isaacerror.ErrorUnauthorizedUser, internalError)
				c.Abort()
				c.JSON(http.StatusUnauthorized, message)
				return
			}

			// Check channel permission if only "USER_COMMON".
			//  - Not check channel permission if "THIRD_PARTY" because only can use RESTFul API, GET /api/v1/channels.
			if payload.USERTYPE == constants.DBUserTypeCodeCommon {
				err = checkChannelPermission(c, payload)
				if err != nil {
					internalError := err.Error()
					logger.Error(internalError)
					message := isaacerror.GetAPIError(isaacerror.ErrorFailToGetThatUnauthorizedChannel, internalError)
					c.Abort()
					c.JSON(http.StatusForbidden, message)
					return
				}
			}
		}
	}
}

// makeLoggerFormat Logging Request
func makeLoggerFormat(c *gin.Context) string {
	logFormat := c.Request.Method + ", " + c.Request.RequestURI
	if strings.Contains(c.Request.RequestURI, "channel") {
		var id string
		if strings.Contains(c.Request.RequestURI, "nodes") {
			id = c.Query("channel")
		} else {
			id = c.Param(constants.RequestResourceID)
		}
		if id != "" {
			var channelTB *db.CONFIGURATION_CHANNEL_TB
			channelTB = db.GetConfigurationChannelInfo(id)
			logFormat = logFormat + ", " + channelTB.CHANNEL_NAME
		}
	} else if strings.Contains(c.Request.RequestURI, "users") {
		id := c.Param(constants.RequestResourceID)
		if id != "" {
			// Get one user resource in DB.
			var userTB db.USER_INFO_TB
			userTB = *db.GetUserInfoByPK(id)
			logFormat = logFormat + ", " + userTB.USER_ID
		}
	} else if strings.Contains(c.Request.RequestURI, "nodes") {
		id := c.Param(constants.RequestResourceID)
		if id != "" {
			var nodeTB *db.CONFIGURATION_NODE_TB
			nodeTB = db.GetConfigurationNodeInfoByNodePK(id)
			logFormat = logFormat + ", " + nodeTB.NODE_NAME
		}
	}
	return logFormat
}

func generateToken(user db.USER_INFO_TB) (string, time.Time, error) {
	// Set the session time in JWT token.
	iat := time.Now().UTC()
	iatString := iat.Format(time.RFC3339)

	var exp time.Time
	// If user didn't set password before, token expire time is 1minutes.
	if user.TIME_WHEN_USER_MODIFIED_PASSWORD.IsZero() {
		exp = iat.Add(time.Minute * time.Duration(constants.InitTokenExpirationTimeInMin))
	} else {
		exp = iat.Add(time.Minute * time.Duration(constants.DefaultTokenExpirationTimeInMin))
	}
	expString := exp.Format(time.RFC3339)

	// Get user permission in DB.
	userPermissionTB := db.GetUserPermissions(user.PK)
	var permissionList []string
	permissionList = make([]string, 0)
	for _, value := range userPermissionTB {
		if value.PERMISSION_CHECK {
			permissionList = append(permissionList, value.PERMISSION_ALIAS)
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		constants.JWTPayloadKeyIat:        iatString,
		constants.JWTPayloadKeyExp:        expString,
		constants.JWTPayloadKeyUser:       user.PK,
		constants.JWTPayloadKeyUserType:   user.TYPE_CODE,
		constants.JWTPayloadKeyPermission: permissionList,
	})

	// Sign and get the complete encoded token as a string using the secret.
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		internalError := isaacerror.SysErrFailedToGenerateSignedString
		logger.Error(internalError.Error(), user.USER_ID)
		return "", iat, internalError
	}

	return tokenString, iat, nil
}

// Verify JWT token and get JWT payload.
func parseToken(c *gin.Context) (JWTpayload, error) {
	var payload JWTpayload

	// Get the JWT token sting in Authorization header.
	var header = c.Request.Header.Get(constants.HTTPHeaderAuthorization)
	if len(header) <= 7 {
		return payload, isaacerror.SysErrInvalidToken
	}

	var authorizationJWTType = header[:7]
	if authorizationJWTType != constants.HTTPAuthorizationJWTType {
		return payload, isaacerror.SysErrInvalidToken
	}

	var tokenString = header[7:]

	// Parse and validate JWT token.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, isaacerror.SysErrFailToVerifyToken
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return jwtSecret, nil
	})
	// Failed parsing JWT.
	if err != nil {
		return payload, err
	}

	// Get JWT payload.
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		if claims[constants.JWTPayloadKeyExp] == nil || claims[constants.JWTPayloadKeyExp] == "" ||
			claims[constants.JWTPayloadKeyIat] == nil || claims[constants.JWTPayloadKeyIat] == "" ||
			claims[constants.JWTPayloadKeyUser] == nil || claims[constants.JWTPayloadKeyUser] == "" ||
			claims[constants.JWTPayloadKeyUserType] == nil || claims[constants.JWTPayloadKeyUserType] == "" ||
			claims[constants.JWTPayloadKeyPermission] == nil {
			return payload, isaacerror.SysErrInvalidJWTPayload
		}
	} else {
		return payload, isaacerror.SysErrInvalidToken
	}

	payload.IAT = claims[constants.JWTPayloadKeyIat].(string)
	payload.EXP = claims[constants.JWTPayloadKeyExp].(string)
	payload.USER = claims[constants.JWTPayloadKeyUser].(string)
	payload.USERTYPE = claims[constants.JWTPayloadKeyUserType].(string)

	// Convert interface to string slice
	permissionList := utility.ConvertInterfaceToStringSlice(claims[constants.JWTPayloadKeyPermission])
	if permissionList == nil {
		return payload, isaacerror.SysErrInvalidToken
	}
	payload.PERMISSION = permissionList

	// Make sure it exists user.
	userTB := db.GetUserInfoByPK(payload.USER)
	if userTB.USER_ID == "" {
		return payload, isaacerror.SysErrInvalidToken
	}
	payload.USERID = userTB.USER_ID

	// Check permission.
	if userTB.TYPE_CODE != payload.USERTYPE {
		return payload, isaacerror.SysErrInvalidToken
	}

	return payload, nil
}

// Check user if can use API.
func checkAPIPermission(c *gin.Context, payload JWTpayload) error {

	// Get user type.
	userType := payload.USERTYPE

	if userType == constants.DBUserTypeCodeAdmin {
		return nil
	}

	// Set authorized API list for third party user according to the set value.
	// Third party user be authorized only GET HTTP method regardless of API type.
	thirdPartyAuthorizedAPIList := map[string][]string{}
	for _, value := range configuration.Conf().Authorization.ThirdPartyUserAPI {
		apiPath := "/" + value
		if !utility.IsExistValueInList(apiPath, thirdPartyAllowableAPIList) {
			return isaacerror.SysErrUsedUnauthorizedAPI
		}
		thirdPartyAuthorizedAPIList[apiPath] = []string{constants.HTTPMethodGET}
	}
	userTypeList[constants.DBUserTypeCodeThirdParty] = thirdPartyAuthorizedAPIList

	// Get authorized API list.
	authorizedAPIList := userTypeList[userType]

	// Get API type.
	APIType, err := getAPIType(c)
	if err != nil {
		return err
	}

	// If using unauthorized API, error.
	authorizedMethodList, exists := authorizedAPIList[APIType]
	if !exists {
		return isaacerror.SysErrUsedUnauthorizedAPI
	}

	// If using unauthorized http method on APIType, error.
	if !utility.IsExistValueInList(c.Request.Method, authorizedMethodList) {
		return isaacerror.SysErrUsedUnauthorizedAPI
	}

	if userType == constants.DBUserTypeCodeCommon {
		switch APIType {
		case constants.UsersAPIBaseURL: // users API.
			// 1. Can be used only to myself.
			id := c.Param(constants.RequestResourceID)
			if id == "" {
				return isaacerror.SysErrInvalidParameter
			}
			if id != payload.USER {
				return isaacerror.SysErrUsedUnauthorizedAPI
			}
		case constants.NodesAPIBaseURL: // nodes API.
			// 1. Can be used to must have permission to use it.
			if !utility.IsExistValueInList(constants.DBUserPermissionNode, payload.PERMISSION) {
				return isaacerror.SysErrUsedUnauthorizedAPI
			}

			// 2. Node Get-List API can be used only when select channel.
			channelID := c.Query(constants.RequestParamChannel)
			if channelID == "" {
				return isaacerror.SysErrUsedUnauthorizedAPI
			}
		case constants.BlockAPIBaseURL, constants.TxAPIBaseURL: // blocks API or txs API.
			// 1. Blocks, txs Get-List and Get API can be used only when select channel.
			channelID := c.Query(constants.RequestParamChannel)
			if channelID == "" {
				return isaacerror.SysErrUsedUnauthorizedAPI
			}
		case constants.SymptomAPIBaseURL: // peer symptom API.
			// 1. Can be used to must have permission to use it.
			if !utility.IsExistValueInList(constants.DBUserPermissionMonitoringLog, payload.PERMISSION) {
				return isaacerror.SysErrUsedUnauthorizedAPI
			}
		}
	}

	return nil
}

func checkChannelPermission(c *gin.Context, payload JWTpayload) error {
	// Only check channel permission for 'Get' API and channel relevant API.
	if c.Request.Method != constants.HTTPMethodGET {
		return nil
	}

	// Get API type.
	APIType, err := getAPIType(c)
	if err != nil {
		return err
	}

	if !utility.IsExistValueInList(APIType, channelPermissionAPIList) {
		return nil
	}

	// Get channel list with permission.
	permissionChannelList := getPermissionChannelList(payload.USER)

	switch APIType {
	case constants.ChannelsAPIBaseURL: // channels API.
		channelID := c.Param(constants.RequestResourceID)
		if channelID == "" {
			c.Set(constants.ContextKeyPermissionChannelList, permissionChannelList)
		} else {
			if !utility.IsExistValueInList(channelID, permissionChannelList) {
				return isaacerror.SysErrFailToGetThatUnauthorizedChannel
			}
		}
	case constants.NodesAPIBaseURL: // nodes API.
		// Get channel ID in query
		channelID := c.Query(constants.RequestParamChannel)
		if channelID != "" {
			if !utility.IsExistValueInList(channelID, permissionChannelList) {
				return isaacerror.SysErrFailToGetThatUnauthorizedChannel
			}
		}
	case constants.BlockAPIBaseURL, constants.TxAPIBaseURL: // blocks API, txs API.
		// Get channel ID in query
		channelID := c.Query(constants.RequestParamChannel)
		if channelID != "" {
			if !utility.IsExistValueInList(channelID, permissionChannelList) {
				return isaacerror.SysErrFailToGetThatUnauthorizedChannel
			}
		}
	case constants.SymptomAPIBaseURL: // symptom API.
		c.Set(constants.ContextKeyPermissionChannelList, permissionChannelList)
	}

	return nil
}

// Get API type with the URL path.
func getAPIType(c *gin.Context) (string, error) {
	path := c.Request.URL.Path
	var APIType string

	// Get API type with the URL path.
	pathSplit := strings.Split(path, constants.URLPathSeparator)

	// Error if path length is 3.
	if len(pathSplit)-1 >= constants.ValidPathLength {
		APIType = "/" + pathSplit[3]
	} else {
		return "", isaacerror.SysErrInvalidURLPath
	}

	return APIType, nil
}

func checkExpiredToken(payload JWTpayload) error {
	// Check expire token.
	exp, err := time.Parse(time.RFC3339, payload.EXP)
	if err != nil {
		return isaacerror.SysErrInvalidExpireTime

	}
	now := time.Now().UTC()

	if now.After(exp) {
		return isaacerror.SysErrExpiredToken
	}

	return nil
}

func getPermissionChannelList(userPK string) []string {
	channelTB := db.GetUserPermissionChannels(userPK)

	permissionChannelList := make([]string, 0)

	for _, value := range channelTB {
		permissionChannelList = append(permissionChannelList, value.CHANNEL_PK)
	}

	return permissionChannelList
}
