package isaacerror

// APIError : Restful API Error Response Message list
type APIError struct {
	Errors []ErrorContent `json:"errors"`
}

// ErrorContent : Restful API Error Response Message
type ErrorContent struct {
	InternalMessage string `json:"internalMessage"`
	UserMessage     string `json:"userMessage"`
	MoreInfo        string `json:"moreInfo"`
}

const ErrorUnknown = "ErrorUnknown"
const ErrorInvalidToken = "ErrorInvalidToken"
const ErrorInvalidParameter = "ErrorInvalidParameter"
const ErrorFailedUserLogin = "ErrorFailedUserLogin"
const ErrorUnauthorizedUser = "ErrorUnauthorizedUser"
const ErrorFailToGetPrometheusData = "ErrorFailToGetPrometheusData"
const ErrorExpiredToken = "ErrorExpiredToken"
const ErrorUsedUnsupportedContentType = "ErrorUsedUnsupportedContentType"
const ErrorFailToReissueToken = "ErrorFailToReissueToken"

// nodes
const ErrorNoNodeInDB = "ErrorNoNodeInDB"
const ErrorFailToInsertNode = "ErrorFailToInsertNode"
const ErrorFailTodUpdateNode = "ErrorFailTodUpdateNode"
const ErrorFailToDeleteNode = "ErrorFailToDeleteNode"
const ErrorDuplicatedNodeName = "ErrorDuplicatedNodeName"

// channels
const ErrorNoChannelInDB = "ErrorNoChannelInDB"
const ErrorFailToInsertChannel = "ErrorFailToInsertChannel"
const ErrorFailToUpdateChannel = "ErrorFailToUpdateChannel"
const ErrorFailToDeleteChannel = "ErrorFailToDeleteChannel"
const ErrorDuplicatedChannelName = "ErrorDuplicatedChannelName"
const ErrorSelectedNodeThatNotExistInDB = "ErrorSelectedNodeThatNotExistInDB"
const ErrorFailToGetThatUnauthorizedChannel = "ErrorFailToGetThatUnauthorizedChannel"

// users
const ErrorNoUserInDB = "ErrorNoUserInDB"
const ErrorFailToInsertUser = "ErrorFailToInsertUser"
const ErrorFailToUpdateUser = "ErrorFailToUpdateUser"
const ErrorFailToDeleteUser = "ErrorFailToDeleteUser"
const ErrorDuplicatedUserID = "ErrorDuplicatedUserID"
const ErrorSelectedChannelThatNotExistInDB = "ErrorSelectedChannelThatNotExistInDB"
const ErrorSelectedPermissionThatNotExistInDB = "ErrorSelectedPermissionThatNotExistInDB"
const ErrorFailToValidPassword = "ErrorFailToValidPassword"

// Blocks
const ErrorFailToQueryBlockList = "ErrorFailToQueryBlockList"

//Txs
const ErrorFailToQueryTxList = "ErrorFailToQueryTxList"
const ErrorFailToQueryTx = "ErrorFailToQueryTx"

// Data
const ErrorNotSupportedDataName = "ErrorNotSupportedDataName"
const ErrorFailToGetLoginLogoImage = "ErrorFailToGetLoginLogoImage"

// Symptom
const ErrorFailToQueryPeerSymptom = "ErrorFailToQueryPeerSymptom"

// GetAPIError : Return RESTful API Error Response Message List
func GetAPIError(message string, internalMessage string) APIError {
	var errorcontent ErrorContent

	errorcontent.UserMessage = message
	errorcontent.InternalMessage = internalMessage
	errorcontent.MoreInfo = "http://localhost:6553/swagger/index.html"

	var apiError APIError
	apiError.Errors = make([]ErrorContent, 1)

	apiError.Errors[0] = errorcontent

	return apiError
}

// AddAPIErrorContent : Add RESTful API Error Response Message
func AddAPIErrorContent(apiError APIError, message string, internalMessage string) APIError {
	var errorcontent ErrorContent

	errorcontent.UserMessage = message
	errorcontent.InternalMessage = internalMessage
	errorcontent.MoreInfo = "http://localhost:6553/swagger/index.html"

	apiError.Errors = append(apiError.Errors, errorcontent)

	return apiError
}
