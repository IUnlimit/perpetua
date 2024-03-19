package perp

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestProxy(t *testing.T) {
	http.HandleFunc("/", proxyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type VerifyPost struct {
	Collect     string `json:"collect"`
	Tlg         int32  `json:"tlg"`
	Eks         string `json:"eks"`
	Sess        string `json:"sess"`
	Ans         string `json:"ans"`
	PowAnswer   string `json:"pow_answer"`
	PowCalcTime int32  `json:"pow_calc_time"`
}

type VerifyResponse struct {
	Randstr    string `json:"randstr"`
	Ticket     string `json:"ticket"`
	ErrorCode  string `json:"errorCode"`
	ErrMessage string `json:"errMessage"`
	Sess       string `json:"sess"`
}

// aaa.show()
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("拦截到请求：" + r.URL.String())
	// 将目标URL替换为实际的验证码服务地址
	target := "https://ti.qq.com/safe/tools/captcha/sms-verify-login?aid=2086582797&login_appid=1600001615&sid=5909489600706368179&uin=3012218237"
	//verify := "https://t.captcha.qq.com/cap_union_new_verify"

	proxyReq, err := http.NewRequest(r.Method, target, r.Body)
	if err != nil {
		log.Printf("Error creating proxy request: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 设置必要的请求头（如Host、Referer等）
	for k, v := range r.Header {
		proxyReq.Header[k] = v
	}
	proxyReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36") // 根据实际情况设置

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("Error during proxy request: %v", err)
		http.Error(w, "Failed to connect to the remote server", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// 将响应内容复制回客户端
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	w.Write(bodyBytes)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
