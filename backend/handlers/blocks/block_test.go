package blocks

import (
	"encoding/json"
	"motherbear/backend/constants"
	"motherbear/backend/polarbear"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
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

func TestGetListHandler(t *testing.T) {
	setup()
	defer tearDown()

	//Check the current test data.
	currentHeight := polarbear.GetCurrentBlockHeightInDB("channel1")
	assert.Equal(t, currentHeight, int64(250))

	// Prepare Block GET-LIST handler.
	router := gin.Default()
	router.GET(constants.BlockGETListAPIURL, GetHandlerList)
	router.GET(constants.BlockGETAPIURL, GetHandler)
	w := httptest.NewRecorder()

	// Request the URL.
	request, _ := http.NewRequest(
		constants.HTTPMethodGET,
		constants.BlockGETListAPIURL,
		nil)
	q := request.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "10")
	q.Add("channel", "channel1")
	request.URL.RawQuery = q.Encode()
	router.ServeHTTP(w, request)

	var resultList BlockResponseList
	_ = json.Unmarshal(w.Body.Bytes(), &resultList)

	// Test the response.
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, len(resultList.Data), 10)
	assert.NotEqual(t, resultList.Total, 10)
	assert.Equal(t, len(resultList.Data[0].ConfirmedTx), 3)
	assert.Equal(t, resultList.Data[0].BlockHeight, int64(240))
	assert.Equal(t, strconv.Itoa(resultList.Total), w.Header().Get("X-Total-Count"))

	// Try to query by block hash.
	w2 := httptest.NewRecorder()
	blockHash := resultList.Data[2].BlockHash
	r2, _ := http.NewRequest(
		constants.HTTPMethodGET,
		constants.BlockAPIBaseURL+"/"+blockHash,
		nil)
	q2 := r2.URL.Query()
	q2.Add("channel", "channel1")
	r2.URL.RawQuery = q2.Encode()
	router.ServeHTTP(w2, r2)
	var resultBody BlockResponse
	_ = json.Unmarshal(w2.Body.Bytes(), &resultBody)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, blockHash, resultBody.BlockHash)
	assert.Equal(t, resultList.Data[2].ConfirmedTx[0].Data, resultBody.ConfirmedTx[0].Data)

	// Try to query by block height.
	w3 := httptest.NewRecorder()
	blockHeight := resultList.Data[2].BlockHeight
	r3, _ := http.NewRequest(
		constants.HTTPMethodGET,
		constants.BlockAPIBaseURL+"/"+strconv.FormatInt(blockHeight, 10),
		nil)
	q3 := r3.URL.Query()
	q3.Add("channel", "channel1")
	r3.URL.RawQuery = q3.Encode()
	router.ServeHTTP(w3, r3)
	var resultBody2 BlockResponse
	_ = json.Unmarshal(w3.Body.Bytes(), &resultBody2)

	assert.Equal(t, http.StatusOK, w3.Code)
	assert.Equal(t, blockHeight, resultBody2.BlockHeight)
	assert.Equal(t, blockHash, resultBody2.BlockHash)
	assert.Equal(t, resultList.Data[2].ConfirmedTx[0].Data, resultBody2.ConfirmedTx[0].Data)
}
