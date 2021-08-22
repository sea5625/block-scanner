package blocks

import (
	"fmt"
	"motherbear/backend/constants"
	"motherbear/backend/db"
	"motherbear/backend/handlers/txs"
	"motherbear/backend/isaacerror"
	"motherbear/backend/logger"
	"motherbear/backend/polarbear"
	"motherbear/backend/utility"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// BlockResponseList is the response for GET LIST
type BlockResponseList struct {
	Data  []BlockResponse `json:"data"`
	Total int             `json:"total" example:"1" format:"int32"`
}

//BlockResponse is the response for GET
type BlockResponse struct {
	BlockHash   string           `json:"blockHash,omitempty" example:"0x586e5b26c51a8d07c6071a510ce7a26bf681342faa443180615c525a41934516"`
	Timestamp   string           `json:"timeStamp,omitempty" example:"2006-01-02T15:04:05Z07:00" `
	BlockHeight int64            `json:"blockHeight,omitempty" format:"int64"`
	PeerID      string           `json:"peerID,omitempty" example:"0x586e5b26c51a8d07c6071a510ce7a26bf681342faa443180615c525a41934516"`
	Signature   string           `json:"signature,omitempty" `
	ConfirmedTx []txs.TxResponse `json:"confirmedTx"`
}

func convertPolarbearBlockToBlockResponse(block *polarbear.Block, out *BlockResponse) {
	out.BlockHash = block.BlockHash
	out.BlockHeight = block.BlockHeight
	out.Signature = block.Signature
	out.PeerID = block.PeerID
	out.Timestamp = block.Timestamp.Format(time.RFC3339)

	// Copy Tx list.
	out.ConfirmedTx = make([]txs.TxResponse, len(block.Txs))

	for j := 0; j < len(block.Txs); j++ {
		out.ConfirmedTx[j].Status = block.Txs[j].Status
		out.ConfirmedTx[j].Timestamp = block.Txs[j].Timestamp.Format(time.RFC3339)
		out.ConfirmedTx[j].From = block.Txs[j].From
		out.ConfirmedTx[j].To = block.Txs[j].To
		out.ConfirmedTx[j].TxHash = block.Txs[j].TxHash
		out.ConfirmedTx[j].Data = block.Txs[j].Data
	}
	return
}

// GetHandlerList godoc
// @Tags Blocks
// @Summary GET handler of blocks
// @Description Get many block  resources.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param channel query string  true "Channel to query. Can be channel name or PK of channel."
// @Param limit query integer true "Identify the number of results returned from a result set."
// @Param offset query integer  true "Identify the starting point to return data from a result set."
// @Success 200 {object} blocks.BlockResponseList "Result for many block resources"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /blocks [get]
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

	// If user queried channel as PK, then get the real name of channel.
	channelName, err = db.ConvertPKCHtoChannelName(channelName)
	if err != nil {
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToQueryBlockList, err.Error())
		c.JSON(http.StatusBadRequest, message)
		return
	}
	logger.Infof("Blocks requested, limit:%d, offset:%d in %s", limit, offset, channelName)

	// Query data.
	var blocksInCh []polarbear.Block
	var count int64
	count, err = polarbear.QueryBlocksInChannel(channelName, limit, offset, &blocksInCh)
	if err != nil {
		internalError := isaacerror.SysErrFailToQueryBlocksInChannel.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToQueryBlockList, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	var resp BlockResponseList
	resp.Data = make([]BlockResponse, len(blocksInCh))
	resp.Total = int(count)

	for i := 0; i < len(blocksInCh); i++ {
		convertPolarbearBlockToBlockResponse(&blocksInCh[i], &resp.Data[i])
	}

	// Put total information.
	bytes := "bytes 0-" + strconv.Itoa(resp.Total) + "/" + strconv.Itoa(resp.Total)
	c.Header(constants.HTTPHeaderContentRange, bytes)
	c.Header(constants.HTTPHeaderXTotalCount, strconv.Itoa(resp.Total))

	// Return body.
	c.JSON(http.StatusOK, resp)
}

// GetHandler godoc
// @Tags Blocks
// @Summary GET handler of blocks
// @Description Get the one block resource by block hash.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param id path string true "Block ID. Can be block height (decimal) or block hash (hex)"
// @Param channel query string  true "Channel to query. Can be channel name or PK of channel."
// @Success 200 {object} blocks.BlockResponse "Result for a block resources"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /blocks/{id} [get]
func GetHandler(c *gin.Context) {

	// Check the parameters.
	blockID := c.Param(constants.RequestParamBlockID)
	channelName := c.Query(constants.RequestQueryChannel)
	if channelName == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	// If user queried channel as PK, then get the real name of channel.
	var err error
	channelName, err = db.ConvertPKCHtoChannelName(channelName)
	if err != nil {
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToQueryBlockList, err.Error())
		c.JSON(http.StatusBadRequest, message)
		return
	}

	// Check block ID is block hash or height.
	var blockInChannel polarbear.Block
	if strings.Contains(blockID, "0x") {
		blockHash := blockID
		logger.Info("Blocks requested with hash : %s in %s  ", blockHash, channelName)

		// Query data.
		err = polarbear.QueryBlockInChannelByHash(channelName, blockHash, &blockInChannel)
		if err != nil {
			internalError := isaacerror.SysErrFailToQueryBlocksInChannel.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorFailToQueryBlockList, internalError)
			c.JSON(http.StatusInternalServerError, message)
			return
		}

	} else {
		if blockHeight, err := strconv.ParseInt(blockID, 10, 64); blockHeight > 0 && err == nil {

			// Query data.
			err = polarbear.QueryBlockByHeightInChannel(channelName, blockHeight, &blockInChannel)
			if err != nil {
				internalError := isaacerror.SysErrFailToQueryBlocksInChannel.Error()
				logger.Error(internalError)
				message := isaacerror.GetAPIError(isaacerror.ErrorFailToQueryBlockList, internalError)
				c.JSON(http.StatusInternalServerError, message)
				return
			}

		} else {
			errMsg := fmt.Sprintf("%s is not right block height or block hash. ", blockID)
			logger.Error(errMsg)
			message := isaacerror.GetAPIError(isaacerror.ErrorFailToQueryBlockList, errMsg)
			c.JSON(http.StatusBadRequest, message)
			return
		}

	}

	var resp BlockResponse
	convertPolarbearBlockToBlockResponse(&blockInChannel, &resp)

	// Return body.
	c.JSON(http.StatusOK, resp)
}
