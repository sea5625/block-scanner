package symptom

import (
	"gopkg.in/go-playground/assert.v1"
	"motherbear/backend/polarbear"
	"testing"
	"time"
)

var symptomTestData = []polarbear.Symptom{
	{
		Channel: "loopchain_default1",
		Msg: "[node4]response time slowly [25.106827] sec",
		SymptomType: "Slow response",
		Timestamp: time.Now(),
	},
	{
		Channel: "loopchain_default2",
		Msg: "[node4]response time slowly [30.106827] sec",
		SymptomType: "Slow response",
		Timestamp: time.Now(),
	},
	{
		Channel: "loopchain_default3",
		Msg: "[node4]response time slowly [35.106827] sec",
		SymptomType: "Slow response",
		Timestamp: time.Now(),
	},
	{
		Channel: "loopchain_default2",
		Msg: "[node4]response time slowly [45.106827] sec",
		SymptomType: "Slow response",
		Timestamp: time.Now(),
	},
	{
		Channel: "loopchain_default3",
		Msg: "[node4]response time slowly [55.106827] sec",
		SymptomType: "Slow response",
		Timestamp: time.Now(),
	},
	{
		Channel: "loopchain_default3",
		Msg: "[node4]response time slowly [65.106827] sec",
		SymptomType: "Slow response",
		Timestamp: time.Now(),
	},
	{
		Channel: "loopchain_default1",
		Msg: "[node4]response time slowly [75.106827] sec",
		SymptomType: "Slow response",
		Timestamp: time.Now(),
	},
}

func TestGetHandler(t *testing.T) {
	// Query data.
	var resp PeerResponseList
	resp.Data = make([]PeerSymptomResponse, len(symptomTestData))
	resp.Total = len(symptomTestData)

	for i := 0; i < len(symptomTestData); i++ {
		convertPeerSymptomResponse(&symptomTestData[i], &resp.Data[i])
	}

	assert.Equal(t, len(symptomTestData), len(resp.Data))
}