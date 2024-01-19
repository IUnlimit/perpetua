package utils

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/cheggaaa/pb/v3"
)

func BuildURLParams(baseURL string, params map[string]string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	query := u.Query()
	for key, value := range params {
		query.Set(key, value)
	}

	u.RawQuery = query.Encode()
	return u.String(), nil
}

func GetJson(url string, headers map[string]string, v any) error {
	req, _ := http.NewRequest("GET", url, nil)
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}
	return nil
}

func DownloadFile(url string, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = download(resp, filePath, -1)
	if err != nil {
		return err
	}
	return nil
}

func DownloadFileWithHeaders(url string, filePath string, headers map[string]string, fileSize int64) error {
	req, _ := http.NewRequest("GET", url, nil)
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	err = download(resp, filePath, fileSize)
	if err != nil {
		return err
	}
	return nil
}

func download(resp *http.Response, filePath string, fileSize int64) error {
	if fileSize == -1 {
		fileSize = resp.ContentLength
	}
	log.Debug("contentLength: %dB", fileSize)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	bar := pb.Full.Start64(fileSize)
	bar.Set(pb.Bytes, true)

	reader := bar.NewProxyReader(resp.Body)
	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	defer bar.Finish()
	return nil
}
