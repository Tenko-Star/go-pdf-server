package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Get(requestUrl string, params map[string]string, headers map[string]string) ([]byte, error) {
	// 处理url
	buf := strings.Builder{}
	requestUrl = strings.Trim(url.PathEscape(requestUrl), "?&")
	// 写入首部
	buf.WriteString(requestUrl)
	if strings.Index(requestUrl, "?") != -1 {
		// 如果存在参数则添加&符号
		buf.WriteString("&")
	} else {
		// 不存在则添加?
		buf.WriteString("?")
	}

	// 写入参数
	for key, value := range params {
		buf.WriteString(fmt.Sprintf("%s=%s&", url.QueryEscape(key), url.QueryEscape(value)))
	}

	// 创建请求
	req, err := http.NewRequest("GET", buf.String(), nil)
	if err != nil {
		return nil, err
	}

	// 处理请求头
	for header, value := range headers {
		req.Header.Add(header, value)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	ret, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return ret, errors.New(response.Status)
	}

	return ret, nil
}

func PostWithJson(requestUrl string, params map[string]interface{}, headers map[string]string) ([]byte, error) {
	var (
		err          error
		client       = http.Client{}
		request      *http.Request
		response     *http.Response
		responseBody []byte
	)

	// 处理json
	jsonBuffer, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	request, err = http.NewRequest(http.MethodPost, requestUrl, strings.NewReader(string(jsonBuffer)))
	if err != nil {
		return nil, err
	}

	// 处理请求头
	request.Header.Add("Content-Type", "application/json")
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	response, err = client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return responseBody, errors.New("fail to callback")
	}

	return responseBody, nil
}
