package nodetype

import (
	"github.com/gin-gonic/gin"
	"motherbear/backend/configuration"
	"motherbear/backend/constants"
	"net/http"
	"strconv"
)

// TypeResponse is the response for GET (node type)
type TypeResponse struct {
	Data DataResponse `json:"data"`
}

// DataResponse data is the response for GET
type DataResponse struct {
	NodeType   string `json:"nodeType,omitempty" example:"loopchain"`
	Prometheus string `json:"prometheus,omitempty" example:"127.0.0.1:9090/api/v1/query"`
	JobName    string `json:"jobName,omitempty" example:"goloop"`
}

// GetHandler godoc
// @Tags Prometheus
// @Summary GET handler of prometheus
// @Description Get node type and prometheus IP + Query API.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication header ('Bearer '+ JWT token)"
// @Success 200 {object} nodetype.TypeResponse "Result for node type and prometheus IP + Query API"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /prometheus [get]
func GetHandler(c *gin.Context) {

	var resp TypeResponse

	resp.Data.NodeType = configuration.QueryNodeType()
	resp.Data.Prometheus = configuration.QueryPrometheusIP()
	if resp.Data.NodeType == constants.NodeType2 {
		resp.Data.JobName = configuration.Conf().Prometheus.JobNameOfgoloop
	}

	// Put total information.
	bytes := "bytes 0-" + strconv.Itoa(1) + "/" + strconv.Itoa(1)
	c.Header(constants.HTTPHeaderContentRange, bytes)
	c.Header(constants.HTTPHeaderXTotalCount, strconv.Itoa(1))

	// Return body.
	c.JSON(http.StatusOK, resp)

}
