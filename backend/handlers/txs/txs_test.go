package txs

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
	"motherbear/backend/constants"
	"motherbear/backend/isaacerror"
	"motherbear/backend/polarbear"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func generateTestBlockdataInDB(channelName string, height int) error {
	for h := 1; h <= height; h++ {
		var data map[string]interface{}
		var blockRecord polarbear.Block

		_, err := polarbear.GenerateTestBlockJSONData(3, int64(h), &data)
		if err != nil {
			return err
		}
		if err := polarbear.AddBlockRecordFromJSONResponse(
			data,
			"",
			channelName,
			&blockRecord); err != nil {
			return err
		}
	}

	return nil
}

func setup() {
	// Init polarbear module.
	_ = polarbear.Init(constants.DBTypeSqlite3, ":memory:")

	// Generate pseudo test data in each block data.
	_ = generateTestBlockdataInDB("channel1", 250)
	_ = generateTestBlockdataInDB("channel2", 100)
	_ = generateTestBlockdataInDB("channel3", 120)

}

func tearDown() {
}

func TestGetHandlerList(t *testing.T) {
	setup()
	defer tearDown()

	//Check the current test data.
	currentHeight := polarbear.GetCurrentBlockHeightInDB("channel1")
	assert.Equal(t, currentHeight, int64(250))

	// Prepare Block GET-LIST handler.
	router := gin.Default()
	router.GET(constants.TxGETListAPIURL, GetHandlerList)
	router.GET(constants.TxGETAPIURL, GetHandler)
	w := httptest.NewRecorder()

	// Request the URL.
	request, _ := http.NewRequest(
		constants.HTTPMethodGET,
		constants.TxGETListAPIURL,
		nil)
	q := request.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "10")
	q.Add("channel", "channel1")
	request.URL.RawQuery = q.Encode()
	router.ServeHTTP(w, request)

	var resultList TxResponseList
	_ = json.Unmarshal(w.Body.Bytes(), &resultList)

	// Test the response.
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, len(resultList.Data), 10)
	assert.NotEqual(t, resultList.Total, 10)
	assert.Equal(t, strconv.Itoa(resultList.Total), w.Header().Get("X-Total-Count"))

	// Test the tx response.
	// Try to query by block hash.
	w2 := httptest.NewRecorder()
	txHash := resultList.Data[2].TxHash
	r2, _ := http.NewRequest(
		constants.HTTPMethodGET,
		constants.TxAPIBaseURL+"/"+txHash,
		nil)
	q2 := r2.URL.Query()
	q2.Add("channel", "channel1")
	r2.URL.RawQuery = q2.Encode()
	router.ServeHTTP(w2, r2)

	var resultBody TxResponse
	_ = json.Unmarshal(w2.Body.Bytes(), &resultBody)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, txHash, resultBody.TxHash)
	assert.Equal(t, resultList.Data[2].Data, resultBody.Data)
}

// Test to get transaction list to meet the timestamp in search items.
func TestGetHandlerList2(t *testing.T) {
	setup()
	defer tearDown()

	//Check the current test data.
	currentHeight := polarbear.GetCurrentBlockHeightInDB("channel1")
	assert.Equal(t, currentHeight, int64(250))

	// Prepare Block GET-LIST handler.
	router := gin.Default()
	router.GET(constants.TxGETListAPIURL, GetHandlerList)
	router.GET(constants.TxGETAPIURL, GetHandler)
	w := httptest.NewRecorder()

	// Request the URL.
	request, _ := http.NewRequest(
		constants.HTTPMethodGET,
		constants.TxGETListAPIURL,
		nil)

	now := time.Now().UTC()
	nowString := now.Format(time.RFC3339)

	q := request.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "10")
	q.Add("channel", "channel1")
	q.Add("from", nowString)
	q.Add("to", nowString)
	request.URL.RawQuery = q.Encode()
	router.ServeHTTP(w, request)

	var resultList TxResponseList
	_ = json.Unmarshal(w.Body.Bytes(), &resultList)

	// Test the response.
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, len(resultList.Data), 0)
	assert.Equal(t, resultList.Total, 0)
	assert.Equal(t, strconv.Itoa(resultList.Total), w.Header().Get("X-Total-Count"))
}

// Test to return error when send only 'from' in search items.
func TestGetHandlerList3(t *testing.T) {
	setup()
	defer tearDown()

	//Check the current test data.
	currentHeight := polarbear.GetCurrentBlockHeightInDB("channel1")
	assert.Equal(t, currentHeight, int64(250))

	// Prepare Block GET-LIST handler.
	router := gin.Default()
	router.GET(constants.TxGETListAPIURL, GetHandlerList)
	router.GET(constants.TxGETAPIURL, GetHandler)
	w := httptest.NewRecorder()

	// Request the URL.
	request, _ := http.NewRequest(
		constants.HTTPMethodGET,
		constants.TxGETListAPIURL,
		nil)

	now := time.Now().UTC()
	nowString := now.Format(time.RFC3339)

	q := request.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "10")
	q.Add("channel", "channel1")
	q.Add("from", nowString)
	request.URL.RawQuery = q.Encode()
	router.ServeHTTP(w, request)

	var err isaacerror.APIError
	_ = json.Unmarshal(w.Body.Bytes(), &err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, isaacerror.ErrorFailToQueryTxList, err.Errors[0].UserMessage)
}
