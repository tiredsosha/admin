package protocols

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/tiredsosha/admin/tools/logger"
)

type UserCommand struct {
	Req string
}

type Response struct {
	Rl0string string `xml:"rl0string"`
}

// func SendGet(url string) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		logger.Error.Println(err)
// 		return
// 	}
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		logger.Error.Println(err)
// 	}
// 	logger.Debug.Printf("responce - %q to post req from %q\n", string(body), url)
// }

func SendGet(url string, timeout int) {
	// Create an HTTP client with a timeout

	var timeoutDur time.Duration = time.Duration(timeout) * time.Second

	client := &http.Client{
		Timeout: timeoutDur, // Set your desired timeout duration
	}

	resp, err := client.Get(url)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	defer resp.Body.Close() // Make sure to close the response body

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Println(err)
	}
	logger.Debug.Printf("response - %q to get req from %q\n", string(body), url)
}

func SendPost(url string, reqData string, timeout int) {
	// Define the timeout duration
	var timeoutDur time.Duration = time.Duration(timeout) * time.Second

	// Create an HTTP client with timeout
	client := &http.Client{
		Timeout: timeoutDur,
	}

	// Prepare the JSON payload
	values := map[string]string{"command": reqData}
	jsonReq, _ := json.Marshal(values)
	responseBody := bytes.NewBuffer(jsonReq)

	// Create the POST request
	req, err := http.NewRequest("POST", url, responseBody)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Println(err)
	}
	logger.Debug.Printf("response - %q to post req from %q\n", string(body), url)
}

// func SendPost(url string, reqData string) {
// 	values := map[string]string{"command": reqData}
// 	jsonReq, _ := json.Marshal(values)
// 	responseBody := bytes.NewBuffer(jsonReq)
// 	resp, err := http.Post(url, "application/json", responseBody)
// 	if err != nil {
// 		logger.Error.Println(err)
// 		return
// 	}
// 	defer resp.Body.Close()
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		logger.Error.Println(err)
// 	}
// 	logger.Debug.Printf("responce - %q to post req from %q\n", string(body), url)
// }

func GetRelay(url string, timeout int) int {
	status := 521
	var timeoutDur time.Duration = time.Duration(timeout) * time.Second

	// Create an HTTP client with timeout
	client := &http.Client{
		Timeout: timeoutDur,
	}

	// Build the GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error.Println("error creating request:", err)
		return 520
	}

	// Send the GET request
	resp, err := client.Do(req)
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

	// Parse XML response
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
		status = 200
	}
	logger.Debug.Println("relay status -", value)

	return status
}

// func GetRelay(url string) int {
// 	status := 200

// 	// Send GET request
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		logger.Error.Println("error making GET request:", err)
// 		return 520
// 	}
// 	defer resp.Body.Close()
// 	// Read response body
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		logger.Error.Println("error reading response body:", err)
// 		return 520
// 	}
// 	// Parse XML
// 	var response Response
// 	err = xml.Unmarshal(body, &response)
// 	if err != nil {
// 		logger.Error.Println("error unmarshalling XML:", err)
// 		return 520
// 	}
// 	// Convert string to int
// 	value, err := strconv.Atoi(response.Rl0string)
// 	if err != nil {
// 		logger.Error.Println("error converting string to int:", err)
// 		return 520
// 	}
// 	if value == 1 {
// 		status = 521
// 	}
// 	logger.Debug.Println("relay status -", value)

// 	return status
// }

// func GetPC(url string) int {
// 	status := 200

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		logger.Error.Println(err)
// 		return 521
// 	}
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		logger.Error.Println(err)
// 		return 521
// 	}
// 	logger.Debug.Println("relay status -", body)
// 	return status

// }

func GetPC(url string, timeout int) int {
	status := 200
	var timeoutDur time.Duration = time.Duration(timeout) * time.Second

	// Create an HTTP client with timeout
	client := &http.Client{
		Timeout: timeoutDur,
	}

	// Build the GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error.Println(err)
		return 521
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		logger.Error.Println(err)
		return 521
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Println(err)
		return 521
	}

	logger.Debug.Println("relay status -", string(body))
	return status
}
