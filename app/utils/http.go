package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// client 是一个全局的 HTTP 客户端，设置了超时和一个不验证 TLS 证书的传输层配置。
// 注意：InsecureSkipVerify 为 true 可能存在安全风险，建议在生产环境中谨慎使用。
var client = http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// DoHttpPost 是一个泛型函数，发送 HTTP POST 请求，并将响应体解析进泛型变量 t 中。
// 参数:
//
//	url: 请求地址
//	headerMap: 请求头的 key/value 映射
//	body: 请求体
//
// 返回值：泛型数据和可能发生的错误
func DoHttpPost[T interface{}](url string, headerMap map[string]string, body io.Reader) (T, error) {
	// 创建一个新的 POST 请求
	req, err := http.NewRequest("POST", url, body)
	var t T // 用于存储解析后的响应数据

	if err != nil {
		return t, err
	}

	// 添加自定义请求头
	req.Header.Add("Content-Type", "application/json")
	for key, value := range headerMap {
		req.Header.Set(key, value)
	}

	// 发送请求
	res, err := client.Do(req)
	if err != nil {
		return t, err
	}
	defer res.Body.Close()

	// 检查响应状态码，非 200 状态码下直接返回错误
	if res.StatusCode != http.StatusOK {
		return t, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// 尝试将响应体解析为 JSON，并解码到 t 中
	err = json.NewDecoder(res.Body).Decode(&t)
	return t, err
}
