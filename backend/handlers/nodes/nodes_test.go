package nodes

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
	"motherbear/backend/configuration"
	"motherbear/backend/constants"
	"motherbear/backend/db"
	"motherbear/backend/isaacerror"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

const dbPath string = ":memory:"
const confFilePath string = "testConfiguration.yaml"

var nodeTestData = []configuration.Nodes{
	{
		Name: "node1",
		IP:   "https://int-test-ctz.solidwallet.io",
	},
	{
		Name: "node2",
		IP:   "https://int-test-ctz.solidwallet.io",
	},
	{
		Name: "node3",
		IP:   "https://int-test-ctz.solidwallet.io",
	},
}
var channelTestData = []configuration.Channels{
	{
		Name:  "channel1",
		Nodes: []string{nodeTestData[0].Name, nodeTestData[1].Name},
	},
}
var nodeResponseList ResponseList          // Correct nodes list.
var nodeInChannelResponseList ResponseList // Correct nodes list in channel.

func TestGetHandler(t *testing.T) {
	Setup()
	defer Teardown()

	// Test get nodes list
	router := gin.Default()
	router.GET(constants.NodesGetListAPIURL, GetHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodGET, constants.NodesGetListAPIURL, nil)
	router.ServeHTTP(w, request)

	var resultList ResponseList
	_ = json.Unmarshal(w.Body.Bytes(), &resultList)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, strconv.Itoa(nodeResponseList.Total), w.Header().Get("X-Total-Count"))
	assert.Equal(t, true, checkNodeResponseList(nodeResponseList, resultList))
}

// Test GET node API.
func TestGetHandler2(t *testing.T) {
	Setup()
	defer Teardown()

	// Test get node list
	router := gin.Default()
	router.GET(constants.NodesGetListAPIURL, GetHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodGET, constants.NodesGetListAPIURL, nil)
	router.ServeHTTP(w, request)

	var resultList ResponseList
	_ = json.Unmarshal(w.Body.Bytes(), &resultList)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, strconv.Itoa(nodeResponseList.Total), w.Header().Get("X-Total-Count"))
	assert.Equal(t, true, checkNodeResponseList(nodeResponseList, resultList))

	// Test get one node
	node := resultList.Data[0].ID
	router = gin.Default()
	router.GET(constants.NodesGetAPIURL, GetHandler)

	w = httptest.NewRecorder()
	request, _ = http.NewRequest(constants.HTTPMethodGET, constants.NodesAPIBaseURL+"/"+node, nil)
	router.ServeHTTP(w, request)

	var result Response
	_ = json.Unmarshal(w.Body.Bytes(), &result)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, nodeResponseList.Data[0].Name, result.Data.Name)
	assert.Equal(t, nodeResponseList.Data[0].IP, result.Data.IP)
}

// Test get nodes list in channel.
func TestGetHandler3(t *testing.T) {
	t.Skip("Skipping test RequestPrometheusQuery. Should prepare prometheus server to test.")
	Setup()
	defer Teardown()

	// Test get nodes list in channel.
	channelTB := db.GetConfigurationChannelTable()
	channelID := channelTB[0].CHANNEL_PK
	router := gin.Default()
	router.GET(constants.ChannelsGetListAPIURL, GetHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodGET, constants.NodesGetListAPIURL+"?"+constants.RequestParamChannel+"="+channelID, nil)
	router.ServeHTTP(w, request)

	var resultList ResponseList
	_ = json.Unmarshal(w.Body.Bytes(), &resultList)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, strconv.Itoa(nodeInChannelResponseList.Total), w.Header().Get("X-Total-Count"))
	assert.Equal(t, true, checkNodeResponseList(nodeInChannelResponseList, resultList))
}

// Test to return error if using invalid channel.
func TestGetHandler4(t *testing.T) {
	Setup()
	defer Teardown()

	// Test get nodes list in channel.
	channelID := "asdf"
	router := gin.Default()
	router.GET(constants.NodesGetListAPIURL, GetHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodGET, constants.NodesGetListAPIURL+"?"+constants.RequestParamChannel+"="+channelID, nil)
	router.ServeHTTP(w, request)

	var err isaacerror.APIError
	_ = json.Unmarshal(w.Body.Bytes(), &err)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, isaacerror.ErrorNoChannelInDB, err.Errors[0].UserMessage)
}

func TestPostHandler(t *testing.T) {
	Setup()
	defer Teardown()

	// Set create node data.
	var addData NodeData
	addData.Name = "node4"
	addData.IP = "https://int-test-ctz.solidwallet.io"
	addData.ID = "PKND_0000000000000004"
	nodeResponseList.Data = append(nodeResponseList.Data, addData)
	nodeResponseList.Total = len(nodeResponseList.Data)

	// Set request data.
	var requestData Request
	requestData.Data.Name = addData.Name
	requestData.Data.IP = addData.IP
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.POST(constants.NodesPostAPIURL, PostHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPOST, constants.NodesPostAPIURL, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	var result Response
	_ = json.Unmarshal(w.Body.Bytes(), &result)

	// Check response data.
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, addData.Name, result.Data.Name)
	assert.Equal(t, addData.IP, result.Data.IP)

	// Check database.
	var node []db.CONFIGURATION_NODE_TB
	node = db.GetConfigurationNodeTable()
	assert.Equal(t, true, checkNodeDBList(nodeResponseList, node))

	// Check configuration.
	conf := configuration.InitConfigData(confFilePath)
	assert.Equal(t, true, checkNodeConfigList(nodeResponseList, conf.Node))
}

// Test to return error when duplicate node name.
func TestPostHandler2(t *testing.T) {
	Setup()
	defer Teardown()

	// Set create node data.
	var addData NodeData
	addData.Name = "node3"
	addData.IP = "https://int-test-ctz.solidwallet.io"
	addData.ID = "PKND_0000000000000004"
	nodeResponseList.Data = append(nodeResponseList.Data, addData)
	nodeResponseList.Total = len(nodeResponseList.Data)

	// Set request data.
	var requestData Request
	requestData.Data.Name = addData.Name
	requestData.Data.IP = addData.IP
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.POST(constants.NodesPostAPIURL, PostHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPOST, constants.NodesPostAPIURL, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	// Check duplicate node name.
	var err isaacerror.APIError
	_ = json.Unmarshal(w.Body.Bytes(), &err)

	assert.Equal(t, 409, w.Code)
	assert.Equal(t, isaacerror.ErrorDuplicatedNodeName, err.Errors[0].UserMessage)
}

// Test to return error when invalid node name.
func TestPostHandler3(t *testing.T) {
	Setup()
	defer Teardown()

	// Set create node data.
	var addData NodeData
	addData.Name = "node4!@#"
	addData.IP = "https://int-test-ctz.solidwallet.io"
	addData.ID = "PKND_0000000000000004"
	nodeResponseList.Data = append(nodeResponseList.Data, addData)
	nodeResponseList.Total = len(nodeResponseList.Data)

	// Set request data.
	var requestData Request
	requestData.Data.Name = addData.Name
	requestData.Data.IP = addData.IP
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.POST(constants.NodesPostAPIURL, PostHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPOST, constants.NodesPostAPIURL, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	// Check duplicate node name.
	var err isaacerror.APIError
	_ = json.Unmarshal(w.Body.Bytes(), &err)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, isaacerror.ErrorInvalidParameter, err.Errors[0].UserMessage)
}

func TestPutHandler(t *testing.T) {
	Setup()
	defer Teardown()

	// Set update node data.
	whereName := "node1"
	updateNodeTB := db.GetConfigurationNodeInfoByNodeName(whereName)
	whereID := updateNodeTB.NODE_PK

	var updateData NodeData
	updateData.Name = "node4"
	updateData.IP = "1.1.1.1"
	updateData.ID = whereID
	for i, value := range nodeResponseList.Data {
		if value.Name == whereName {
			nodeResponseList.Data[i].Name = updateData.Name
			nodeResponseList.Data[i].IP = updateData.IP
		}
	}

	// Set request data.
	var requestData Request
	requestData.Data.Name = updateData.Name
	requestData.Data.IP = updateData.IP
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.PUT(constants.NodesPutAPIURL, PutHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPUT, constants.NodesAPIBaseURL+"/"+whereID, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	var result Response
	_ = json.Unmarshal(w.Body.Bytes(), &result)

	// Check response data.
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, updateData.Name, result.Data.Name)
	assert.Equal(t, updateData.IP, result.Data.IP)

	// Check database.
	var node []db.CONFIGURATION_NODE_TB
	node = db.GetConfigurationNodeTable()
	assert.Equal(t, true, checkNodeDBList(nodeResponseList, node))

	// Check configuration.
	conf := configuration.InitConfigData(confFilePath)
	assert.Equal(t, true, checkNodeConfigList(nodeResponseList, conf.Node))
}

// Test to return error when duplicate node name.
func TestPutHandler2(t *testing.T) {
	Setup()
	defer Teardown()

	// Set update node data.
	whereName := "node1"
	updateNodeTB := db.GetConfigurationNodeInfoByNodeName(whereName)
	whereID := updateNodeTB.NODE_PK

	var updateData NodeData
	updateData.Name = "node3"
	updateData.IP = "1.1.1.1"
	updateData.ID = whereID
	for i, value := range nodeResponseList.Data {
		if value.Name == whereName {
			nodeResponseList.Data[i].Name = updateData.Name
			nodeResponseList.Data[i].IP = updateData.IP
		}
	}

	// Set request data.
	var requestData Request
	requestData.Data.Name = updateData.Name
	requestData.Data.IP = updateData.IP
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.PUT(constants.NodesPutAPIURL, PutHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPUT, constants.NodesAPIBaseURL+"/"+whereID, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	// Check duplicate node name.
	var err isaacerror.APIError
	_ = json.Unmarshal(w.Body.Bytes(), &err)

	assert.Equal(t, 409, w.Code)
	assert.Equal(t, isaacerror.ErrorDuplicatedNodeName, err.Errors[0].UserMessage)
}

func TestDeleteHandler(t *testing.T) {
	Setup()
	defer Teardown()

	// Set delete node data.
	deleteName := "node1"
	deleteNodeTB := db.GetConfigurationNodeInfoByNodeName(deleteName)
	deleteID := deleteNodeTB.NODE_PK

	var deleteData NodeData
	deleteData.Name = deleteName
	deleteData.IP = deleteNodeTB.NODE_IP
	deleteData.ID = deleteNodeTB.NODE_PK
	for i, value := range nodeResponseList.Data {
		if value.Name == deleteName {
			copy(nodeResponseList.Data[i:], nodeResponseList.Data[i+1:])
			nodeResponseList.Data[len(nodeResponseList.Data)-1] = NodeData{}
			nodeResponseList.Data = nodeResponseList.Data[:len(nodeResponseList.Data)-1]
		}
	}
	nodeResponseList.Total -= 1

	router := gin.Default()
	router.DELETE(constants.NodesDeleteAPIURL, DeleteHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodDELETE, constants.NodesAPIBaseURL+"/"+deleteID, nil)
	router.ServeHTTP(w, request)

	var result Response
	_ = json.Unmarshal(w.Body.Bytes(), &result)

	// Check response data.
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, deleteData.Name, result.Data.Name)
	assert.Equal(t, deleteData.IP, result.Data.IP)

	// Check database.
	var node []db.CONFIGURATION_NODE_TB
	node = db.GetConfigurationNodeTable()
	assert.Equal(t, true, checkNodeDBList(nodeResponseList, node))

	// Check configuration.
	conf := configuration.InitConfigData(confFilePath)
	assert.Equal(t, true, checkNodeConfigList(nodeResponseList, conf.Node))
}

func checkNodeResponseList(var1, var2 ResponseList) bool {
	var1Len := len(var1.Data)
	var2Len := len(var2.Data)

	// Check length.
	if var1Len != var2Len {
		return false
	}

	if var1.Total != var2.Total {
		return false
	}

	// Check node name and ip.
	for _, value1 := range var1.Data {
		check := false

		for _, value2 := range var2.Data {
			if value1.Name == value2.Name {
				if value1.IP != value2.IP {
					return false
				}

				check = true
				break
			}
		}

		if !check {
			return false
		}
	}

	return true
}

func checkNodeDBList(var1 ResponseList, var2 []db.CONFIGURATION_NODE_TB) bool {
	var1Len := len(var1.Data)
	var2Len := len(var2)

	// Check length.
	if var1Len != var2Len {
		return false
	}

	if var1.Total != var2Len {
		return false
	}

	// Check node name and ip.
	for _, value1 := range var1.Data {
		check := false

		for _, value2 := range var2 {
			if value1.Name == value2.NODE_NAME {
				if value1.IP != value2.NODE_IP {
					return false
				}

				check = true
				break
			}
		}

		if !check {
			return false
		}
	}

	return true
}

func checkNodeConfigList(var1 ResponseList, var2 []configuration.Nodes) bool {
	var1Len := len(var1.Data)
	var2Len := len(var2)

	// Check length.
	if var1Len != var2Len {
		return false
	}

	if var1.Total != var2Len {
		return false
	}

	// Check node name and ip.
	for _, value1 := range var1.Data {
		check := false

		for _, value2 := range var2 {
			if value1.Name == value2.Name {
				if value1.IP != value2.IP {
					return false
				}

				check = true
				break
			}
		}

		if !check {
			return false
		}
	}

	return true
}

func Setup() {
	var conf configuration.Configuration

	// Add node and channel configuration.
	conf.Node = nodeTestData
	conf.Channel = channelTestData

	configuration.ChangeConfigFile(confFilePath, &conf)
	configuration.InitConfigData(confFilePath)

	// Create database.
	db.InitDB("sqlite3", dbPath)
	db.InitCreateTable()

	// Get node data in database and convert node data to response data.
	nodeTB := db.GetConfigurationNodeTable()
	nodeResponseList.Data = make([]NodeData, len(nodeTB))
	nodeResponseList.Total = len(nodeTB)
	for i, value := range nodeTB {
		nodeResponseList.Data[i].ID = value.NODE_PK
		nodeResponseList.Data[i].Name = value.NODE_NAME
		nodeResponseList.Data[i].IP = value.NODE_IP
	}

	// Convert node data in channel to response data.
	nodeInChannelResponseList.Data = make([]NodeData, len(channelTestData[0].Nodes))
	nodeInChannelResponseList.Total = len(channelTestData[0].Nodes)
	for i, nodeInChannel := range channelTestData[0].Nodes {
		for _, value := range nodeTB {
			if nodeInChannel == value.NODE_NAME {
				nodeInChannelResponseList.Data[i].Name = value.NODE_NAME
				nodeInChannelResponseList.Data[i].ID = value.NODE_PK
				nodeInChannelResponseList.Data[i].IP = value.NODE_IP
			}
		}
	}
}

func Teardown() {
	_ = os.Remove(confFilePath)
}
