package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ajg/form"
	errorpkg "github.com/devpayments/common/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	NOT_ASSIGNED                              = ""
	DEFAULT_TIMEOUT                           = 20 * time.Second
	AUTHORIZATION_TOKEN_PREFIX                = "Bearer"
	CONTENT_TYPE_APP_JSON         ContentType = "application/json"
	CONTENT_TYPE_FORM_URL_ENCODED ContentType = "application/x-www-form-urlencoded"
)

type ContentType string

type HttpRequest struct {
	endpoint      string            // the api endpoint
	method        string            // the http request method, get/post
	body          interface{}       // the request body, should be json serializable
	queryParams   map[string]string // the request query url params
	authUserName  string            // basic auth user name
	authPassword  string            // basic auth password
	authToken     string            // to be used in cases where token based authentication is required
	host          string            // to specify host header
	customHeaders map[string]string // to set any custom headers, if any
	contentType   ContentType       // to set content type
	timeout       time.Duration     // to set custom request timeout if needed, default is 20 secs
}

func NewHttpRequest(endpoint string, method string) *HttpRequest {
	return &HttpRequest{endpoint: endpoint, method: method, body: nil, queryParams: make(map[string]string, 0), customHeaders: make(map[string]string, 0), contentType: CONTENT_TYPE_APP_JSON}
}

func (a *HttpRequest) SetBasicAuth(username string, password string) {
	a.authPassword = password
	a.authUserName = username
}

func (a *HttpRequest) SetAuthToken(authToken string) {
	a.authToken = authToken
}

func (a *HttpRequest) SetBody(requestBody interface{}) {
	a.body = requestBody
}

func (a *HttpRequest) SetQueryParams(requestQueryParams map[string]string) {
	a.queryParams = requestQueryParams
}

func (a *HttpRequest) SetHost(host string) {
	a.host = host
}

func (a *HttpRequest) SetContentType(contentType ContentType) {
	a.contentType = contentType
}

func (a *HttpRequest) SetTimeout(timeout time.Duration) {
	a.timeout = timeout
}

func (a *HttpRequest) AddQueryParam(key, val string) {
	if a.queryParams == nil {
		a.queryParams = make(map[string]string, 0)
	}
	a.queryParams[key] = val
}

func (a *HttpRequest) AddHeader(key, val string) {
	if a.customHeaders == nil {
		a.customHeaders = make(map[string]string, 0)
	}
	a.customHeaders[key] = val
}

func MakeApiCallWithRetries(ctx context.Context, request *HttpRequest, responseAddr interface{}, retriesCount int) error {
	for i := 0; i <= retriesCount; i++ {
		err := MakeApiCall(ctx, request, responseAddr)
		v, ok := err.(errorpkg.CustomError)

		var customError *errorpkg.CustomError
		errors.As(err, &customError)

		if ok && (v.ErrorCode() == errorpkg.API_REQUEST_ERROR || v.ErrorCode() == errorpkg.API_REQUEST_STATUS_ERROR) {
			continue
		} else {
			return err
		}
	}
	return nil
}

func MakeApiCall(ctx context.Context, request *HttpRequest, responseAddr interface{}) error {
	body, err := MakeApiCallWithRawResponse(ctx, request)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &responseAddr)
	if err != nil {
		return errorpkg.NewCustomError(errorpkg.JSON_SERIALIZATION_ERROR, err.Error()).
			WithParam("response", string(body)).
			WithParam("request", fmt.Sprintf("%+v", request))
	}

	return nil
}

func MakeApiCallWithRawResponse(ctx context.Context, request *HttpRequest) ([]byte, error) {
	var requestBuffer io.Reader

	if request.body != nil {
		switch request.contentType {
		case CONTENT_TYPE_FORM_URL_ENCODED:
			values, err := form.EncodeToValues(request.body)
			if err != nil {
				err = errorpkg.NewCustomError(errorpkg.FORM_SERIALIZATION_ERROR, err.Error()).
					WithParam("request", fmt.Sprintf("%+v", request))
				return nil, err
			}
			requestBuffer = strings.NewReader(values.Encode())
		case CONTENT_TYPE_APP_JSON:
			requestJSONForm, err := json.Marshal(request.body)
			if err != nil {
				err = errorpkg.NewCustomError(errorpkg.JSON_SERIALIZATION_ERROR, err.Error()).
					WithParam("request", fmt.Sprintf("%+v", request))
				return nil, err
			}
			requestBuffer = bytes.NewBuffer(requestJSONForm)
		}
	} else {
		requestBuffer = nil
	}

	url, err := url.Parse(request.endpoint)
	if err != nil {
		err = errorpkg.NewCustomError(errorpkg.API_URL_PARSING_ERROR, err.Error()).
			WithParam("request", fmt.Sprintf("%+v", request))
		return nil, err
	}

	httpRequest, err := http.NewRequest(request.method, url.String(), requestBuffer)
	if request.queryParams != nil {
		q := httpRequest.URL.Query()
		for k, v := range request.queryParams {
			q.Add(k, v)
		}
		httpRequest.URL.RawQuery = q.Encode()
	}

	if err != nil {
		err = errorpkg.NewCustomError(errorpkg.API_REQUEST_CREATION_ERROR, err.Error()).
			WithParam("request", fmt.Sprintf("%+v", request))
		return nil, err
	}

	client := &http.Client{Timeout: DEFAULT_TIMEOUT}
	if request.timeout != time.Duration(0) {
		client.Timeout = request.timeout
	}

	httpRequest.Header.Set("Content-Type", string(request.contentType))

	if request.authToken != NOT_ASSIGNED {
		httpRequest.Header.Set("Authorization", fmt.Sprintf("%s %s", AUTHORIZATION_TOKEN_PREFIX, request.authToken))
	}

	if request.host != NOT_ASSIGNED {
		httpRequest.Host = request.host
	}

	if request.authUserName != NOT_ASSIGNED && request.authPassword != NOT_ASSIGNED {
		httpRequest.SetBasicAuth(request.authUserName, request.authPassword)
	}

	if request.customHeaders != nil {
		for k, v := range request.customHeaders {
			httpRequest.Header.Add(k, v)
		}
	}

	response, httpErr := client.Do(httpRequest)

	if httpErr != nil {
		err = errorpkg.NewCustomError(errorpkg.API_REQUEST_ERROR, err.Error()).
			WithParam("request", fmt.Sprintf("%+v", request))
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		responseMap := map[string]string{}
		err = json.Unmarshal(body, &responseMap)

		err = errorpkg.NewCustomError(errorpkg.FORM_SERIALIZATION_ERROR, err.Error()).
			WithParam("request", fmt.Sprintf("%+v", request)).
			WithParam("response", string(body)).
			WithParam("response-json", string(body))
		return nil, err
	}

	body, _ := ioutil.ReadAll(response.Body)
	return body, nil
}
