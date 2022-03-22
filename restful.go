package goutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type contentType string

const (
	ContentTypeJSON              contentType = "application/json"
	ContentTypeFormURLEncoded    contentType = "application/x-www-form-urlencoded"
	ContentTypeTextHTML          contentType = "text/html"
	ContentTypeMultipartFormData contentType = "multipart/form-data"
)

type RequestPayload struct {
	Path        string
	QueryParam  map[string]string
	Method      string
	ContentType contentType
	Payload     interface{}
	PayloadFile interface{}
	Response    interface{}
	Headers     map[string]string
}

type BasicAuthPayload struct {
	Path     string
	Method   string
	Payload  BasicAuth
	Response interface{}
	Headers  map[string]string
}

type BasicAuth struct {
	Username string
	Password string
}

type RESTful struct {
	retry   int
	baseURL string
}

func NewRESTful(baseURL string, retry int) *RESTful {
	return &RESTful{
		retry:   retry,
		baseURL: baseURL,
	}
}

func setQueryParam(baseURL, pathURL string, queryParam map[string]string) (urlRequest string, err error) {

	urlRequest = baseURL + pathURL
	urls, err := url.Parse(urlRequest)
	if err != nil {
		return urlRequest, err
	}

	q := urls.Query()
	for key, val := range queryParam {
		q.Add(key, val)
	}
	urls.RawQuery = q.Encode()

	return urls.String(), err
}

func (r *RESTful) Request(req RequestPayload) (statusCode int, err error) {

	for i := 0; i < r.retry; i++ {

		var client = &http.Client{}
		var request *http.Request
		urlRequest, err := setQueryParam(r.baseURL, req.Path, req.QueryParam)
		if err != nil {
			return statusCode, err
		}

		// setup request
		if req.ContentType == ContentTypeJSON {

			payload := &strings.Reader{}

			if req.Payload != nil {

				payloadJSON, err := json.Marshal(req.Payload)
				if err != nil {
					return statusCode, err
				}

				payload = strings.NewReader(string(payloadJSON))
			}

			request, err = http.NewRequest(req.Method, urlRequest, payload)
			if err != nil {
				return statusCode, err
			}

			request.Header.Add("content-type", string(ContentTypeJSON))

		} else if req.ContentType == ContentTypeFormURLEncoded {

			if req.Payload == nil {
				return statusCode, errors.New("payload is nil")
			}

			if reflect.TypeOf(req.Payload).String() == "url.Values" {

				values := req.Payload.(url.Values)
				request, err = http.NewRequest(req.Method, urlRequest, strings.NewReader(values.Encode()))
				if err != nil {
					return statusCode, err
				}

				request.Header.Add("Content-Type", string(ContentTypeFormURLEncoded))

			} else {
				return statusCode, errors.New("payload isn't url.Values")
			}

		} else if req.ContentType == ContentTypeMultipartFormData {

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			if req.Payload == nil && req.PayloadFile == nil {
				return statusCode, errors.New("payload is nil")
			}

			if req.Payload != nil && reflect.TypeOf(req.Payload).String() == "url.Values" {
				values := req.Payload.(url.Values)
				for key, value := range values {
					for _, i := range value {
						writer.WriteField(key, i)
					}
				}
			}

			if req.PayloadFile != nil && reflect.TypeOf(req.PayloadFile).String() == "url.Values" {
				values := req.PayloadFile.(url.Values)
				for key, value := range values {
					for _, i := range value {
						file, err := os.Open(i)
						if err != nil {
							return statusCode, err
						}
						defer file.Close()

						part, err := writer.CreateFormFile(key, filepath.Base(i))
						if err != nil {
							return statusCode, err
						}

						io.Copy(part, file)
					}
				}
			}

			writer.Close()

			request, err = http.NewRequest(req.Method, urlRequest, body)
			if err != nil {
				return statusCode, err
			}

			request.Header.Add("Content-Type", writer.FormDataContentType())

		} else {
			return statusCode, errors.New("you must to choice the content type")
		}

		// action request
		setHeader(request, req.Headers)
		response, err := client.Do(request)
		if err != nil {
			return statusCode, err
		}

		// getting content type
		// we don't need handle this error because isn't important
		contentType, _, _ := mime.ParseMediaType(response.Header.Get("content-type"))

		// if statuscode isn't ok, we'll retry the request
		// but, if the retry still failed we'll return the error
		if ((i + 1) == r.retry) && response.StatusCode >= http.StatusBadRequest {
			errorMessage := fmt.Sprintf("Status code: %d / %s", response.StatusCode, http.StatusText(response.StatusCode))
			if contentType == string(ContentTypeJSON) {
				resBytes, _ := ioutil.ReadAll(response.Body)
				errorMessage = fmt.Sprintf("%s, Response: %s", errorMessage, string(resBytes))
			}

			return response.StatusCode, errors.New(errorMessage)
		} else if response.StatusCode >= http.StatusBadRequest {
			continue
		}

		// setup response
		if contentType == string(ContentTypeJSON) {

			if req.Response != nil {

				resBytes, err := ioutil.ReadAll(response.Body)
				if err != nil {
					return response.StatusCode, err
				}

				if err = json.Unmarshal(resBytes, req.Response); err != nil {
					return response.StatusCode, err
				}
			}
		}

		return response.StatusCode, err
	}

	return statusCode, errors.New("retry need to setup greater than 0")
}

func (r *RESTful) RequestBasicAuth(req BasicAuthPayload) (statusCode int, err error) {

	for i := 0; i < r.retry; i++ {

		var client = &http.Client{}
		var request *http.Request
		var urlRequest = r.baseURL + req.Path

		if req.Payload.Username == "" {
			return statusCode, errors.New("username must be filled")
		}

		if req.Payload.Password == "" {
			return statusCode, errors.New("password must be filled")
		}

		request, err = http.NewRequest(req.Method, urlRequest, nil)
		if err != nil {
			return statusCode, err
		}

		setHeader(request, req.Headers)
		request.SetBasicAuth(req.Payload.Username, req.Payload.Password)

		// action request
		response, err := client.Do(request)
		if err != nil {
			return statusCode, err
		}

		// if statuscode isn't ok, we'll retry the request
		// but, if the retry still failed we'll return the error
		if (i + 1) == r.retry {
			return response.StatusCode, errors.New(http.StatusText(response.StatusCode))
		} else if response.StatusCode >= http.StatusBadRequest {
			continue
		}

		// getting content type
		contentType, _, err := mime.ParseMediaType(response.Header.Get("content-type"))
		if err != nil {
			return response.StatusCode, err
		}

		// setup response
		if contentType == string(ContentTypeJSON) {

			if req.Response != nil {

				resBytes, err := ioutil.ReadAll(response.Body)
				if err != nil {
					return response.StatusCode, err
				}

				if err = json.Unmarshal(resBytes, req.Response); err != nil {
					return response.StatusCode, err
				}
			}
		}

		return response.StatusCode, err
	}

	return statusCode, errors.New("retry need to setup greater than 0")
}

func setHeader(request *http.Request, headers map[string]string) {
	for h, i := range headers {
		request.Header.Add(h, i)
	}
}
