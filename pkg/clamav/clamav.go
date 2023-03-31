package clamav

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type Response struct {
	Status      string `json:"status"`
	Description string `json:"description"`
}

const (
	host = "localhost"
	port = "9000"
)

func Interact(filePath string, logger *zap.Logger) {

	logger.Info("scanning " + filePath)

	url := fmt.Sprintf("http://%s:%s/scan", host, port)

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	file, errFile1 := os.Open(filePath)
	if errFile1 != nil {
		logger.Warn("error: " + errFile1.Error())
		return
	}

	defer file.Close()

	part1, errFile1 := writer.CreateFormFile("file", filepath.Base(filePath))
	if errFile1 != nil {
		logger.Warn("error: " + errFile1.Error())
		return
	}

	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		logger.Warn("error: " + errFile1.Error())
		return
	}

	err := writer.Close()
	if err != nil {
		logger.Warn("error: " + errFile1.Error())
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		logger.Warn("error: " + errFile1.Error())
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Warn("error: " + errFile1.Error())
		return
	}

	logger.Info(string(body))

}
