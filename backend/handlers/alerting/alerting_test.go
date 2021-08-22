package alerting

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
	"motherbear/backend/configuration"
	"motherbear/backend/constants"
	"motherbear/backend/db"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var channelTestData = []configuration.Channels{
	{
		Name:  "channel1",
		Nodes: []string{"node1", "node2"},
	},
	{
		Name:  "channel2",
		Nodes: []string{"node1", "node2"},
	},
	{
		Name:  "channel3",
		Nodes: []string{"node1", "node2"},
	},
}

const dbPath string = ":memory:"
const confFilePath string = "testConfiguration.yaml"

func TestGetHandler(t *testing.T) {
	Setup()
	defer Teardown()

	// Get alerting data in database.
	alertTB := db.GetAlertConfigTable()

	// Test get alerting data.
	router := gin.Default()
	router.GET(constants.AlertingGetListAPIURL, GetHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodGET, constants.AlertingGetListAPIURL, nil)
	router.ServeHTTP(w, request)

	var resultList ResponseList
	_ = json.Unmarshal(w.Body.Bytes(), &resultList)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, strconv.Itoa(len(alertTB)), w.Header().Get("X-Total-Count"))
	assert.Equal(t, len(alertTB), len(resultList.Data))
	channelCount := 0
	for _, value := range alertTB {
		for _, result := range resultList.Data {
			if value.CHANNEL_PK == result.ID {
				assert.Equal(t, value.MAX_TIME_SEC_FOR_UNSYNC, result.UnsyncBlockToleranceTime)
				assert.Equal(t, value.MAX_TIME_SEC_FOR_RESPONSE, result.SlowResponseTime)
				channelCount++
			}
		}
	}
	assert.Equal(t, len(alertTB), channelCount)
}

func TestPutHandler(t *testing.T) {
	Setup()
	defer Teardown()

	var requestData Request

	unsyncTime1 := 360
	responseTime1 := 10
	unsyncTime2 := 420
	responseTime2 := 20

	// Get alerting data in database.
	alertTB := db.GetAlertConfigTable()

	// Get channels map that can be convert channel pk to channel name.
	channelPKToName := db.GetChannelPKToNameMap()

	requestData.Data = make([]AlertingData, len(alertTB))
	requestData.Total = len(alertTB)

	for i, value := range alertTB {
		requestData.Data[i].ID = value.CHANNEL_PK
		requestData.Data[i].Name = channelPKToName[value.CHANNEL_PK]

		if requestData.Data[i].Name == "channel1" {
			requestData.Data[i].UnsyncBlockToleranceTime = unsyncTime1
			requestData.Data[i].SlowResponseTime = responseTime1
		} else if requestData.Data[i].Name == "channel2" {
			requestData.Data[i].UnsyncBlockToleranceTime = unsyncTime2
			requestData.Data[i].SlowResponseTime = responseTime2
		} else {
			requestData.Data[i].UnsyncBlockToleranceTime = value.MAX_TIME_SEC_FOR_UNSYNC
			requestData.Data[i].SlowResponseTime = value.MAX_TIME_SEC_FOR_RESPONSE
		}
	}
	requestJSON, _ := json.Marshal(requestData)

	// Test get alerting data.
	router := gin.Default()
	router.PUT(constants.AlertingGetListAPIURL, PutHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPUT, constants.AlertingGetListAPIURL, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	// Get alerting data in database.
	alertTB = db.GetAlertConfigTable()

	var resultList ResponseList
	_ = json.Unmarshal(w.Body.Bytes(), &resultList)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, strconv.Itoa(len(alertTB)), w.Header().Get("X-Total-Count"))
	assert.Equal(t, len(alertTB), len(resultList.Data))
	channelCount := 0
	for _, value := range alertTB {
		for _, result := range resultList.Data {
			if value.CHANNEL_PK == result.ID {
				if result.Name == "channel1" {
					assert.Equal(t, value.MAX_TIME_SEC_FOR_UNSYNC, unsyncTime1)
					assert.Equal(t, value.MAX_TIME_SEC_FOR_RESPONSE, responseTime1)
				} else if result.Name == "channel2" {
					assert.Equal(t, value.MAX_TIME_SEC_FOR_UNSYNC, unsyncTime2)
					assert.Equal(t, value.MAX_TIME_SEC_FOR_RESPONSE, responseTime2)
				} else {
					assert.Equal(t, value.MAX_TIME_SEC_FOR_UNSYNC, result.UnsyncBlockToleranceTime)
					assert.Equal(t, value.MAX_TIME_SEC_FOR_RESPONSE, result.SlowResponseTime)
				}
				channelCount++
			}
		}
	}
	assert.Equal(t, len(alertTB), channelCount)
}

func Setup() {
	var conf configuration.Configuration

	// Add channel configuration.
	conf.Channel = channelTestData

	configuration.ChangeConfigFile(confFilePath, &conf)
	configuration.InitConfigData(confFilePath)

	// Create database.
	db.InitDB("sqlite3", dbPath)
	db.InitCreateTable()
}

func Teardown() {
	_ = os.Remove(confFilePath)
}
