package http

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func UploadFile(path string, requestUrl string, fieldName string) ([]byte, error) {
	var (
		err          error
		file         *os.File
		client       = http.Client{}
		request      *http.Request
		response     *http.Response
		responseBody []byte
		bodyBuffer   = &bytes.Buffer{}
		bodyWriter   = multipart.NewWriter(bodyBuffer)
	)

	file, err = os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileWriter, err := bodyWriter.CreateFormFile(fieldName, path)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return nil, err
	}
	_ = bodyWriter.Close()

	request, err = http.NewRequest(http.MethodPost, requestUrl, bodyBuffer)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	response, err = client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Println(response.Status)
		return nil, errors.New(response.Status)
	}

	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func Download(requestUrl string, destination io.Writer) error {
	var (
		err      error
		response *http.Response
	)

	response, err = http.Get(requestUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(destination, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func DownloadToTemp(requestUrl string, tempDir string) (string, error) {
	var (
		fileName string
		ext      string
		fileInfo os.FileInfo
		temp     *os.File
		m        hash.Hash

		err error
	)

	fileInfo, err = os.Stat(tempDir)
	if err != nil {
		// 文件夹不存在时
		err = os.MkdirAll(tempDir, os.ModePerm)
		if err != nil {
			// 创建文件夹失败时
			return "", err
		}
	} else if !fileInfo.IsDir() {
		return "", errors.New("the path 'tempDir' is not a dir")
	}

	strArr := strings.Split(requestUrl, ".")
	ext = strArr[len(strArr)-1]
	m = md5.New()
	m.Write([]byte(requestUrl))
	fileName = fmt.Sprintf("%s/%s.%s", tempDir, hex.EncodeToString(m.Sum(nil)), ext)

	// 检查缓存
	fileInfo, err = os.Stat(fileName)
	if err == nil {
		// todo ref+1
		return fileName, nil
	}

	temp, err = os.Create(fileName)
	if err != nil {
		return "", err
	}

	err = Download(requestUrl, temp)
	if err != nil {
		_ = os.Remove(fileName)

		return "", err
	}

	return fileName, nil
}
