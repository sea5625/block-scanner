package constants

const ServerPort = ":6553"

const URLPathSeparator = "/"
const APIVersionURL = "/api/v1"

const HTTPMethodGET = "GET"
const HTTPMethodPOST = "POST"
const HTTPMethodPUT = "PUT"
const HTTPMethodDELETE = "DELETE"

// HTTP Header Key
const HTTPHeaderAuthorization = "Authorization"
const HTTPHeaderContentType = "Content-Type"
const HTTPHeaderContentRange = "Content-Range"
const HTTPHeaderXTotalCount = "X-Total-Count"
const HTTPHeaderContentLength = "Content-Length"

// HTTP Header content
const HTTPContentTypeApplicationJson = "application/json"
const HTTPAuthorizationJWTType = "Bearer "

// HTTP request data key.
const RequestResourceID = "id"
const RequestParamChannel = "channel"
const RequestParamBlockID = "blockid"
const RequestParamTxHash = "txhash"
const RequestQueryLimit = "limit"
const RequestQueryOffset = "offset"
const RequestQueryFrom = "from"
const RequestQueryTo = "to"
const RequestQueryChannel = "channel"
const RequestQueryStatus = "status"
const RequestQueryBlockHeight = "blockHeight"
const RequestQueryFromAddress = "fromAddress"
const RequestQueryToAddress = "toAddress"
const RequestQueryData = "data"

// Gin context data key.
const ContextKeyPermissionChannelList = "permissionChannelList"
const ContextKeyUserID = "userID"

const AdminID = "admin"
const AdminPK = "PKID_0000000000000000"

// Users API URL.
const UsersAPIBaseURL = "/users"
const UsersGetListAPIURL = UsersAPIBaseURL
const UsersGetAPIURL = UsersAPIBaseURL + "/:id"
const UsersPostAPIURL = UsersAPIBaseURL
const UsersPutAPIURL = UsersAPIBaseURL + "/:id"
const UsersDeleteAPIURL = UsersAPIBaseURL + "/:id"

// Channels API URL.
const ChannelsAPIBaseURL = "/channels"
const ChannelsGetListAPIURL = ChannelsAPIBaseURL
const ChannelsGetAPIURL = ChannelsAPIBaseURL + "/:id"
const ChannelsPostAPIURL = ChannelsAPIBaseURL
const ChannelsPutAPIURL = ChannelsAPIBaseURL + "/:id"
const ChannelsDeleteAPIURL = ChannelsAPIBaseURL + "/:id"

// Nodes API URL.
const NodesAPIBaseURL = "/nodes"
const NodesGetListAPIURL = NodesAPIBaseURL
const NodesGetAPIURL = NodesAPIBaseURL + "/:id"
const NodesPostAPIURL = NodesAPIBaseURL
const NodesPutAPIURL = NodesAPIBaseURL + "/:id"
const NodesDeleteAPIURL = NodesAPIBaseURL + "/:id"

// Alerting API URL.
const AlertingAPIBaseURL = "/alerting"
const AlertingGetListAPIURL = AlertingAPIBaseURL
const AlertingPutAPIURL = AlertingAPIBaseURL

// Settings API URL.
const SettingsAPIBaseURL = "/settings"
const SettingsGetListAPIURL = SettingsAPIBaseURL
const SettingsPutAPIURL = SettingsAPIBaseURL

// Auth API URL.
const AuthAPIBaseURL = "/auth"
const AuthLoginAPIURL = AuthAPIBaseURL + "/login"
const AuthReissueTokenAPIURL = AuthAPIBaseURL + "/token"

// Block API URL
const BlockAPIBaseURL = "/blocks"
const BlockGETListAPIURL = BlockAPIBaseURL
const BlockGETAPIURL = BlockAPIBaseURL + "/:blockid"

// TX API URL
const TxAPIBaseURL = "/txs"
const TxGETListAPIURL = TxAPIBaseURL
const TxGETAPIURL = TxAPIBaseURL + "/:txhash"

// Resources API URL
const ResourcesAPIBaseURL = "/resources"
const ResourcesGETAPIURL = ResourcesAPIBaseURL + "/:id"

// Symptom API URL
const SymptomAPIBaseURL = "/symptom"
const SymptomGETAPIURL = SymptomAPIBaseURL

// Prometheus API URL
const PrometheusAPIBaseURL = "/prometheus"
const PrometheusGETAPIURL = PrometheusAPIBaseURL

// File path.
const ConfigFolderName = "config"
const ConfigFileName = "configuration.yaml"
const DBFolderName = "data"

// DB.
const DBTypeMysql = "mysql"
const DBTypeSqlite3 = "sqlite3"
const DBUserTypeCodeAdmin = "USER_ADMIN"
const DBUserTypeCodeCommon = "USER_COMMON"
const DBUserTypeCodeThirdParty = "THIRD_PARTY"
const DBUserPermissionNode = "Node"
const DBUserPermissionMonitoringLog = "MonitoringLog"
const DBDefaultUnsyncBlockToleranceTime = 360
const DBDefaultSlowResponseTime = 5
const DBDefaultUnsyncBlockDifference = 100

// Logger
const LoggerServerUser = "Isaac Server"
const SlowResponse = "Slow response"
const UnsyncBlock = "Unsync block"

// Auth.
const InitTokenExpirationTimeInMin = 1
const DefaultTokenExpirationTimeInMin = 5
const ValidPathLength = 3
const JWTPayloadKeyIat = "iat"
const JWTPayloadKeyExp = "exp"
const JWTPayloadKeyUser = "user"
const JWTPayloadKeyUserType = "userType"
const JWTPayloadKeyPermission = "permission"

// Allow to get resource id.
const ResourcesIDLoginLogoImage = "loginLogoImage"

// Node Type
const NodeType1 = "loopchain"
const NodeType2 = "goloop"
const NodeUnknown = "unknownType"