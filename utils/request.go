package utils

// 封装 web 请求相关的工具函数
import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// 封装 web 请求相关的工具函数
func Request(url string, method string, body any, headers map[string]string) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		switch b := body.(type) {
		case []byte:
			reqBody = bytes.NewReader(b)
		case string:
			reqBody = bytes.NewReader([]byte(b))
		default:
			jsonBytes, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			reqBody = bytes.NewReader(jsonBytes)
		}
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	// headers 允许为 nil，内部自动处理
	if headers == nil {
		headers = map[string]string{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36",
		}
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	// 默认 Content-Type
	if req.Header.Get("Content-Type") == "" && body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return respBytes, &HttpStatusError{StatusCode: resp.StatusCode, Body: respBytes}
	}
	return respBytes, nil
}

// HttpStatusError 用于返回非 2xx 状态码的错误
type HttpStatusError struct {
	StatusCode int
	Body       []byte
}

func (e *HttpStatusError) Error() string {
	return http.StatusText(e.StatusCode)
}
