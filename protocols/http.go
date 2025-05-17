package protocols

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"strconv"

	"github.com/tiredsosha/admin/tools/logger"
)

type UserCommand struct {
	Req string
}

type Response struct {
	Rl0string string `xml:"rl0string"`
}

func SendGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Println(err)
	}
	logger.Debug.Printf("responce - %q to post req from %q\n", string(body), url)
}

func SendPost(url string, reqData string) {
	values := map[string]string{"command": reqData}
	jsonReq, _ := json.Marshal(values)
	responseBody := bytes.NewBuffer(jsonReq)
	resp, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Println(err)
	}
	logger.Debug.Printf("responce - %q to post req from %q\n", string(body), url)
}

func GetRelay(url string) int {
	status := 200

	// Send GET request
	resp, err := http.Get(url)
	if err != nil {
		logger.Error.Println("error making GET request:", err)
		return 520
	}
	defer resp.Body.Close()
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Println("error reading response body:", err)
		return 520
	}
	// Parse XML
	var response Response
	err = xml.Unmarshal(body, &response)
	if err != nil {
		logger.Error.Println("error unmarshalling XML:", err)
		return 520
	}
	// Convert string to int
	value, err := strconv.Atoi(response.Rl0string)
	if err != nil {
		logger.Error.Println("error converting string to int:", err)
		return 520
	}
	if value == 1 {
		status = 521
	}
	logger.Debug.Println("relay status -", value)

	return status
}

func GetPC(url string) int {
	status := 200

	resp, err := http.Get(url)
	if err != nil {
		logger.Error.Println(err)
		return 521
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Println(err)
		return 521
	}
	logger.Debug.Println("relay status -", body)
	return status

}
