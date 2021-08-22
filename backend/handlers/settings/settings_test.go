package settings

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
	"motherbear/backend/configuration"
	"motherbear/backend/constants"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const confFilePath string = "testConfiguration.yaml"

var settingsTestData = configuration.Etc{
	SessionTimeout: 30,
}

func TestGetHandler(t *testing.T) {
	Setup()
	defer Teardown()

	// Test get configuration value.
	router := gin.Default()
	router.GET(constants.SettingsGetListAPIURL, GetHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodGET, constants.SettingsGetListAPIURL, nil)
	router.ServeHTTP(w, request)

	var result Response
	_ = json.Unmarshal(w.Body.Bytes(), &result)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, settingsTestData.SessionTimeout, result.Data.SessionTimeout)
}

func TestPutHandler(t *testing.T) {
	Setup()
	defer Teardown()

	// Set update configuration value.
	updateSessionTimeOut := 60

	// Set request data.
	var requestData Request
	requestData.Data.SessionTimeout = updateSessionTimeOut
	requestJSON, _ := json.Marshal(requestData)

	router := gin.Default()
	router.PUT(constants.SettingsPutAPIURL, PutHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodPUT, constants.SettingsPutAPIURL, bytes.NewBuffer(requestJSON))
	router.ServeHTTP(w, request)

	var result Response
	_ = json.Unmarshal(w.Body.Bytes(), &result)

	// Check response data.
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, requestData.Data.SessionTimeout, result.Data.SessionTimeout)
}

func Setup() {
	// Add configuration value.
	var conf configuration.Configuration
	conf.ETC = settingsTestData

	configuration.ChangeConfigFile(confFilePath, &conf)
	configuration.InitConfigData(confFilePath)
}

func Teardown() {
	_ = os.Remove(confFilePath)
}
