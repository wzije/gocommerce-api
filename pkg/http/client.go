package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ecommerce-api/pkg/config"
	"github.com/ecommerce-api/pkg/exception"
	"github.com/ecommerce-api/pkg/helper"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func AuthRequestClient(URL string, method string, bodyParams *[]byte, expectStatusCode int) (*BodyRemote, error) {
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)

	reqHeader, err := GetReqHeader()

	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+reqHeader.TokenRaw)
	request.Header.Add("X-Header-Outlet-Id", reqHeader.OutletId)
	request.Header.Add("X-Header-Settings", reqHeader.Setting)
	request.Header.SetContentType("application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.SetMethod(method)
	request.SetHost(config.Config.AppUrl)
	request.SetRequestURI(URL)

	//info request
	info := fmt.Sprintf("send request %s to %s", method, URL)
	fmt.Println(info)
	logrus.Info(info)

	if bodyParams != nil {
		request.SetBody(*bodyParams)
	}

	//prepare response
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)

	//request remotely
	err = fasthttp.Do(request, response)
	if err != nil {
		fmt.Printf("Client get failed: %s\n", err)
		return nil, err
	}

	if response.StatusCode() != expectStatusCode && (response.StatusCode() < 200 && response.StatusCode() >= 300) {
		fmt.Printf("Expected status code %d but got %d\n", fasthttp.StatusOK, response.StatusCode())
		var e BodyRemote
		err := json.Unmarshal(response.Body(), &e)
		if err != nil {
			return nil, err
		}

		code, err := helper.AnyToInt(e["code"])
		if err != nil {
			return nil, err
		}

		return nil, exception.New(*code, (e)["message"].(string), (e)["data"])
	}

	// Verify the content type
	contentType := response.Header.Peek("Content-Type")
	if bytes.Index(contentType, []byte("application/json")) != 0 {
		fmt.Printf("Expected content type application/json but got %s\n", contentType)
		return nil, err
	}

	// Do we need to decompress the responses?
	contentEncoding := response.Header.Peek("Content-Encoding")
	var responseBody []byte
	if bytes.EqualFold(contentEncoding, []byte("gzip")) {
		fmt.Println("Unzipping...")
		responseBody, _ = response.BodyGunzip()
	} else {
		responseBody = response.Body()
	}

	//encode responseBody with BodyRemote
	var data BodyRemote
	if err := json.Unmarshal(responseBody, &data); err != nil {
		return nil, err
	}

	code, err := helper.AnyToInt(data["code"])
	if err != nil {
		return nil, err
	}

	if *code < 200 && *code >= 300 {
		message := (data)["message"].(string)
		return nil, exception.New(*code, message, nil)
	}

	return &data, nil

}

func RequestClient(URL string, method string, bodyParams *[]byte, expectStatusCode int) (*BodyRemote, error) {
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)

	var err error

	request.Header.SetContentType("application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.SetMethod(method)
	request.SetHost(config.Config.AppUrl)
	request.SetRequestURI(URL)

	//info request
	info := fmt.Sprintf("send request %s to %s", method, URL)
	fmt.Println(info)
	logrus.Info(info)

	if bodyParams != nil {
		request.SetBody(*bodyParams)
	}

	//prepare response
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)

	//request remotely
	err = fasthttp.Do(request, response)
	if err != nil {
		fmt.Printf("Client get failed: %s\n", err)
		return nil, err
	}

	if response.StatusCode() != expectStatusCode && (response.StatusCode() < 200 && response.StatusCode() >= 300) {
		fmt.Printf("Expected status code %d but got %d\n", fasthttp.StatusOK, response.StatusCode())
		var e BodyRemote
		err := json.Unmarshal(response.Body(), &e)
		if err != nil {
			return nil, err
		}

		code, err := helper.AnyToInt(e["code"])
		if err != nil {
			return nil, err
		}

		return nil, exception.New(*code, (e)["message"].(string), (e)["data"])
	}

	// Verify the content type
	contentType := response.Header.Peek("Content-Type")
	if bytes.Index(contentType, []byte("application/json")) != 0 {
		fmt.Printf("Expected content type application/json but got %s\n", contentType)
		return nil, err
	}

	// Do we need to decompress the responses?
	contentEncoding := response.Header.Peek("Content-Encoding")
	var responseBody []byte
	if bytes.EqualFold(contentEncoding, []byte("gzip")) {
		fmt.Println("Unzipping...")
		responseBody, _ = response.BodyGunzip()
	} else {
		responseBody = response.Body()
	}

	//encode responseBody with BodyRemote
	var data BodyRemote
	if err := json.Unmarshal(responseBody, &data); err != nil {
		return nil, err
	}

	code, err := helper.AnyToInt(data["code"])
	if err != nil {
		return nil, err
	}

	if *code < 200 && *code >= 300 {
		message := (data)["message"].(string)
		return nil, exception.New(*code, message, nil)
	}

	return &data, nil

}
