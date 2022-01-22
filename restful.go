package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type contentType string

const (
	ContentTypeJSON           contentType = "application/json"
	ContentTypeFormURLEncoded contentType = "application/x-www-form-urlencoded"
	ContentTypeTextHTML       contentType = "text/html"
)

type RequestPayload struct {
	Path        string
	Method      string
	ContentType contentType
	Payload     interface{}
	Response    interface{}
}

type restful struct {
	retry   int
	baseURL string
}

func NewRestful(baseURL string, retry int) *restful {
	return &restful{
		retry:   retry,
		baseURL: baseURL,
	}
}

func (r *restful) Request(req RequestPayload) (statusCode int, err error) {

	for i := 0; i < r.retry; i++ {

		var client *http.Client
		var request *http.Request
		var urlRequest = r.baseURL + req.Path

		// setup request
		if req.ContentType == ContentTypeJSON {

			var payload *strings.Reader

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

			request.Header.Add("content-type", string(req.ContentType))

		} else if req.ContentType == ContentTypeFormURLEncoded {

			if req.Payload != nil {

				if reflect.TypeOf(req.Payload).String() == "url.Values" {

					values := req.Payload.(url.Values)
					request, err = http.NewRequest(http.MethodPost, urlRequest, strings.NewReader(values.Encode()))
					if err != nil {
						return statusCode, err
					}

					request.Header.Add("Content-Type", string(ContentTypeFormURLEncoded))

				} else {
					return statusCode, errors.New("payload isn't url.Values")
				}
			} else {
				return statusCode, errors.New("payload is nil")
			}
		} else {
			return statusCode, errors.New("you must to choice the content type")
		}

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
