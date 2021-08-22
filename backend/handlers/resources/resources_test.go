package resources

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
	"io/ioutil"
	"motherbear/backend/configuration"
	"motherbear/backend/constants"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var etcTestData = configuration.Etc{
	LoginLogoImagePath: "logoImage.png",
}

const confFilePath string = "testConfiguration.yaml"

func TestGetHandler(t *testing.T) {
	Setup()
	defer Teardown()

	// Test get channel list
	router := gin.Default()
	router.GET(constants.ResourcesGETAPIURL, GetHandler)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(constants.HTTPMethodGET, constants.ResourcesAPIBaseURL+"/"+constants.ResourcesIDLoginLogoImage, nil)
	router.ServeHTTP(w, request)

	var result Response
	_ = json.Unmarshal(w.Body.Bytes(), &result)

	assert.Equal(t, 200, w.Code)
}

func Setup() {
	var conf configuration.Configuration

	// Add configuration.
	conf.ETC = etcTestData

	configuration.ChangeConfigFile(confFilePath, &conf)
	configuration.InitConfigData(confFilePath)

	imageBase64 := "R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
	imageByte, _ := base64.StdEncoding.DecodeString(imageBase64)
	err := ioutil.WriteFile(conf.ETC.LoginLogoImagePath, imageByte, 0644)
	if err != nil {
		err.Error()
	}
}

func Teardown() {
	_ = os.Remove(confFilePath)
	_ = os.Remove(etcTestData.LoginLogoImagePath)
}
