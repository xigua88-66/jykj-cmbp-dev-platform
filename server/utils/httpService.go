package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// HttpService 函数调整为可以处理GET请求，其中查询参数编码到URL中
func HttpService(requestURL string, method string, queryParameters []byte, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	var req *http.Request
	var err error

	if method == "GET" && queryParameters != nil {
		// 将请求体参数转换为查询字符串
		var queryParams map[string]interface{}

		json.Unmarshal(queryParameters, &queryParams) // 确保无误地解析

		// 构建查询字符串
		queryValues := url.Values{}
		for key, value := range queryParams {
			if value == nil {
				continue
			}
			// 检查value是否为空字符串
			valueStr, ok := value.(string)
			if ok && valueStr == "" {
				continue
			}
			snakeKey := CamelToSnake(key)
			queryValues.Add(snakeKey, fmt.Sprintf("%v", value))
		}

		// 将查询字符串附加到URL
		fullURL := fmt.Sprintf("%s?%s", requestURL, queryValues.Encode())
		fmt.Println("完整请求地址：", fullURL)
		req, err = http.NewRequest(method, fullURL, nil) // GET请求没有body
	} else {
		// 对于非GET请求，正常处理请求体
		req, err = http.NewRequest(method, requestURL, bytes.NewBuffer(queryParameters))
	}

	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	fmt.Println(req.Method, req.Header, req.URL, req.Form)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return bodyBytes, nil
	}
	fmt.Println(req)
	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
