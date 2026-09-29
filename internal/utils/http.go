package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/IUnlimit/perpetua/internal/erren"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/cheggaaa/pb/v3"
)

func ParseWebSocketURL(wsURL string) (host string, port int, suffix string, err error) {
	u, err := url.Parse(wsURL)
	if err != nil {
		return "", 0, "", err
	}

	hostname, portStr, err := net.SplitHostPort(u.Host)
	if err != nil {
		// 没有显式端口，则判断默认端口
		hostname = u.Host
		switch u.Scheme {
		case "ws":
			portStr = "80"
		case "wss":
			portStr = "443"
		default:
			portStr = ""
		}
	}

	suffix = u.Path
	if !strings.HasPrefix(suffix, "/") {
		suffix = "/" + suffix
	}
	port, err = strconv.Atoi(portStr)
	if err != nil {
		return "", 0, "", err
	}

	return hostname, port, suffix, nil
}

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
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

// CheckWebsocket checks whether the websocket endpoint accepts a handshake with the given access token
func CheckWebsocket(ws string, accessToken string, timeout time.Duration) error {
	conn, err := DialWebsocket(ws, accessToken, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

// DialWebsocket dials a OneBot forward websocket, carrying the access token (if any) as a Bearer header.
// Handshake failures are enriched with the HTTP status and response body so that auth/path problems
// (e.g. SnowLuma's 401 Unauthorized / 400 Bad path) can be told apart.
func DialWebsocket(wsUrl string, accessToken string, timeout time.Duration) (*websocket.Conn, error) {
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: timeout,
	}

	header := http.Header{}
	if len(accessToken) != 0 {
		header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	conn, resp, err := dialer.Dial(wsUrl, header)
	if err != nil {
		return nil, describeHandshakeError(err, resp)
	}
	return conn, nil
}

// describeHandshakeError appends status code, body and a troubleshooting hint to a bad handshake error
func describeHandshakeError(err error, resp *http.Response) error {
	if !errors.Is(err, websocket.ErrBadHandshake) || resp == nil {
		return err
	}

	var body string
	if resp.Body != nil {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		body = strings.TrimSpace(string(data))
	}

	var hint string
	switch resp.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		hint = ", please check `external-access-token`"
	case http.StatusBadRequest, http.StatusNotFound:
		hint = ", please check the path of `external-web-socket` (SnowLuma default: /)"
	}
	return fmt.Errorf("%w (status: %s, body: %q%s)", err, resp.Status, body, hint)
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

// GoodResponse Try to return data
// entry: mapKey1, mapValue1 ...
func GoodResponse(c *gin.Context, entry ...any) {
	if len(entry)%2 != 0 {
		BadResponse(c, errors.New(fmt.Sprintf("错误的 map 参数个数: %d", len(entry))))
		return
	}
	m := make(map[string]any)
	for i := 0; i < len(entry); i += 2 {
		m[entry[i].(string)] = entry[i+1]
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  "ok",
		RetCode: 0,
		Data:    m,
	})
}

func GoodResponseArray(c *gin.Context, array any) {
	c.JSON(http.StatusOK, model.Response{
		Status:  "ok",
		RetCode: 0,
		Data:    array,
	})
}

func download(resp *http.Response, filePath string, fileSize int64) error {
	if fileSize == -1 {
		fileSize = resp.ContentLength
	}
	log.Debugf("Download content length: %d", fileSize)

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
