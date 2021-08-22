package channels

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
	"motherbear/backend/configuration"
	"motherbear/backend/constants"
	"motherbear/backend/db"
	"motherbear/backend/isaacerror"
	"motherbear/backend/prometheus/prom_crawler"
	"motherbear/backend/utility"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
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
		Nodes: []string{nodeTestData[0].Name, nodeTestData[1].Name, nodeTestData[2].Name},
	},
	{
		Name:  "channel2",
		Nodes: []string{nodeTestData[0].Name, nodeTestData[1].Name, nodeTestData[2].Name},
	},
	{
		Name:  "channel3",
		Nodes: []string{nodeTestData[0].Name, nodeTestData[1].Name, nodeTestData[2].Name},
	},
}
var prometheusTestData = configuration.Prometheus{
	PrometheusISAAC:    "http://localhost:9090",
	PrometheusExternal: "http://localhost:9090",
	QueryPath:          "/api/v1/query",
	CrawlingInterval:   5,
}
var etcTestData = configuration.Etc{
	SessionTimeout: 30,
	Language:       "ko",
}

var channelResponse Response         // Correct channel API response data
var channelResponseList ResponseList // Correct channel list API response data

func TestGetHandler(t *testing.T) {
	t.Skip("Skipping test RequestPrometheusQuery. Should prepare prometheus server to test.")
	Setup()
	defer Teardown()
	prom_crawler.BeginToCrawl()

	// Sleeping during crawling.
	time.Sleep(time.Duration(configuration.Conf().Prometheus.CrawlingInterval+1) * time.Second)

	// Test get channel list
	router := gin.Default()
	router.GET(constants.ChannelsGetListAPIURL, GetHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodGET, constants.ChannelsGetListAPIURL, nil)
	router.ServeHTTP(w, request)

	var resultList ResponseList
	_ = json.Unmarshal(w.Body.Bytes(), &resultList)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, strconv.Itoa(channelResponseList.Total), w.Header().Get("X-Total-Count"))
	assert.Equal(t, true, isEqualChannelResponseList(channelResponseList, resultList))

	// Test get one channel
	channel := resultList.Data[0].ID
	router = gin.Default()
	router.GET(constants.ChannelsGetAPIURL, GetHandler)

	w = httptest.NewRecorder()
	request, _ = http.NewRequest(constants.HTTPMethodGET, constants.ChannelsAPIBaseURL+"/"+channel, nil)
	router.ServeHTTP(w, request)

	var result Response
	_ = json.Unmarshal(w.Body.Bytes(), &result)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, channelResponseList.Data[0].Name, result.Data.Name)
	assert.Equal(t, true, isEqualNodeData(channelResponseList.Data[0].Nodes, result.Data.Nodes))
}

// Test to get channel list only in channel permission list.
func TestGetHandler2(t *testing.T) {
	t.Skip("Skipping test RequestPrometheusQuery. Should prepare prometheus server to test.")

	Setup()
	defer Teardown()
	prom_crawler.BeginToCrawl()

	// Sleeping during crawling.
	time.Sleep(time.Duration(configuration.Conf().Prometheus.CrawlingInterval+1) * time.Second)

	// Test get channel list only in channel permission list.
	var channelPermissionList []string
	for i := 0; i < len(channelResponseList.Data)-1; i++ {
		channelPermissionList = append(channelPermissionList, channelResponseList.Data[i].ID)
	}
	channelResponseList.Total = len(channelResponseList.Data) - 1
	channelResponseList.Data = channelResponseList.Data[:len(channelResponseList.Data)-1]

	router := gin.Default()
	router.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set(constants.ContextKeyPermissionChannelList, channelPermissionList)
			return
		}
	}())

	router.GET(constants.ChannelsGetListAPIURL, GetHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodGET, constants.ChannelsGetListAPIURL, nil)
	router.ServeHTTP(w, request)

	var result ResponseList
	_ = json.Unmarshal(w.Body.Bytes(), &result)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, len(channelPermissionList), len(result.Data))
	assert.Equal(t, true, isEqualChannelResponseList(channelResponseList, result))
}

func TestPostHandler(t *testing.T) {
	Setup()
	defer Teardown()

	// Set create channel data.
	var addData ChannelData
	addData.Name = "channel4"
	addData.ID = "PKCH_0000000000000004"
	addData.Nodes = make([]NodeData, len(channelResponseList.Data[0].Nodes))
	copy(addData.Nodes, channelResponseList.Data[0].Nodes)
	channelResponseList.Data = append(channelResponseList.Data, addData)
	channelResponseList.Total = len(channelResponseList.Data)

	// Set request data.
	var requestData Request
	requestData.Data.Name = addData.Name
	requestData.Data.Nodes = make([]string, len(addData.Nodes))
	for i, value := range addData.Nodes {
		requestData.Data.Nodes[i] = value.ID
	}
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.POST(constants.ChannelsPostAPIURL, PostHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPOST, constants.ChannelsPostAPIURL, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	var result Response
	_ = json.Unmarshal(w.Body.Bytes(), &result)

	// Check response data.
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, addData.Name, result.Data.Name)
	assert.Equal(t, true, isEqualNodeData(addData.Nodes, result.Data.Nodes))

	// Check database.
	assert.Equal(t, true, isEqualDB(channelResponseList))

	// Check configuration.
	conf := configuration.InitConfigData(confFilePath)
	assert.Equal(t, true, isEqualChannelConfigList(channelResponseList, conf.Channel))
}

// Test to return error when duplicate channel name.
func TestPostHandler2(t *testing.T) {
	Setup()
	defer Teardown()

	// Set create channel data.
	var addData ChannelData
	addData.Name = "channel1"
	addData.ID = "PKCH_0000000000000004"
	addData.Nodes = make([]NodeData, len(channelResponseList.Data[0].Nodes))
	copy(addData.Nodes, channelResponseList.Data[0].Nodes)
	channelResponseList.Data = append(channelResponseList.Data, addData)
	channelResponseList.Total = len(channelResponseList.Data)

	// Set request data.
	var requestData Request
	requestData.Data.Name = addData.Name
	requestData.Data.Nodes = make([]string, len(addData.Nodes))
	for i, value := range addData.Nodes {
		requestData.Data.Nodes[i] = value.ID
	}
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.POST(constants.ChannelsPostAPIURL, PostHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPOST, constants.ChannelsPostAPIURL, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	// Check duplicate node name.
	var err isaacerror.APIError
	_ = json.Unmarshal(w.Body.Bytes(), &err)

	assert.Equal(t, 409, w.Code)
	assert.Equal(t, isaacerror.ErrorDuplicatedChannelName, err.Errors[0].UserMessage)
}

// Test to return error when invalid channel name.
func TestPostHandler3(t *testing.T) {
	Setup()
	defer Teardown()

	// Set create channel data.
	var addData ChannelData
	addData.Name = "channel4!@#"
	addData.ID = "PKCH_0000000000000004"
	addData.Nodes = make([]NodeData, len(channelResponseList.Data[0].Nodes))
	copy(addData.Nodes, channelResponseList.Data[0].Nodes)
	channelResponseList.Data = append(channelResponseList.Data, addData)
	channelResponseList.Total = len(channelResponseList.Data)

	// Set request data.
	var requestData Request
	requestData.Data.Name = addData.Name
	requestData.Data.Nodes = make([]string, len(addData.Nodes))
	for i, value := range addData.Nodes {
		requestData.Data.Nodes[i] = value.ID
	}
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.POST(constants.ChannelsPostAPIURL, PostHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPOST, constants.ChannelsPostAPIURL, bytes.NewBuffer(requestJSON))
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

	// Set update channel data.
	whereName := "channel1"
	updateChannelTB := db.GetConfigurationChannelInfoByName(whereName)
	whereID := updateChannelTB.CHANNEL_PK

	var updateData ChannelData
	updateData.Name = "channel4"
	updateData.ID = whereID
	updateData.Nodes = make([]NodeData, len(channelResponseList.Data[0].Nodes)-1)
	for i := 0; i < len(updateData.Nodes); i++ {
		updateData.Nodes[i] = channelResponseList.Data[0].Nodes[i+1]
	}
	for i, value := range channelResponseList.Data {
		if value.Name == whereName {
			channelResponseList.Data[i].Name = updateData.Name
			channelResponseList.Data[i].Nodes = updateData.Nodes
		}
	}

	// Set request data.
	var requestData Request
	requestData.Data.Name = updateData.Name
	requestData.Data.Nodes = make([]string, len(updateData.Nodes))
	for i, value := range updateData.Nodes {
		requestData.Data.Nodes[i] = value.ID
	}
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.PUT(constants.ChannelsPutAPIURL, PutHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPUT, constants.ChannelsAPIBaseURL+"/"+whereID, bytes.NewBuffer(requestJSON))

	router.ServeHTTP(w, request)

	var result Response
	_ = json.Unmarshal(w.Body.Bytes(), &result)

	// Check response data.
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, updateData.Name, result.Data.Name)
	assert.Equal(t, true, isEqualNodeData(updateData.Nodes, result.Data.Nodes))

	// Check database.
	assert.Equal(t, true, isEqualDB(channelResponseList))

	// Check configuration.
	conf := configuration.InitConfigData(confFilePath)
	assert.Equal(t, true, isEqualChannelConfigList(channelResponseList, conf.Channel))
}

// Test to return error when duplicate channel name.
func TestPutHandler2(t *testing.T) {
	Setup()
	defer Teardown()

	// Set update channel data.
	whereName := "channel1"
	updateChannelTB := db.GetConfigurationChannelInfoByName(whereName)
	whereID := updateChannelTB.CHANNEL_PK

	var updateData ChannelData
	updateData.Name = "channel3"
	updateData.ID = whereID
	updateData.Nodes = make([]NodeData, len(channelResponseList.Data[0].Nodes)-1)
	for i := 0; i < len(updateData.Nodes); i++ {
		updateData.Nodes[i] = channelResponseList.Data[0].Nodes[i+1]
	}
	for i, value := range channelResponseList.Data {
		if value.Name == whereName {
			channelResponseList.Data[i].Name = updateData.Name
			channelResponseList.Data[i].Nodes = updateData.Nodes
		}
	}

	// Set request data.
	var requestData Request
	requestData.Data.Name = updateData.Name
	requestData.Data.Nodes = make([]string, len(updateData.Nodes))
	for i, value := range updateData.Nodes {
		requestData.Data.Nodes[i] = value.ID
	}
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.PUT(constants.ChannelsPutAPIURL, PutHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPUT, constants.ChannelsAPIBaseURL+"/"+whereID, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	// Check duplicate node name.
	var err isaacerror.APIError
	_ = json.Unmarshal(w.Body.Bytes(), &err)

	assert.Equal(t, 409, w.Code)
	assert.Equal(t, isaacerror.ErrorDuplicatedChannelName, err.Errors[0].UserMessage)
}

func TestDeleteHandler(t *testing.T) {
	Setup()
	defer Teardown()

	// Set delete node data.
	deleteName := "channel1"
	deleteChannelTB := db.GetConfigurationChannelInfoByName(deleteName)
	deleteID := deleteChannelTB.CHANNEL_PK

	var deleteData ChannelData
	deleteData.Name = deleteName
	deleteData.ID = deleteChannelTB.CHANNEL_PK
	deleteChannelMappingTB := db.GetChannelPermissionNodes(deleteID)
	deleteData.Nodes = make([]NodeData, len(deleteChannelMappingTB))
	nodePKToData := db.GetNodePKToDataMap()
	for i, value := range deleteChannelMappingTB {
		deleteData.Nodes[i].ID = value.NODE_PK
		deleteData.Nodes[i].Name = nodePKToData[value.NODE_PK].NODE_NAME
		deleteData.Nodes[i].IP = nodePKToData[value.NODE_PK].NODE_IP
	}
	for i, value := range channelResponseList.Data {
		if value.Name == deleteName {
			copy(channelResponseList.Data[i:], channelResponseList.Data[i+1:])
			channelResponseList.Data[len(channelResponseList.Data)-1] = ChannelData{}
			channelResponseList.Data = channelResponseList.Data[:len(channelResponseList.Data)-1]
		}
	}
	channelResponseList.Total -= 1

	router := gin.Default()
	router.DELETE(constants.ChannelsDeleteAPIURL, DeleteHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodDELETE, constants.ChannelsAPIBaseURL+"/"+deleteID, nil)
	router.ServeHTTP(w, request)

	var result Response
	_ = json.Unmarshal(w.Body.Bytes(), &result)

	// Check response data.
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, deleteData.Name, result.Data.Name)
	assert.Equal(t, true, isEqualNodeData(deleteData.Nodes, result.Data.Nodes))

	// Check database.
	assert.Equal(t, true, isEqualDB(channelResponseList))
	assert.Equal(t, true, isEqualDeletedDB(channelResponseList))

	// Check configuration.
	conf := configuration.InitConfigData(confFilePath)
	assert.Equal(t, true, isEqualChannelConfigList(channelResponseList, conf.Channel))
}

func isEqualChannelResponseList(var1, var2 ResponseList) bool {
	var1Len := len(var1.Data)
	var2Len := len(var2.Data)

	// Check length.
	if var1Len != var2Len {
		return false
	}

	if var1.Total != var2.Total {
		return false
	}

	// Check channel name and ip, .
	for _, value1 := range var1.Data {
		check := false

		for _, value2 := range var2.Data {
			if value1.Name == value2.Name {
				// Check node name and ip.
				if !isEqualNodeData(value1.Nodes, value2.Nodes) {
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

func isEqualDB(var1 ResponseList) bool {
	channelTB := db.GetConfigurationChannelTable()
	alertTB := db.GetAlertConfigTable()
	visibilityTB := db.GetVisibilityConfigTable()

	var1Len := len(var1.Data)
	channelTBLen := len(channelTB)

	// Check length.
	if var1.Total != channelTBLen {
		return false
	}

	if var1Len != channelTBLen {
		return false
	}

	if var1Len != len(alertTB) {
		return false
	}

	if var1Len != len(visibilityTB) {
		return false
	}

	// Check channel name.
	for _, value1 := range var1.Data {
		check := false

		for _, value2 := range channelTB {
			if value1.Name == value2.CHANNEL_NAME {
				check = true
				break
			}
		}
		if !check {
			return false
		}
		check = false

		getChannelTB := db.GetConfigurationChannelInfoByName(value1.Name)

		for _, value2 := range alertTB {
			if getChannelTB.CHANNEL_PK == value2.CHANNEL_PK {
				check = true
				break
			}
		}
		if !check {
			return false
		}
		check = false

		for _, value2 := range visibilityTB {
			if getChannelTB.CHANNEL_PK == value2.CHANNEL_PK {
				check = true
				break
			}
		}
		if !check {
			return false
		}
	}

	nodeNameToPK := db.GetNodeNameToPKMap()

	// Check NODE_CHANNEL_MAPPING table.
	for _, value := range var1.Data {
		getChannelTB := db.GetConfigurationChannelInfoByName(value.Name)
		mappingTB := db.GetChannelPermissionNodes(getChannelTB.CHANNEL_PK)

		if len(value.Nodes) != len(mappingTB) {
			return false
		}

		for _, nodeValue := range value.Nodes {
			check := false

			for _, mappingValue := range mappingTB {
				if nodeNameToPK[nodeValue.Name] == mappingValue.NODE_PK {
					check = true
					break
				}
			}

			if !check {
				return false
			}
		}
	}

	return true
}

func isEqualDeletedDB(var1 ResponseList) bool {
	userMappingTB := db.GetUserPermissionChannelsTable()
	groupMappingTB := db.GetGroupPermissionChannelsTable()

	var1Len := len(var1.Data)

	var userMappingList []string
	for _, value := range userMappingTB {
		if !utility.IsExistValueInList(value.CHANNEL_PK, userMappingList) {
			userMappingList = append(userMappingList, value.CHANNEL_PK)
		}
	}

	if var1Len != len(userMappingList) {
		return false
	}

	if var1Len != len(groupMappingTB) {
		return false
	}

	// Check channel name.
	for _, value1 := range var1.Data {
		check := false

		getChannelTB := db.GetConfigurationChannelInfoByName(value1.Name)
		if !utility.IsExistValueInList(getChannelTB.CHANNEL_PK, userMappingList) {
			return false
		}

		for _, value2 := range groupMappingTB {
			if getChannelTB.CHANNEL_PK == value2.CHANNEL_PK {
				check = true
				break
			}
		}
		if !check {
			return false
		}
		check = false
	}

	return true
}

func isEqualChannelConfigList(var1 ResponseList, var2 []configuration.Channels) bool {
	var1Len := len(var1.Data)
	var2Len := len(var2)

	// Check length.
	if var1Len != var2Len {
		return false
	}

	if var1.Total != var2Len {
		return false
	}

	// Check channel name and ip.
	for _, value1 := range var1.Data {
		check := false

		for _, value2 := range var2 {
			if value1.Name == value2.Name {

				if !isEqualNodeDataConf(value1.Nodes, value2.Nodes) {
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

func isEqualNodeData(var1 []NodeData, var2 []NodeData) bool {
	var1Len := len(var1)
	var2Len := len(var2)

	// Check length.
	if var1Len != var2Len {
		return false
	}

	// Check node name and ip.
	for _, nodeValue1 := range var1 {
		check := false

		for _, nodeValue2 := range var2 {
			if nodeValue1.Name == nodeValue2.Name {
				if nodeValue1.IP != nodeValue2.IP {
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

func isEqualNodeDataConf(var1 []NodeData, var2 []string) bool {
	var1Len := len(var1)
	var2Len := len(var2)

	// Check length.
	if var1Len != var2Len {
		return false
	}

	// Check node name and ip.
	for _, nodeValue1 := range var1 {
		check := false

		for _, nodeValue2 := range var2 {
			if nodeValue1.Name == nodeValue2 {
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
	conf.Prometheus = prometheusTestData
	conf.ETC = etcTestData

	configuration.ChangeConfigFile(confFilePath, &conf)
	configuration.InitConfigData(confFilePath)

	// Create database.
	db.InitDB("sqlite3", dbPath)
	db.InitCreateTable()

	channelTB := db.GetConfigurationChannelTable()
	nodeTB := db.GetConfigurationNodeTable()

	channelResponseList.Total = len(channelTB)
	channelResponseList.Data = make([]ChannelData, channelResponseList.Total)
	for i, value := range channelTB {
		channelResponseList.Data[i].ID = value.CHANNEL_PK
		channelResponseList.Data[i].Name = value.CHANNEL_NAME
		channelResponseList.Data[i].Total = len(nodeTB)
		channelResponseList.Data[i].Nodes = make([]NodeData, channelResponseList.Data[i].Total)
		for j, nodeValue := range nodeTB {
			channelResponseList.Data[i].Nodes[j].ID = nodeValue.NODE_PK
			channelResponseList.Data[i].Nodes[j].Name = nodeValue.NODE_NAME
			channelResponseList.Data[i].Nodes[j].IP = nodeValue.NODE_IP
		}
	}

	// Insert channel data.
	groupMappingTB := make([]db.CHANNEL_GROUP_MAPPING_TB, len(channelTB))
	for i, value := range channelTB {
		db.InsertUserPermissionChannels("PKUS000000000000"+strconv.Itoa(i), value.CHANNEL_PK)
		groupMappingTB[i].CHANNEL_PK = value.CHANNEL_PK
		groupMappingTB[i].GROUP_PK = "PKGR000000000000" + strconv.Itoa(i)
		db.DBgorm().Create(&groupMappingTB[i])
	}
}

func Teardown() {
	_ = os.Remove(confFilePath)
}
