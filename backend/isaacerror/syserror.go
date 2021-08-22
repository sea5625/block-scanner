package isaacerror

import (
	"github.com/pkg/errors"
)

var (
	SysUnknownError                         = errors.New("Unknown Error.")
	SysErrInvalidParameter                  = errors.New("Invalid Parameter.")
	SysErrInvalidURLPath                    = errors.New("Invalid URL path.")
	SysErrUsedUnsupportedContentType        = errors.New("Used unsupported Content-Type.")
	SysErrFailToParseStringToInt            = errors.New("Fail to parse string to int.")
	SysErrFailToParseTimeStringToTimeObject = errors.New("Fail to parse time string to time object.")
	SysErrFailToGetFirstNodeIP              = errors.New("Fail to get first node IP.")
	SysErrNotExistChannelNameInTheGoloop    = errors.New("Not exist channel name in the goloop.")
	SysErrFailToParseURL                    = errors.New("Fail to Parse URL.")

	// Authentication error
	SysErrFailToValidPassword          = errors.New("Fail to valid password.")
	SysErrFailedToGenerateSignedString = errors.New("Failed to generate signed String.")
	SysErrInvalidToken                 = errors.New("Invalid Token.")
	SysErrFailedParsingJWT             = errors.New("Failed parsing JWT.")
	SysErrInvalidJWTPayload            = errors.New("Invalid JWT payload.")
	SysErrInvalidExpireTime            = errors.New("Invalid expire time.")
	SysErrExpiredToken                 = errors.New("Expired token.")
	SysErrUsedUnauthorizedAPI          = errors.New("Used unauthorized API.")
	SysErrFailToVerifyToken            = errors.New("Fail to verify JWT Token.")
	SysErrNoUserID                     = errors.New("No user ID.")

	// Nodes error
	SysErrNoNodeInDB         = errors.New("No node in DB.")
	SysErrFailToInsertNode   = errors.New("Fail to insert the node.")
	SysErrFailToUpdateNode   = errors.New("Fail to update the node.")
	SysErrFailToDeleteNode   = errors.New("Fail to delete the node.")
	SysErrDuplicatedNodeName = errors.New("Duplicated node name.")

	// Channels error.
	SysErrNoChannelInDB                    = errors.New("No channel in DB.")
	SysErrFailToInsertChannel              = errors.New("Fail to insert the channel.")
	SysErrFailToUpdateChannel              = errors.New("Fail to update the channel.")
	SysErrFailToDeleteChannel              = errors.New("Fail to delete the channel.")
	SysErrDuplicatedChannelName            = errors.New("Duplicated channel name.")
	SysErrSelectedNodeThatNotExistInDB     = errors.New("Selected node that not exist in DB.")
	SysErrFailToGetThatUnauthorizedChannel = errors.New("Fail to get that unauthorized channel.")

	// User error.
	SysErrNoUserInDB                         = errors.New("No user in DB.")
	SysErrFailToInsertUser                   = errors.New("Fail to insert the user.")
	SysErrFailToUpdateUser                   = errors.New("Fail to update the user.")
	SysErrFailToDeleteUser                   = errors.New("Fail to delete the user.")
	SysErrDuplicatedUserID                   = errors.New("Duplicated user ID.")
	SysErrSelectedChannelThatNotExistInDB    = errors.New("Selected channel that not exist in DB.")
	SysErrSelectedPermissionThatNotExistInDB = errors.New("Selected permission that not exist in DB.")

	//Prometheus error
	SysErrFailToConnectionPrometheus    = errors.New("Fail to connection prometheus server.")
	SysErrFailToReadBodyPrometheus      = errors.New("Fail to read response body received from prometheus server.")
	SysErrFailToReadBodyClosePrometheus = errors.New("Fail to close response body received from prometheus server.")
	SysErrFailToUnmarshalPrometheusData = errors.New("Fail to unmarshal prometheus data.")
	SysErrFailToGetPrometheusData       = errors.New("Fail to get prometheus data.")
	SysErrFailToGetPrometheusDataFromLC = errors.New("Fail to get prometheus data from the loopchain.")
	SysErrFailToGetPrometheusDataFromGL = errors.New("Fail to get prometheus data from the goloop.")
	SysErrFailToGetNewPrometheusData    = errors.New("Fail to get new prometheus data.")
	SysErrFailToGetNodeType             = errors.New("Fail to get node type.")
	SysErrNoChannelInISAAC              = errors.New("No channel in ISAAC.")

	// Tx error.
	SysErrInvalidTransactionStatus = errors.New("Invalid transaction status.")

	//Common error
	SysErrNoRequestBody               = errors.New("No request body.")
	SysErrFailConvertingStringToInt   = errors.New("Fail converting string to int.")
	SysErrFailConvertingStringToFloat = errors.New("Fail converting string to float")
	SysErrFailToLoadTimeLocation      = errors.New("Fail to load time location.")

	// PolarBear error
	SysErrFailToInitDBForPolarbear   = errors.New("Fail to initialize DB for Polarbear. ")
	SysErrFailToQueryBlocksInChannel = errors.New("Fail to query the blocks in channel from DB.  ")
	SysErrFailToQueryTxsInChannel    = errors.New("Fail to query the Txs in channel from DB.  ")
	SysErrFailToQueryBlockInChannel  = errors.New("Fail to query the block in channel from DB.  ")
	SysErrFailToQueryTxInChannel     = errors.New("Fail to query the Tx in channel from DB. ")
	SysErrFailToGetBlockData         = errors.New("Cannot get the block by height.")
	SysErrInvalidTimeSearchCondition = errors.New("Invalid time search condition, Both 'from/to' must be present.")

	// Data error.
	SysErrInvalidImageFileExtension      = errors.New("Invalid image file extension.")
	SysErrNotSupportedImageFileExtension = errors.New("Not supported image file extension.")
	SysErrFailToGetLoginLogoImage        = errors.New("Fail to get login logo image.")
	SysErrNotSupportedDataName           = errors.New("Not Supported Data name.")

	// Peer Symptom error.
	SysErrFailToQueryPeerSymptom = errors.New("Fail to query the symptom in peer from DB.  ")

	// goloop Admin API error.
	SysErrFailToConnectionNodeOfGoloop       = errors.New("Fail to connection node of goloop.")
	SysErrFailToReadBodyNodeOfGoloop         = errors.New("Fail to read response body received from node of goloop.")
	SysErrFailToUnmarshalChainsDataOfGoloop  = errors.New("Fail to unmarshal chains data of goloop.")
	SysErrFailToGetChainsDataOfGoloop        = errors.New("Fail to get chains data of goloop.")
	SysErrFailToUnmarshalSystemsDataOfGoloop = errors.New("Fail to unmarshal systems data of goloop.")
	SysErrFailToGetSystemDataOfGoloop        = errors.New("Fail to get system data of goloop.")
)
