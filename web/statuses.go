package web

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tiredsosha/admin/tools/logger"
	"gopkg.in/yaml.v3"
)

func statusPark(c *gin.Context) {
	// Read the JSON data from the file
	jsonBytes, err := os.ReadFile("./configs/status.json")
	if err != nil {
		// Handle error: file not found or read error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read status data"})
		return
	}

	// Respond with the JSON data
	c.Data(http.StatusOK, "application/json", jsonBytes)

}

func convertYamlToJson(yamlData []byte) ([]byte, error) {
	var parsedData any
	// Unmarshal YAML into an interface{}
	if err := yaml.Unmarshal(yamlData, &parsedData); err != nil {
		return nil, err
	}
	// Marshal the data into JSON with indentation
	jsonData, err := json.MarshalIndent(parsedData, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func ConverterYaml() {
	// Read the YAML file
	yamlBytes, err := os.ReadFile("./configs/status.yaml")
	if err != nil {
		logger.Warn.Println(err)
		logger.Error.Fatal("EXITING")
	}

	// Convert YAML content to JSON
	jsonBytes, err := convertYamlToJson(yamlBytes)
	if err != nil {
		logger.Warn.Println(err)
		logger.Error.Fatal("EXITING")
	}

	// Write JSON to a file
	err = os.WriteFile("./configs/status.json", jsonBytes, 0644)
	if err != nil {
		logger.Warn.Println(err)
		logger.Error.Fatal("EXITING")
	}

	updateYaml(jsonBytes)

	logger.Warn.Println("yaml converted to json, check status.json")
}

func updateYaml(jBytes []byte) {
	var data any
	if err := json.Unmarshal(jBytes, &data); err != nil {
		logger.Warn.Println(err)
		logger.Error.Fatal("EXITING")
	}

	// Marshal the data into YAML
	yamlBytesOut, err := yaml.Marshal(data)
	if err != nil {
		logger.Warn.Println(err)
		logger.Error.Fatal("EXITING")
	}

	// Save the YAML to a file
	err = os.WriteFile("./configs/status.yaml", yamlBytesOut, 0644)
	if err != nil {
		logger.Warn.Println(err)
		logger.Error.Fatal("EXITING")
	}
}
