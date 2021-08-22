package channels

import (
	"motherbear/backend/configuration"
	"motherbear/backend/constants"
	"motherbear/backend/db"
	"motherbear/backend/isaacerror"
	"motherbear/backend/logger"
	"motherbear/backend/prometheus"
	"motherbear/backend/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Data struct {
		Name  string   `json:"name" example:"channel1"`
		Nodes []string `json:"nodes" example:"node1,node2"`
	} `json:"data"`
}
type Response struct {
	Data ChannelData `json:"data"`
}

type ResponseList struct {
	Data  []ChannelData `json:"data"`
	Total int           `json:"total" example:"1" format:"int32"`
}

type ChannelData struct {
	ID                string     `json:"id" example:"PKCH_0000000000000001"`
	Name              string     `json:"name" example:"channel1"`
	GoloopChannelID   string     `json:"goloopChannelID" example:"0x822027"`
	Status            *int       `json:"status,omitempty" example:"0" format:"int32"`
	BlockHeight       *uint64    `json:"blockHeight,omitempty" example:"10" format:"uint64"`
	CountOfTX         *uint64    `json:"countOfTX,omitempty" example:"1000" format:"uint64"`
	ResponseTimeInSec *float64   `json:"responseTimeInSec,omitempty" example:"0.001" format:"float64"`
	Total             int        `json:"total" example:"1" format:"int32"`
	Nodes             []NodeData `json:"nodes"`
}

type NodeData struct {
	ID                   string   `json:"id" example:"PKND_0000000000000001"`
	Name                 string   `json:"name" example:"node1"`
	IP                   string   `json:"ip" example:"https://int-test-ctz.solidwallet.io"`
	BlockHeight          *uint64  `json:"blockHeight,omitempty" example:"10" format:"uint64"`
	CountOfTX            *uint64  `json:"countOfTX,omitempty" example:"1000" format:"uint64"`
	CountOfUnconfirmedTX *uint64  `json:"countOfUnconfirmedTX,omitempty" example:"0" format:"int32"`
	ResponseTimeInSec    *float64 `json:"responseTimeInSec,omitempty" example:"0.001" format:"float64"`
	IsLeader             *int     `json:"isLeader,omitempty" example:"0" format:"int32"`
	TimeStamp            string   `json:"timeStamp,omitempty" example:"2006-01-02T15:04:05Z07:00"`
	Status               *int     `json:"status,omitempty" example:"0" format:"int32"`
}

// GetHandlerList godoc
// @Tags Channels
// @Summary GET handler of channels
// @Description Get many channel resources.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Success 200 {object} channels.ResponseList "Result for get many channel resources"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /channels [get]
func GetHandlerList(c *gin.Context) {
}

// GetHandler godoc
// @Tags Channels
// @Summary GET handler of channels
// @Description Get the one channel resource.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param id path string true "Channel ID to be get"
// @Success 200 {object} channels.Response "Result for get the one channel resource"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /channels/{id} [get]
func GetHandler(c *gin.Context) {
	id := c.Param(constants.RequestResourceID)

	var response interface{}
	if id != "" {
		// Get one channel resource.
		var channelTB *db.CONFIGURATION_CHANNEL_TB
		channelTB = db.GetConfigurationChannelInfo(id)
		if channelTB.CHANNEL_PK == "" {
			// No channel in DB.
			internalError := isaacerror.SysErrNoChannelInDB.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorNoChannelInDB, internalError)
			c.JSON(http.StatusInternalServerError, message)
			return
		}

		// Get loopchainChannelData data.
		responseData, err := getResponseDataByChannelPK(id)
		if err != nil {
			// Fail to get prometheus data.
			internalError := isaacerror.SysErrFailToGetPrometheusData.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorFailToGetPrometheusData, internalError)
			c.JSON(http.StatusInternalServerError, message)
			return
		}
		response = responseData
	} else {

		// Get all channel resources.
		var channelTB []db.CONFIGURATION_CHANNEL_TB
		permissionChannelList, exists := c.Get(constants.ContextKeyPermissionChannelList)
		if exists {
			// If permission channel list received from middleware, Channel not in permission channel list is unauthorized channel.
			// If channel is not in permission channel list, exclude.
			permissionChannelListString := utility.ConvertInterfaceToStringSlice(permissionChannelList)
			channelTB = db.GetConfigurationChannelInfoByList(permissionChannelListString)
		} else {
			// If not permission channel list received from middleware, get all channel.
			channelTB = db.GetConfigurationChannelTable()
		}
		channelTBLen := len(channelTB)

		// Convert the channel data to response data.
		var responseDataList ResponseList
		responseDataList.Data = make([]ChannelData, channelTBLen)
		responseDataList.Total = channelTBLen

		for i, value := range channelTB {
			// Get prometheus data.
			responseData, err := getResponseDataByChannelPK(value.CHANNEL_PK)
			if err != nil {
				// Fail to get prometheus data.
				internalError := isaacerror.SysErrFailToGetPrometheusData.Error()
				logger.Error(internalError)
				message := isaacerror.GetAPIError(isaacerror.ErrorFailToGetPrometheusData, internalError)
				c.JSON(http.StatusInternalServerError, message)
				return
			}
			responseDataList.Data[i] = responseData.Data
		}

		response = responseDataList
		bytes := "bytes 0-" + strconv.Itoa(responseDataList.Total) + "/" + strconv.Itoa(responseDataList.Total)
		c.Header(constants.HTTPHeaderContentRange, bytes)
		c.Header(constants.HTTPHeaderXTotalCount, strconv.Itoa(responseDataList.Total))
	}

	c.JSON(http.StatusOK, response)
}

// PostHandler godoc
// @Tags Channels
// @Summary POST handler of channels
// @Description Create the channel.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param body body channels.Request true "The chennel data to be create"
// @Success 200 {object} channels.Response "Result for get the created channel resource"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /channels [post]
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
	if !utility.IsAlphanumericString(data.Data.Name) {
		// Invalid Parameter.
		internalError := isaacerror.SysErrInvalidParameter.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorInvalidParameter, internalError)
		c.JSON(http.StatusBadRequest, message)
		return
	}

	// Check duplicate channel name.
	duplicationChannel := db.GetConfigurationChannelInfoByName(data.Data.Name)
	if duplicationChannel.CHANNEL_NAME != "" {
		// Duplicated channel name.
		internalError := isaacerror.SysErrDuplicatedChannelName.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorDuplicatedChannelName, internalError)
		c.JSON(http.StatusConflict, message)
		return
	}

	if isExistNodesInDB(data.Data.Nodes) {
		// Selected node that not exist in DB.
		internalError := isaacerror.SysErrSelectedNodeThatNotExistInDB.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorSelectedNodeThatNotExistInDB, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	// Insert channel in database.
	id := db.InsertConfigurationChannelTBWithMapping(data.Data.Name, data.Data.Nodes...)
	if id == "" {
		// Failed insert channel.
		internalError := isaacerror.SysErrFailToInsertChannel.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToInsertChannel, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	// Update node configuration.
	changeChannelConfig()

	// Get node-channel mapping data in database.
	mappingTB := db.GetChannelPermissionNodes(id)

	// Get nodes map that can be convert node pk to node data.
	nodePKToData := db.GetNodePKToDataMap()

	// Convert the channel data to response data.
	var response Response
	response.Data.ID = id
	response.Data.Name = data.Data.Name
	response.Data.Nodes = make([]NodeData, len(mappingTB))
	for i, value := range mappingTB {
		response.Data.Nodes[i].ID = value.NODE_PK
		response.Data.Nodes[i].Name = nodePKToData[value.NODE_PK].NODE_NAME
		response.Data.Nodes[i].IP = nodePKToData[value.NODE_PK].NODE_IP
	}

	c.JSON(http.StatusOK, response)
	logger.Infof("Successfully created the channel, %s", response.Data.Name)
}

// PutHandler godoc
// @Tags Channels
// @Summary PUT handler of channels
// @Description Update the channel resource.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param id path string true "Channel ID to be update"
// @Param body body channels.Request true "The channel data to be update"
// @Success 200 {object} channels.Response "Result for get the updated channel resource"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /channels/{id} [put]
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
	if !utility.IsAlphanumericString(data.Data.Name) {
		// Invalid Parameter.
		internalError := isaacerror.SysErrInvalidParameter.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorInvalidParameter, internalError)
		c.JSON(http.StatusBadRequest, message)
		return
	}

	if isExistNodesInDB(data.Data.Nodes) {
		// Selected node that not exist in DB.
		internalError := isaacerror.SysErrSelectedNodeThatNotExistInDB.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorSelectedNodeThatNotExistInDB, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	// Get channel data in database.
	channelTB := db.GetConfigurationChannelInfo(id)
	if channelTB.CHANNEL_PK == "" {
		// No channel in DB.
		internalError := isaacerror.SysErrNoChannelInDB.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorNoChannelInDB, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	// Check duplicate channel name.
	if data.Data.Name != channelTB.CHANNEL_NAME {
		duplicationNode := db.GetConfigurationChannelInfoByName(data.Data.Name)
		if duplicationNode.CHANNEL_NAME != "" {
			// Duplicated channel name.
			internalError := isaacerror.SysErrDuplicatedChannelName.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorDuplicatedChannelName, internalError)
			c.JSON(http.StatusConflict, message)
			return
		}
	}

	// Check if it has changed node list.
	mappingTB := db.GetChannelPermissionNodes(id)
	currentNodeList := make([]string, len(mappingTB))
	for i, value := range mappingTB {
		currentNodeList[i] = value.NODE_PK
	}

	isEqualNode := utility.IsIdenticalSlice(data.Data.Nodes, currentNodeList)

	// Update node in database.
	if isEqualNode {
		err = db.UpdateConfigurationChannelInfo(id, data.Data.Name)
	} else {
		err = db.UpdateConfigurationChannelInfoWithMapping(id, data.Data.Name, data.Data.Nodes...)
	}
	if err != nil {
		// Failed update node.
		internalError := isaacerror.SysErrFailToUpdateChannel.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToUpdateChannel, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	// Update node configuration.
	changeChannelConfig()

	// Get node-channel mapping data in database.
	mappingTB = db.GetChannelPermissionNodes(id)

	// Get nodes map that can be convert node pk to node data.
	nodePKToData := db.GetNodePKToDataMap()

	// Convert the channel data to response data.
	var response Response
	response.Data.ID = id
	response.Data.Name = data.Data.Name
	response.Data.Nodes = make([]NodeData, len(mappingTB))
	for i, value := range mappingTB {
		response.Data.Nodes[i].ID = value.NODE_PK
		response.Data.Nodes[i].Name = nodePKToData[value.NODE_PK].NODE_NAME
		response.Data.Nodes[i].IP = nodePKToData[value.NODE_PK].NODE_IP
	}

	c.JSON(http.StatusOK, response)
	logger.Infof("Successfully updated the channel, %s, %s", data.Data.Name, response.Data.Name)
}

// DeleteHandler godoc
// @Tags Channels
// @Summary DELETE handler of channels
// @Description Delete the channel resource.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param id path string true "Channel ID to be delete"
// @Success 200 {object} channels.Response "Result for get the deleted channel resource"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /channels/{id} [delete]
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

	// Get node data in database to response data.
	channelTB := db.GetConfigurationChannelInfo(id)

	// Get node-channel mapping data in database.
	mappingTB := db.GetChannelPermissionNodes(id)

	// Delete node in database.
	err := db.DeleteConfigurationChannel(id)
	if err != nil {
		// Failed delete node.
		internalError := isaacerror.SysErrFailToDeleteChannel.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToDeleteChannel, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	// delete the YAML channel configuration file.
	changeChannelConfig()

	// Get nodes map that can be convert node pk to node data.
	nodePKToData := db.GetNodePKToDataMap()

	// Convert the channel data to response data.
	var response Response
	response.Data.ID = id
	response.Data.Name = channelTB.CHANNEL_NAME
	response.Data.Nodes = make([]NodeData, len(mappingTB))
	for i, value := range mappingTB {
		response.Data.Nodes[i].ID = value.NODE_PK
		response.Data.Nodes[i].Name = nodePKToData[value.NODE_PK].NODE_NAME
		response.Data.Nodes[i].IP = nodePKToData[value.NODE_PK].NODE_IP
	}

	c.JSON(http.StatusOK, response)
	logger.Infof("Successfully deleted the channel, %s", response.Data.Name)
}

func changeChannelConfig() {
	// Change the node data in configuration file.
	conf := configuration.Conf()

	// Get node-channel mapping data in DB.
	var mappingTB []db.NODE_CHANNEL_MAPPING_TB
	mappingTB = db.GetChannelPermissionNodesTable()

	// Get nodes map that can be convert node pk to node data.
	nodePKToData := db.GetNodePKToDataMap()

	// Get channelTB data in DB and set configuration.
	var channelTB []db.CONFIGURATION_CHANNEL_TB
	channelTB = db.GetConfigurationChannelTable()
	conf.Channel = make([]configuration.Channels, len(channelTB))
	for i, value := range channelTB {
		conf.Channel[i].Name = value.CHANNEL_NAME
		conf.Channel[i].Nodes = make([]string, 0)

		for _, mappingValue := range mappingTB {
			if value.CHANNEL_PK == mappingValue.CHANNEL_PK {
				conf.Channel[i].Nodes = append(conf.Channel[i].Nodes, nodePKToData[mappingValue.NODE_PK].NODE_NAME)
			}
		}
	}

	// Update  the YAML node configuration file.
	configuration.ChangeConfigFile(configuration.GetFilePath(), conf)
}

func isExistNodesInDB(nodeList []string) bool {
	nodeTB := db.GetConfigurationNodeTable()
	nodeListInDB := make([]string, len(nodeTB))
	for i, value := range nodeTB {
		nodeListInDB[i] = value.NODE_NAME
	}

	for _, value := range nodeList {
		if !utility.IsExistValueInList(value, nodeListInDB) {
			return false
		}
	}

	return true
}

func getResponseDataByChannelPK(channelPK string) (*Response, error) {
	var channelTB *db.CONFIGURATION_CHANNEL_TB
	channelTB = db.GetConfigurationChannelInfo(channelPK)

	// Get prometheus data.
	loopchainChannelData, err := prometheus.GetPrometheusChannelData(channelTB.CHANNEL_NAME)
	if err != nil {
		// Fail to get prometheus data.
		return nil, err
	}

	// Get all node channel mapping data.
	var mappingTB []db.NODE_CHANNEL_MAPPING_TB
	mappingTB = db.GetChannelPermissionNodes(channelPK)

	// Get all node resource and generate nodes list that can be convert node pk to node name.
	nodePKToData := db.GetNodePKToDataMap()

	var responseData Response
	responseData.Data.ID = channelPK
	responseData.Data.Name = channelTB.CHANNEL_NAME
	responseData.Data.GoloopChannelID = channelTB.CHANNEL_ID
	responseData.Data.Total = len(mappingTB)
	responseData.Data.Status = &loopchainChannelData.Status
	responseData.Data.Nodes = make([]NodeData, len(mappingTB))

	for i, value := range mappingTB {
		responseData.Data.Nodes[i].ID = value.NODE_PK
		responseData.Data.Nodes[i].Name = nodePKToData[value.NODE_PK].NODE_NAME
		responseData.Data.Nodes[i].IP = nodePKToData[value.NODE_PK].NODE_IP
		for _, prometheusData := range loopchainChannelData.Nodes {

			if responseData.Data.Nodes[i].Name == prometheusData.Name {
				responseData.Data.Nodes[i].BlockHeight = &prometheusData.BlockHeight
				responseData.Data.Nodes[i].CountOfTX = &prometheusData.CountOfTX
				responseData.Data.Nodes[i].CountOfUnconfirmedTX = &prometheusData.CountOfUnconfirmedTX
				responseData.Data.Nodes[i].ResponseTimeInSec = &prometheusData.ResponseTimeInSec
				responseData.Data.Nodes[i].IsLeader = &prometheusData.IsLeader
				responseData.Data.Nodes[i].TimeStamp = prometheusData.TimeStamp
				responseData.Data.Nodes[i].Status = &prometheusData.Status
				break
			}
		}
	}

	for i, value := range loopchainChannelData.Nodes {
		if value.IsLeader == 1 || i == 0 || value.Status != 0 {
			responseData.Data.BlockHeight = &loopchainChannelData.Nodes[i].BlockHeight
			responseData.Data.CountOfTX = &loopchainChannelData.Nodes[i].CountOfTX
			responseData.Data.ResponseTimeInSec = &loopchainChannelData.Nodes[i].ResponseTimeInSec
			if value.Status != 0 {
				break
			}
		}
	}
	return &responseData, nil
}
