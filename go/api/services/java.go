package services

import (
	"Visma/models"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func SendCodeToJava(code models.TaskToSubmit, context *gin.Context) {
	jsonBytes, err := json.Marshal(code)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "bad code"})
		return
	}
	resp, err := http.Post(models.JavaURL, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		context.JSON(http.StatusInternalServerError, err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read response body",
		})
		return
	}
	var responseMap map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseMap)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse response body",
		})
		return
	}
	context.JSON(http.StatusOK, responseMap)
}
