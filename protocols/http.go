package protocols

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/tiredsosha/admin/tools/logger"
)

type User struct {
	req string
}

func SendGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error.Println(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Println(err)
	}
	logger.Debug.Printf("responce - %q to post req from %q\n", string(body), url)
}

func SendPost(reqData User, url string) {
	postBody, _ := json.Marshal(reqData)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		logger.Error.Println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Println(err)
	}
	logger.Debug.Printf("responce - %q to post req from %q\n", string(body), url)
}
