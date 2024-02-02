package utils

import (
	"encoding/json"
	"fmt"
	"github.com/IUnlimit/perpetua/internal/erren"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

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
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Unexpected state code(:%d) with message: %s", resp.StatusCode, string(body)))
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

func CheckPort(host string, port int, timeout time.Duration) error {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

// BadResponse Return error status code and error message
func BadResponse(c *gin.Context, err error) {
	Err := erren.ConvertErr(err)
	c.JSON(http.StatusOK, model.Response{
		Status:  "failed",
		RetCode: Err.ErrCode,
		Msg:     Err.ErrMsg,
	})
}

// SendResponse Try to return data
func SendResponse(c *gin.Context, data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		BadResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  "ok",
		RetCode: 0,
		Data:    string(bytes),
	})
}

func download(resp *http.Response, filePath string, fileSize int64) error {
	if fileSize == -1 {
		fileSize = resp.ContentLength
	}
	log.Debug("Download content length: %d", fileSize)

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
