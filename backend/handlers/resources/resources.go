package resources

import (
	"bufio"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"motherbear/backend/configuration"
	"motherbear/backend/constants"
	"motherbear/backend/isaacerror"
	"motherbear/backend/logger"
	"motherbear/backend/utility"
	"net/http"
	"os"
	"path"
)

type Response struct {
	Data string `json:"data"`
}

var allowImageFileExt = []string{"gif", "jpeg", "jpg", "png", "svg"}
var imageTypeMapping = map[string]string{
	allowImageFileExt[0]: "image/gif",
	allowImageFileExt[1]: "image/jpeg",
	allowImageFileExt[2]: "image/jpeg",
	allowImageFileExt[3]: "image/png",
	allowImageFileExt[4]: "image/svg+xml",
}

// GetHandlerList godoc
// @Tags resources
// @Summary GET handler of resource value in server.
// @Description Get resource value in server.
// @Accept  json
// @Produce  json
// @Param id path string true "Resource value ID to be get - Allow to get Resource ID : loginLogoImage"
// @Success 200 {object} resources.Response "Result for get resource value"
// @Failure 400 {object} isaacerror.APIError "Invalid parameter."
// @Failure 401 {object} isaacerror.APIError "Unauthorized"
// @Failure 500 {object} isaacerror.APIError "Internal server error."
// @Router /resources/{id} [get]
func GetHandler(c *gin.Context) {
	id := c.Param(constants.RequestResourceID)

	var response Response

	switch id {
	case constants.ResourcesIDLoginLogoImage:
		// Get logo image.
		imageFilePath := configuration.Conf().ETC.LoginLogoImagePath
		imageFileExtTemp := path.Ext(imageFilePath)
		if len(imageFileExtTemp) <= 1 {
			// Invalid image file extension.
			internalError := isaacerror.SysErrInvalidImageFileExtension.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorFailToGetLoginLogoImage, internalError)
			c.JSON(http.StatusInternalServerError, message)
			return
		}

		imageFileExt := imageFileExtTemp[1:]
		if !utility.IsExistValueInList(imageFileExt, allowImageFileExt) {
			// Not supported image file extension.
			internalError := isaacerror.SysErrNotSupportedImageFileExtension.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorFailToGetLoginLogoImage, internalError)
			c.JSON(http.StatusInternalServerError, message)
			return
		}

		imageFile, err := os.Open(imageFilePath)
		if err != nil {
			// File to get image file.
			internalError := isaacerror.SysErrFailToGetLoginLogoImage.Error()
			logger.Error(internalError)
			message := isaacerror.GetAPIError(isaacerror.ErrorFailToGetLoginLogoImage, internalError)
			c.JSON(http.StatusInternalServerError, message)
			return
		}

		defer imageFile.Close()

		// Convert image to byte
		fileInfo, err := imageFile.Stat()
		imageSize := fileInfo.Size()
		imageBuf := make([]byte, imageSize)

		fileReader := bufio.NewReader(imageFile)
		fileReader.Read(imageBuf)

		// Encode the image byte by base64.
		imageBase64String := base64.StdEncoding.EncodeToString(imageBuf)

		// Set image type.
		imageType := imageTypeMapping[imageFileExt]

		response.Data = "data:" + imageType + ";base64, " + imageBase64String

	default:
		// Not supported resource id.
		internalError := isaacerror.SysErrNotSupportedDataName.Error()
		logger.Error(internalError)
		message := isaacerror.GetAPIError(isaacerror.ErrorNotSupportedDataName, internalError)
		c.JSON(http.StatusBadRequest, message)
		return
	}

	c.JSON(http.StatusOK, response)
}
