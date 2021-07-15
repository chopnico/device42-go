package device42

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/chopnico/output"
)

const (
	apiVersion     = "1.0"
	defaultLogging = "info"
	defaultTimeout = 60
	apiPath        = "/api/" + apiVersion
)

// the api client
type Api struct {
	options    map[string]interface{}
	httpClient *http.Client
}

// api error
type ApiResponse struct {
	Code    int           `json:"code"`
	Message []interface{} `json:"msg"`
}

// REQUIRED: sets username
func (api *Api) username(v string) *Api {
	api.option("username", v)
	return api
}

// REQUIRED: sets password
func (api *Api) password(v string) *Api {
	api.option("password", v)
	return api
}

// REQUIRED: sets host of firewall endpoint
func (api *Api) url(v string) *Api {
	api.option("url", v)
	return api
}

// sets timeout
func (api *Api) Timeout(v int) *Api {
	api.option("timeout", v)
	return api
}

// sets whether the http client should ignore ssl errors
func (api *Api) IgnoreSslErrors() *Api {
	api.options["ignore-ssl"] = true
	api.httpOptions()
	return api
}

// sets the proxy for the http client
func (api *Api) Proxy(v string) *Api {
	api.options["proxy"] = v
	api.httpOptions()
	return api
}

// sets the info logger
func (api *Api) InfoLogger(v *log.Logger) *Api {
	api.option("logger-info", v)
	return api
}

// sets the debug logger
func (api *Api) DebugLogger(v *log.Logger) *Api {
	api.option("logger-debug", v)
	return api
}

// sets the logging level
func (api *Api) LoggingLevel(v string) *Api {
	api.option("logging-level", v)
	return api
}

// debug writer
func (api *Api) WriteToDebugLog(msg string) {
	logger := api.options["logger-debug"].(*log.Logger)
	logger.Println(msg)
}

// info writer
func (api *Api) WriteToInfoLog(msg string) {
	logger := api.options["logger-info"].(*log.Logger)
	logger.Println(msg)
}

func (api *Api) option(k string, v interface{}) {
	if api.options == nil {
		api.options = make(map[string]interface{})
		api.option("logging-level", defaultLogging)
		api.option("timeout", defaultTimeout)
		api.option("ignore-ssl", false)
		api.option("proxy", "")
	}
	api.options[k] = v
}

// Build HTTP client options
func (api Api) httpOptions() {
	tr := &http.Transport{}

	if api.options["ignore-ssl"].(bool) {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if api.options["proxy"].(string) != "" {
		u, _ := url.Parse(api.options["proxy"].(string))
		tr.Proxy = http.ProxyURL(u)
	}

	api.httpClient.Transport = tr

	api.httpClient.Timeout = time.Duration(api.options["timeout"].(int)) * time.Second
}

// create the default debug logger
func defaultInfoLogger() *log.Logger {
	return log.New(os.Stderr, "[INFO] ", log.LstdFlags)
}

// create the default debug logger
func defaultDebugLogger() *log.Logger {
	return log.New(os.Stderr, "[DEBUG] ", log.LstdFlags)
}

func newApiResponse(b []byte) ApiResponse {
	r := ApiResponse{}
	json.Unmarshal(b, &r)
	return r
}

// checks to see if logging level is set to debug
func (api *Api) isLoggingDebug() bool {
	if api.options["logging-level"].(string) == "debug" {
		return true
	}
	return false
}

// checks to see if logging level is set to info
func (api *Api) isLoggingInfo() bool {
	if api.options["logging-level"].(string) == "info" {
		return true
	}
	return false
}

// Creates a API client that uses basic auth
func NewApiBasicAuth(username string, password string, host string) (*Api, error) {
	api := Api{
		httpClient: &http.Client{},
	}

	api.option("username", username)
	api.option("password", password)
	api.option("url", "https://"+host+apiPath)
	api.option("logger-info", defaultInfoLogger())
	api.option("logger-debug", defaultDebugLogger())

	return &api, nil
}

// The main do function for an api request
// lots of debug logging code in here, but it's also the main method
// for the client.
// REVIEW : refactor
// NOTES: it's pretty ugly
func (api *Api) Do(method, path string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, api.options["url"].(string)+path, body)
	if api.isLoggingDebug() {
		api.WriteToDebugLog("request url : " + req.URL.Host)
		api.WriteToDebugLog("request method : " + req.Method)
		api.WriteToDebugLog("request headers : " + output.FormatItemAsPrettyJson(req.Header))
	}
	if err != nil {
		if api.isLoggingDebug() {
			api.WriteToDebugLog(err.Error())
			return nil, errors.New("debugging")
		}
		return nil, err
	}

	req.SetBasicAuth(api.options["username"].(string), api.options["password"].(string))
	switch method {
	case "POST":
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Accept`", "application/json")
	case "GET":
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Accept`", "application/json")
	}

	resp, err := api.httpClient.Do(req)
	if api.isLoggingDebug() {
		api.WriteToDebugLog("debug http client do")
		if err != nil {
			b, _ := ioutil.ReadAll(resp.Body)
			e := newApiResponse(b)
			api.WriteToDebugLog(fmt.Sprintf("response message : %s", e.Message))
			api.WriteToDebugLog("response status code : " + fmt.Sprintf("%d", e.Code))
		}
		api.WriteToDebugLog("request headers : " + output.FormatItemAsPrettyJson(req.Header))
		api.WriteToDebugLog("request uri : " + req.URL.RequestURI())
		api.WriteToDebugLog("request method : " + req.Method)
	}
	if err != nil {
		if api.isLoggingDebug() {
			api.WriteToDebugLog(err.Error())
			return nil, errors.New("debugging")
		}
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if api.isLoggingDebug() {
			api.WriteToDebugLog(err.Error())
			return nil, errors.New("debugging")
		}
		return nil, err
	}

	defer resp.Body.Close()

	if api.isLoggingDebug() {
		e := newApiResponse(b)

		api.WriteToDebugLog("debug http client do")
		api.WriteToDebugLog(fmt.Sprintf("response message : %s", e.Message))
		api.WriteToDebugLog("response status code : " + fmt.Sprintf("%d", e.Code))
		api.WriteToDebugLog("request headers : " + output.FormatItemAsPrettyJson(req.Header))
		api.WriteToDebugLog("request uri : " + req.URL.RequestURI())
		api.WriteToDebugLog("request method : " + req.Method)
		api.WriteToDebugLog("response headers : " + output.FormatItemAsPrettyJson(resp.Header))
	}

	switch resp.StatusCode {
	case 200:
		return b, nil
	case 400:
		return nil, errors.New("bad request... stop it")
	case 401:
		return nil, errors.New("unauthorized... don't think so")
	case 403:
		return nil, errors.New("forbidden... get out of here")
	case 404:
		return nil, errors.New("resource not found... it's gone")
	case 405:
		return nil, errors.New("method not allowed... what're you trying pull?")
	case 410:
		return nil, errors.New("gone... was it even real?")
	case 500:
		return nil, errors.New("internal server error... our stuff is broke")
	case 503:
		return nil, errors.New("service unavaliable... not sure what's going on")
	default:
		e := newApiResponse(b)
		if api.isLoggingDebug() {
			api.WriteToDebugLog(fmt.Sprintf("%s", e.Message))
			api.WriteToDebugLog(fmt.Sprintf("%s", e.Message))
			return nil, errors.New("debugging")
		}
		return nil, errors.New(fmt.Sprintf("%s", e.Message))
	}
}
