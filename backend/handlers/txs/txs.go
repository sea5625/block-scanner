package txs

import (
	"github.com/gin-gonic/gin"
	"motherbear/backend/constants"
	"motherbear/backend/db"
	"motherbear/backend/isaacerror"
	"motherbear/backend/logger"
	"motherbear/backend/polarbear"
	"motherbear/backend/utility"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TxResponseList struct {
	Data  []TxResponse `json:"data"`
	Total int          `json:"total" example:"1" format:"int32"`
}

type TxResponse struct {
	TxHash      string `json:"txHash" example:"0xf6a9cfccbcb40a8fa2a6226c8087f01917b3f6ab0a44b809874b46b3348aabea" `
	Status      string `json:"status" example:"Success" `
	Timestamp   string `json:"timeStamp" example:"2006-01-02T15:04:05Z07:00"`
	From        string `json:"from"  example:"hx5d91dee6102ead2aca60256cf33ebf9aab102c82"`
	To          string `json:"to"  example:"cx54d95fee187faaea03cee908f50623c8381179d0" `
	BlockHeight int64  `json:"blockHeight"  example:"124" `
	Data        string `json:"data"  example:"{ 'method': 'make' }" `
}

var allowTransactionStatus = []string{"Success", "Failure"}

// GetHandlerList godoc
// @Tags Transactions
// @Summary GET handler of transactions
// @Description Get many transactions  resources.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param channel query string  true "Channel to query. Can be channel name or PK of channel."
// @Param limit query integer true "Identify the number of results returned from a result set."
// @Param offset query integer  true "Identify the starting point to return data from a result set."
// @Param status query string false "'status' be used to search transaction. The kind of 'status' is 'success' and 'failure'."
// @Param blockHeight query string false "'blockHeight' be used to search transaction belong blockHeight."
// @Param from query string false "'from' be used to search timestamp. Used with 'to'."
// @Param to query string false "'to' be used to search timestamp. Used with 'from'."
// @Param fromAddress query string false "'fromAddress' is send address of transaction. Be use to search."
// @Param toAddress query string false "'toAddress' is receive address of transaction. Be use to search."
// @Param data query string false "'data' be used to search specific phrases in data field of transaction."
// @Success 200 {object} txs.TxResponseList "Result for many tx resources"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /txs [get]
func GetHandlerList(c *gin.Context) {
	// Check the parameters.
	offset, limit, err := utility.GetOffsetListFromRequest(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	channelName := c.Query(constants.RequestQueryChannel)
	if channelName == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	// Get search items in http query.
	txSearch, err := getTxSearchItems(c)
	if err != nil {
		internalError := err.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorInvalidParameter, internalError)
		c.JSON(http.StatusBadRequest, message)
		return
	}

	// If user queried channel as PK, then get the real name of channel.
	channelName, err = db.ConvertPKCHtoChannelName(channelName)
	if err != nil {
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToQueryBlockList, err.Error())
		c.JSON(http.StatusBadRequest, message)
		return
	}
	logger.Infof("Blocks requested, limit:%d, offset:%d in %s", limit, offset, channelName)

	// Query Tx list.
	var txInCh []polarbear.Tx
	var count int64
	count, err = polarbear.QueryTxsInChannelBySearch(channelName, limit, offset, txSearch, &txInCh)
	if err != nil {
		internalError := isaacerror.SysErrFailToQueryTxsInChannel.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToQueryTxList, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	var resp TxResponseList
	resp.Data = make([]TxResponse, len(txInCh))
	resp.Total = int(count)

	for i := 0; i < len(txInCh); i++ {
		convertPbTxToTxResponse(&txInCh[i], &resp.Data[i])
	}

	// Put total information.
	bytes := "bytes 0-" + strconv.Itoa(resp.Total) + "/" + strconv.Itoa(resp.Total)
	c.Header(constants.HTTPHeaderContentRange, bytes)
	c.Header(constants.HTTPHeaderXTotalCount, strconv.Itoa(resp.Total))

	// Return body.
	c.JSON(http.StatusOK, resp)
}

func getTxSearchItems(c *gin.Context) (polarbear.TxSearch, error) {
	txSearch := polarbear.TxSearch{
		Status:      "",
		BlockHeight: -1,
		From:        time.Time{},
		To:          time.Time{},
		FromAddress: "",
		ToAddress:   "",
	}
	if status, exist := c.GetQuery(constants.RequestQueryStatus); exist {
		for _, value := range allowTransactionStatus {
			if strings.ToLower(status) == strings.ToLower(value) {
				txSearch.Status = value
			}
		}

		if txSearch.Status == "" {
			return txSearch, isaacerror.SysErrInvalidTransactionStatus
		}
	}

	if blockHeight, exist := c.GetQuery(constants.RequestQueryBlockHeight); exist {
		blockHeightInt, err := strconv.ParseInt(blockHeight, 10, 64)
		if err != nil {
			return txSearch, isaacerror.SysErrFailToParseStringToInt
		}
		txSearch.BlockHeight = blockHeightInt
	}

	if from, exist := c.GetQuery(constants.RequestQueryFrom); exist {
		fromTimer, err := time.Parse(time.RFC3339, from)
		if err != nil {
			return txSearch, isaacerror.SysErrFailToParseTimeStringToTimeObject
		}
		txSearch.From = fromTimer
	}

	if to, exist := c.GetQuery(constants.RequestQueryTo); exist {
		toTimer, err := time.Parse(time.RFC3339, to)
		if err != nil {
			return txSearch, isaacerror.SysErrFailToParseTimeStringToTimeObject
		}
		txSearch.To = toTimer
	}

	if fromAddress, exist := c.GetQuery(constants.RequestQueryFromAddress); exist {
		txSearch.FromAddress = fromAddress
	}

	if toAddress, exist := c.GetQuery(constants.RequestQueryToAddress); exist {
		txSearch.ToAddress = toAddress
	}

	if data, exist := c.GetQuery(constants.RequestQueryData); exist {
		txSearch.Data = data
	}

	return txSearch, nil
}

func convertPbTxToTxResponse(resp *polarbear.Tx, out *TxResponse) {
	out.TxHash = resp.TxHash
	out.Data = resp.Data
	out.From = resp.From
	out.To = resp.To
	out.Timestamp = resp.Timestamp.Format(time.RFC3339)
	out.Status = resp.Status
	out.BlockHeight = resp.BlockHeight
}

// GetHandler godoc
// @Tags Transactions
// @Summary GET handler of transactions
// @Description Get the one transaction resource by tx hash.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param channel query string  true "Channel to query. Can be channel name or PK of channel."
// @Param txhash path string true "Transaction hash"
// @Success 200 {object} txs.TxResponse "Result for get the one transaction resource"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /txs/{txhash} [get]
func GetHandler(c *gin.Context) {

	// Check the parameters.
	txHash := c.Param(constants.RequestParamTxHash)
	channelName := c.Query(constants.RequestQueryChannel)
	if channelName == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	// If user queried channel as PK, then get the real name of channel.
	var err error
	channelName, err = db.ConvertPKCHtoChannelName(channelName)
	if err != nil {
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToQueryTx, err.Error())
		c.JSON(http.StatusBadRequest, message)
		return
	}
	logger.Infof("Tx requested, %s in %s", txHash, channelName)

	var txInChannel polarbear.Tx
	err = polarbear.QueryTxInChannelByHash(channelName, txHash, &txInChannel)
	if err != nil {
		internalError := isaacerror.SysErrFailToQueryTxInChannel.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToQueryTx, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	var resp TxResponse
	convertPbTxToTxResponse(&txInChannel, &resp)
	// Return body.
	c.JSON(http.StatusOK, resp)

}
