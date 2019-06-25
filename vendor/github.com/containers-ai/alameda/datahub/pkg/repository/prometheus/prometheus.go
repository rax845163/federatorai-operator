package prometheus

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/containers-ai/alameda/pkg/utils/log"
	"github.com/pkg/errors"
)

const (
	apiPrefix    = "/api/v1"
	epQuery      = "/query"
	epQueryRange = "/query_range"

	// StatusSuccess Status string literal of prometheus api request
	StatusSuccess = "success"
	// StatusError Status string literal of prometheus api request
	StatusError = "error"
)

var (
	scope              = log.RegisterScope("prometheus", "metrics repository", 0)
	defaultStepTime, _ = time.ParseDuration("30s")
)

// Entity Structure to store metrics data from Prometheus response
type Entity struct {
	Labels map[string]string
	Values []UnixTimeWithSampleValue
}

// Prometheus Prometheus api client
type Prometheus struct {
	config    Config
	client    *http.Client
	transport *http.Transport
}

// New New Prometheus api client with configuration
func New(config Config) (*Prometheus, error) {

	var (
		err error

		requestTimeout   = 30 * time.Second
		handShakeTimeout = 5 * time.Second
	)

	if err = config.Validate(); err != nil {
		return nil, errors.New("create prometheus instance failed: " + err.Error())
	}

	tr := &http.Transport{
		TLSHandshakeTimeout: handShakeTimeout,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.TLSConfig.InsecureSkipVerify,
		},
	}
	client := &http.Client{
		Timeout:   requestTimeout,
		Transport: tr,
	}

	if config.BearerTokenFile != "" {
		token, err := ioutil.ReadFile(config.BearerTokenFile)
		if err != nil {
			return nil, errors.New("create prometheus instance failed: open bearer token file for prometheus failed: " + err.Error())
		}
		config.bearerToken = string(token)
	}

	return &Prometheus{
		config:    config,
		client:    client,
		transport: tr,
	}, nil
}

// Query Query data over a range of time from prometheus
func (p *Prometheus) Query(query string, startTime, timeout *time.Time) (Response, error) {

	var (
		err error

		endpoint        = apiPrefix + epQuery
		queryParameters = url.Values{}

		u            *url.URL
		httpRequest  *http.Request
		httpResponse *http.Response

		response Response
	)

	queryParameters.Set("query", query)

	if startTime != nil {
		queryParameters.Set("time", strconv.FormatInt(int64(startTime.Unix()), 10))
	}

	if timeout != nil {
		queryParameters.Set("timeout", strconv.FormatInt(int64(timeout.Unix()), 10))
	}

	u, err = url.Parse(p.config.URL + endpoint)
	if err != nil {
		return Response{}, errors.New("prometheus query failed: url parse failed: " + err.Error())
	}
	u.RawQuery = queryParameters.Encode()

	httpRequest, err = http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return Response{}, errors.New("Query: " + err.Error())
	}
	if token := p.config.bearerToken; token != "" {
		h := http.Header{
			"Authorization": []string{fmt.Sprintf(" Bearer %s", token)},
		}
		httpRequest.Header = h
	}

	httpResponse, err = p.client.Do(httpRequest)
	if err != nil {
		return Response{}, errors.New("prometheus query failed: send http request failed" + err.Error())
	}
	err = decodeHTTPResponse(httpResponse, &response)
	if err != nil {
		return Response{}, errors.Wrap(err, "prometheus query failed")
	}

	defer p.Close()

	return response, nil
}

// QueryRange Query data over a range of time from prometheus
func (p *Prometheus) QueryRange(query string, startTime, endTime *time.Time, stepTime *time.Duration) (Response, error) {

	var (
		err error

		startTimeString string
		endTimeString   string
		stepTimeString  string

		endpoint        = apiPrefix + epQueryRange
		queryParameters = url.Values{}

		u            *url.URL
		httpRequest  *http.Request
		httpResponse *http.Response

		response Response
	)

	if startTime == nil {
		tmpTime := time.Unix(0, 0)
		startTime = &tmpTime
	}
	startTimeString = strconv.FormatInt(int64(startTime.Unix()), 10)

	if endTime == nil {
		tmpTime := time.Now()
		endTime = &tmpTime
	}
	endTimeString = strconv.FormatInt(int64(endTime.Unix()), 10)

	if stepTime == nil {
		stepTime = &defaultStepTime
	}
	stepTimeString = strconv.FormatInt(int64(stepTime.Nanoseconds()/int64(time.Second)), 10)

	queryParameters.Set("query", query)
	queryParameters.Set("start", startTimeString)
	queryParameters.Set("end", endTimeString)
	queryParameters.Set("step", stepTimeString)

	u, err = url.Parse(p.config.URL + endpoint)
	if err != nil {
		return Response{}, errors.New("prometheus query_range failed: url parse failed: " + err.Error())
	}
	u.RawQuery = queryParameters.Encode()

	httpRequest, err = http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return Response{}, errors.New("Query: " + err.Error())
	}
	if token := p.config.bearerToken; token != "" {
		h := http.Header{
			"Authorization": []string{fmt.Sprintf(" Bearer %s", token)},
		}
		httpRequest.Header = h
	}

	httpResponse, err = p.client.Do(httpRequest)
	if err != nil {
		return Response{}, errors.New("prometheus query_range failed: send http request failed" + err.Error())
	}
	err = decodeHTTPResponse(httpResponse, &response)
	if err != nil {
		return Response{}, errors.Wrap(err, "prometheus query_range failed")
	}

	defer p.Close()

	return response, nil
}

// Close Free resoure used by Prometehus
func (p *Prometheus) Close() {
	p.transport.CloseIdleConnections()
}

func decodeHTTPResponse(httpResponse *http.Response, response *Response) error {

	var (
		err                    error
		httpResponseBody       []byte
		httpResponseBodyReader io.Reader
	)

	defer httpResponse.Body.Close()

	httpResponseBody, err = ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return errors.New("decode http response failed: read http response body failed: " + err.Error())
	}

	httpResponseBodyReader = strings.NewReader(string(httpResponseBody))
	err = json.NewDecoder(httpResponseBodyReader).Decode(&response)
	if err != nil {
		return errors.New("decode http response failed: " + err.Error() + " \n received response: " + string(httpResponseBody))
	}

	return nil
}
