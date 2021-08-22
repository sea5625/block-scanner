package nodes

import (
	"github.com/gin-gonic/gin"
	"motherbear/backend/configuration"
	"motherbear/backend/constants"
	"motherbear/backend/db"
	"motherbear/backend/isaacerror"
	"motherbear/backend/logger"
	"motherbear/backend/prometheus"
	"motherbear/backend/utility"
	"net/http"
	"strconv"
)

type Request struct {
	Data struct {
		Name string `json:"name" example:"node1"`
		IP   string `json:"ip" example:"https://int-test-ctz.solidwallet.io"`
	} `json:"data"`
}

type ResponseList struct {
	Data  []NodeData `json:"data"`
	Total int        `json:"total" example:"1"`
}

type Response struct {
	Data NodeData `json:"data"`
}

type NodeData struct {
	ID                   string   `json:"id" example:"PKND_0000000000000001"`
	Name                 string   `json:"name" example:"node1"`
	IP                   string   `json:"ip" example:"https://int-test-ctz.solidwallet.io"`
	BlockHeight          *uint64  `json:"blockHeight,omitempty" example:"10" format:"uint64"`
	CountOfTX            *uint64  `json:"countOfTX,omitempty" example:"1000" format:"uint64"`
	CountOfUnconfirmedTX *uint64  `json:"countOfUnconfirmedTX,omitempty" example:"0" format:"uint64"`
	ResponseTimeInSec    *float64 `json:"responseTimeInSec,omitempty" example:"0.001" format:"float64"`
	IsLeader             *int     `json:"isLeader,omitempty" example:"0" format:"int32"`
	TimeStamp            string   `json:"timeStamp,omitempty" example:"2006-01-02T15:04:05Z07:00"`
	Status               *int     `json:"status,omitempty" example:"0" format:"int32"`
}

// GetHandlerList godoc
// @Tags Nodes
// @Summary GET handler of nodes
// @Description Get all node resources or many node resources in the channel.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param channel query string false "It is ID of channel to be get nodes list."
// @Success 200 {object} nodes.ResponseList "Result for get many node resource."
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /nodes [get]
func GetHandlerList(c *gin.Context) {

}

// GetHandler godoc
// @Tags Nodes
// @Summary GET handler of nodes
// @Description Get the one node resource.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param id path string true "Node ID to be get"
// @Success 200 {object} nodes.Response "Result for get the one node resource"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /nodes/{id} [get]
func GetHandler(c *gin.Context) {
	id := c.Param(constants.RequestResourceID)

	var response interface{}
	if id != "" {
		// Get one node resource.
		var nodeTB *db.CONFIGURATION_NODE_TB
		nodeTB = db.GetConfigurationNodeInfoByNodePK(id)

		if nodeTB.NODE_PK == "" {
			// No node in DB.
			internalError := isaacerror.SysErrNoNodeInDB.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorNoNodeInDB, internalError)
			c.JSON(http.StatusInternalServerError, message)
			return
		}

		var responseStruct Response
		var data NodeData

		data.ID = nodeTB.NODE_PK
		data.Name = nodeTB.NODE_NAME
		data.IP = nodeTB.NODE_IP

		responseStruct.Data = data

		response = responseStruct
	} else {
		// Get many node resources.

		var responseDataList ResponseList

		channelID := c.Query("channel")
		if channelID == "" {
			// Get all node resource.
			var nodeTB []db.CONFIGURATION_NODE_TB
			nodeTB = db.GetConfigurationNodeTable()
			nodeTBLen := len(nodeTB)

			responseDataList.Data = make([]NodeData, nodeTBLen)
			responseDataList.Total = nodeTBLen

			for i := 0; i < nodeTBLen; i++ {
				responseDataList.Data[i].ID = nodeTB[i].NODE_PK
				responseDataList.Data[i].Name = nodeTB[i].NODE_NAME
				responseDataList.Data[i].IP = nodeTB[i].NODE_IP
			}
		} else {
			// Get many node resources in the channel.

			channelTB := db.GetConfigurationChannelInfo(channelID)
			if channelTB.CHANNEL_PK == "" {
				// No channel in DB.
				internalError := isaacerror.SysErrNoChannelInDB.Error()
				logger.Error(internalError)
				message := isaacerror.GetAPIError(isaacerror.ErrorNoChannelInDB, internalError)
				c.JSON(http.StatusInternalServerError, message)
				return
			}

			// Get node-channel mapping data in database.
			mappingTB := db.GetChannelPermissionNodes(channelID)
			nodeLen := len(mappingTB)

			// Get nodes map that can be convert node pk to node data.
			nodePKToData := db.GetNodePKToDataMap()

			// Get prometheus data.
			loopchainChannelData, err := prometheus.GetPrometheusChannelData(channelTB.CHANNEL_NAME)
			if err != nil {
				// Fail to get prometheus data.
				internalError := isaacerror.SysErrFailToGetPrometheusData.Error()
				logger.Error(internalError)
				message := isaacerror.GetAPIError(isaacerror.ErrorFailToGetPrometheusData, internalError)
				c.JSON(http.StatusInternalServerError, message)
				return
			}

			// Convert the node data to response data.
			responseDataList.Data = make([]NodeData, nodeLen)
			responseDataList.Total = nodeLen

			for i, value := range mappingTB {
				responseDataList.Data[i].ID = value.NODE_PK
				responseDataList.Data[i].Name = nodePKToData[value.NODE_PK].NODE_NAME
				responseDataList.Data[i].IP = nodePKToData[value.NODE_PK].NODE_IP

				for _, prometheusData := range loopchainChannelData.Nodes {
					if responseDataList.Data[i].Name == prometheusData.Name {
						responseDataList.Data[i].BlockHeight = &prometheusData.BlockHeight
						responseDataList.Data[i].CountOfTX = &prometheusData.CountOfTX
						responseDataList.Data[i].CountOfUnconfirmedTX = &prometheusData.CountOfUnconfirmedTX
						responseDataList.Data[i].ResponseTimeInSec = &prometheusData.ResponseTimeInSec
						responseDataList.Data[i].IsLeader = &prometheusData.IsLeader
						responseDataList.Data[i].TimeStamp = prometheusData.TimeStamp
						responseDataList.Data[i].Status = &prometheusData.Status
						break
					}
				}
			}
		}

		response = responseDataList
		bytes := "bytes 0-" + strconv.Itoa(responseDataList.Total) + "/" + strconv.Itoa(responseDataList.Total)
		c.Header(constants.HTTPHeaderContentRange, bytes)
		c.Header(constants.HTTPHeaderXTotalCount, strconv.Itoa(responseDataList.Total))
	}

	c.JSON(http.StatusOK, response)
}

// PostHandler godoc
// @Tags Nodes
// @Summary POST handler of nodes
// @Description Create the node resource.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param body body nodes.Request true "The node data to be create"
// @Success 200 {object} nodes.Response "Result for get the created node resource"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /nodes [post]
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
	if !utility.IsAlphanumericString(data.Data.Name) || data.Data.IP == "" {
		// Invalid Parameter.
		internalError := isaacerror.SysErrInvalidParameter.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorInvalidParameter, internalError)
		c.JSON(http.StatusBadRequest, message)
		return
	}

	// Check duplicate node name.
	duplicationNode := db.GetConfigurationNodeInfoByNodeName(data.Data.Name)
	if duplicationNode.NODE_NAME != "" {
		// Duplicated node name.
		internalError := isaacerror.SysErrDuplicatedNodeName.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorDuplicatedNodeName, internalError)
		c.JSON(http.StatusConflict, message)
		return
	}

	// Insert node in database.
	id := db.InsertConfigurationNode(data.Data.Name, data.Data.IP)
	if id == "" {
		// Failed insert node.
		internalError := isaacerror.SysErrFailToInsertNode.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToInsertNode, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	// Update node configuration.
	changeNodeConfig()

	// Convert node data to response data.
	var response Response
	response.Data.ID = id
	response.Data.Name = data.Data.Name
	response.Data.IP = data.Data.IP

	c.JSON(http.StatusOK, response)
	logger.Infof("Successfully created the node, %s", response.Data.Name)
}

// PutHandler godoc
// @Tags Nodes
// @Summary PUT handler of nodes
// @Description Update the node resource.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param id path string true "Node ID to be update"
// @Param body body nodes.Request true "The node data to be update"
// @Success 200 {object} nodes.Response "Result for get the updated node resource"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /nodes/{id} [put]
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
	if !utility.IsAlphanumericString(data.Data.Name) || data.Data.IP == "" {
		// Invalid Parameter.
		internalError := isaacerror.SysErrInvalidParameter.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorInvalidParameter, internalError)
		c.JSON(http.StatusBadRequest, message)
		return
	}

	// Get node data in database.
	nodeTB := db.GetConfigurationNodeInfoByNodePK(id)
	if nodeTB.NODE_PK == "" {
		// No channel in DB.
		internalError := isaacerror.SysErrNoNodeInDB.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorNoNodeInDB, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	// Check duplicate node name.
	if data.Data.Name != nodeTB.NODE_NAME {
		duplicationNode := db.GetConfigurationNodeInfoByNodeName(data.Data.Name)
		if duplicationNode.NODE_NAME != "" {
			// Duplicated node name.
			internalError := isaacerror.SysErrDuplicatedNodeName.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorDuplicatedNodeName, internalError)
			c.JSON(http.StatusConflict, message)
			return
		}
	}

	// Update node in database.
	err = db.UpdateConfigurationNodeInfo(id, data.Data.Name, data.Data.IP)
	if err != nil {
		// Failed update node.
		internalError := isaacerror.SysErrFailToUpdateNode.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorFailTodUpdateNode, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	// Update node configuration.
	changeNodeConfig()

	// Convert node data to response data.
	var response Response
	response.Data.ID = id
	response.Data.Name = data.Data.Name
	response.Data.IP = data.Data.IP

	c.JSON(http.StatusOK, response)
	logger.Infof("Successfully updated the node, %s, %s", nodeTB.NODE_NAME, response.Data.Name)
}

// DeleteHandler godoc
// @Tags Nodes
// @Summary DELETE handler of nodes
// @Description Delete the node resource.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param id path string true "Node ID to be delete"
// @Success 200 {object} nodes.Response "Result for get the deleted node resource"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /nodes/{id} [delete]
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
	nodeTB := db.GetConfigurationNodeInfoByNodePK(id)

	// Delete node in database.
	err := db.DeleteConfigurationNodeInfo(id)
	if err != nil {
		// Failed delete node.
		internalError := isaacerror.SysErrFailToDeleteNode.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToDeleteNode, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	// Update  the YAML node configuration file.
	changeNodeConfig()

	// Convert the node data to response data.
	var response Response
	response.Data.ID = id
	response.Data.Name = nodeTB.NODE_NAME
	response.Data.IP = nodeTB.NODE_IP

	c.JSON(http.StatusOK, response)
	logger.Infof("Successfully deleted the node, %s", response.Data.Name)
}

func changeNodeConfig() {
	// Change the node data in configuration file.
	conf := configuration.Conf()

	// Get nodeTB data in DB and set configuration.
	var nodeTB []db.CONFIGURATION_NODE_TB
	nodeTB = db.GetConfigurationNodeTable()
	conf.Node = make([]configuration.Nodes, len(nodeTB))
	for i, value := range nodeTB {
		conf.Node[i].Name = value.NODE_NAME
		conf.Node[i].IP = value.NODE_IP
	}

	// Update  the YAML node configuration file.
	configuration.ChangeConfigFile(configuration.GetFilePath(), conf)
}
