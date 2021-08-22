package symptom

import (
	"github.com/gin-gonic/gin"
	"motherbear/backend/constants"
	"motherbear/backend/db"
	"motherbear/backend/isaacerror"
	"motherbear/backend/logger"
	"motherbear/backend/polarbear"
	"motherbear/backend/utility"
	"net/http"
	"strconv"
	"time"
)

// BlockResponseList is the response for GET LIST
type PeerResponseList struct {
	Data  []PeerSymptomResponse `json:"data"`
	Total int             		`json:"total" example:"1" format:"int32"`
}

//BlockResponse is the response for GET
type PeerSymptomResponse struct {
	Channel   string           	`json:"channel,omitempty" example:"loopchain_default"`
	Msg   string           		`json:"msg,omitempty" example:"[node4]response time slowly [25.106827] sec"`
	Symptom string         		`json:"symptom,omitempty" example:"Slow response"`
	TimeStamp      string       `json:"timeStamp,omitempty" timeStamp:"2019-08-06T19:29:17+09:00"`
}

func convertPeerSymptomResponse(symptom *polarbear.Symptom, out *PeerSymptomResponse) {
	out.Channel = symptom.Channel
	out.Msg = symptom.Msg
	out.Symptom = symptom.SymptomType
	out.TimeStamp = symptom.Timestamp.Format(time.RFC3339)
	return
}

// GetHandlerList godoc
// @Tags Symptom
// @Summary GET handler of symptom
// @Description Get many peer symptom.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Param limit query integer true "Identify the number of results returned from a result set."
// @Param offset query integer  true "Identify the starting point to return data from a result set."
// @Param from query timeStamp false "Start date range for date range search."
// @Param to query timeStamp  false "End date range for date range search."
// @Success 200 {object} symptom.PeerResponseList "Result for many peer symptom"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /symptom [get]
func GetHandlerList(c *gin.Context) {

	// Check the parameters.
	offset, limit, from, to, err := utility.GetOffsetListDateSearchFromRequest(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	logger.Infof("peer symptom requested, limit:%d, offset:%d, from:%s, to:%s", limit, offset, from, to)

	// Get all channel resources.
	var permissionChannelListString []string
	permissionChannelList, exists := c.Get(constants.ContextKeyPermissionChannelList)
	if exists {
		// If permission channel list received from middleware, Channel not in permission channel list is unauthorized channel.
		// If channel is not in permission channel list, exclude.
		permissionChannelListString = utility.ConvertInterfaceToStringSlice(permissionChannelList)
	} else {
		channelTB := db.GetUserPermissionChannels(constants.AdminPK)
		permissionAdminChannelList := make([]string, 0)

		for _, value := range channelTB {
			permissionAdminChannelList = append(permissionAdminChannelList, value.CHANNEL_PK)
		}
		permissionChannelListString = utility.ConvertInterfaceToStringSlice(permissionAdminChannelList)
	}

	// Query data.
	var peerSymptom []polarbear.Symptom
	var count int64
	count, err = polarbear.QueryPeerSymptomListTable(limit, offset, from, to, permissionChannelListString, &peerSymptom)
	if err != nil {
		internalError := isaacerror.SysErrFailToQueryPeerSymptom.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorFailToQueryPeerSymptom, internalError)
		c.JSON(http.StatusInternalServerError, message)
		return
	}

	var resp PeerResponseList
	resp.Data = make([]PeerSymptomResponse, len(peerSymptom))
	resp.Total = int(count)

	for i := 0; i < len(peerSymptom); i++ {
		convertPeerSymptomResponse(&peerSymptom[i], &resp.Data[i])
	}

	// Put total information.
	bytes := "bytes 0-" + strconv.Itoa(resp.Total) + "/" + strconv.Itoa(resp.Total)
	c.Header(constants.HTTPHeaderContentRange, bytes)
	c.Header(constants.HTTPHeaderXTotalCount, strconv.Itoa(resp.Total))

	// Return body.
	c.JSON(http.StatusOK, resp)
}