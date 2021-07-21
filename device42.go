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
	"github.com/google/uuid"
)

const (
	apiVersion     = "1.0"
	defaultLogging = "info"
	defaultTimeout = 60
	apiPath        = "/api/" + apiVersion
)

// API type
type API struct {
	options    map[string]interface{}
	httpClient *http.Client
}

// APIResponse type
type APIResponse struct {
	Code    int         `json:"code"`
	Message interface{} `json:"msg"`
}

// APIContextKey is a helper string key for contexts
type APIContextKey string

func (api *API) username(v string) *API {
	api.option("username", v)
	return api
}

func (api *API) password(v string) *API {
	api.option("password", v)
	return api
}

func (api *API) url(v string) *API {
	api.option("url", v)
	return api
}

// Timeout will set client timeout
func (api *API) Timeout(v int) *API {
	api.option("timeout", v)
	return api
}

// IgnoreSSLErrors will tell the HTTP client to ignore SSL errors
func (api *API) IgnoreSSLErrors() *API {
	api.options["ignore-ssl"] = true
	api.httpOptions()
	return api
}

// Proxy will tell the HTTP client which proxy address should be used
func (api *API) Proxy(v string) *API {
	api.options["proxy"] = v
	api.httpOptions()
	return api
}

// InfoLogger sets a custom InfoLogger
func (api *API) InfoLogger(v *log.Logger) *API {
	api.option("logger-info", v)
	return api
}

// DebugLogger sets a custom DebugLogger
func (api *API) DebugLogger(v *log.Logger) *API {
	api.option("logger-debug", v)
	return api
}

// LoggingLevel sets the log level (info, debug)
func (api *API) LoggingLevel(v string) *API {
	api.option("logging-level", v)
	return api
}

// WriteToDebugLog will write entries to the debug logger
func (api *API) WriteToDebugLog(msg string) {
	logger := api.options["logger-debug"].(*log.Logger)
	logger.Println(msg)
}

// WriteToInfoLog will write to the info logger
func (api *API) WriteToInfoLog(msg string) {
	logger := api.options["logger-info"].(*log.Logger)
	logger.Println(msg)
}

// option sets (and maybe create) options
func (api *API) option(k string, v interface{}) {
	if api.options == nil {
		api.options = make(map[string]interface{})
		api.option("logging-level", defaultLogging)
		api.option("timeout", defaultTimeout)
		api.option("ignore-ssl", false)
		api.option("proxy", "")
	}
	api.options[k] = v
}

// httpOptions will build the HTTP client
func (api API) httpOptions() {
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

// defaultInfoLogger creates an info logger
func defaultInfoLogger() *log.Logger {
	return log.New(os.Stderr, "[INFO] ", log.LstdFlags)
}

// defaultDebugLogger creates a debug logger
func defaultDebugLogger() *log.Logger {
	return log.New(os.Stderr, "[DEBUG] ", log.LstdFlags)
}

// newAPIResponse creates a new api response
func newAPIResponse(b []byte) APIResponse {
	r := APIResponse{}
	json.Unmarshal(b, &r)
	return r
}

// IsLoggingDebug checks if logging level is set to debug
func (api *API) IsLoggingDebug() bool {
	if api.options["logging-level"].(string) == "debug" {
		return true
	}
	return false
}

// IsLoggingInfo checks if logging level is set to info
func (api *API) IsLoggingInfo() bool {
	if api.options["logging-level"].(string) == "info" {
		return true
	}
	return false
}

// NewAPIBasicAuth creates a new api client using basic authentication
func NewAPIBasicAuth(username string, password string, host string) (*API, error) {
	api := API{
		httpClient: &http.Client{},
	}

	api.option("username", username)
	api.option("password", password)
	api.option("url", "https://"+host+apiPath)
	api.option("logger-info", defaultInfoLogger())
	api.option("logger-debug", defaultDebugLogger())

	return &api, nil
}

// Do is a wrapper function for the httpClient Do function
// REVIEW : refactor
// NOTES: it's pretty ugly
func (api *API) Do(method, path string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, api.options["url"].(string)+path, body)
	if api.IsLoggingDebug() {
		api.WriteToDebugLog("request url : " + req.URL.Host)
		api.WriteToDebugLog("request method : " + req.Method)
		api.WriteToDebugLog("request headers : " + output.FormatItemAsPrettyJson(req.Header))
	}
	if err != nil {
		if api.IsLoggingDebug() {
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
	req.Header.Add("Client-Transaction-ID", uuid.New().String())

	resp, err := api.httpClient.Do(req)
	if api.IsLoggingDebug() {
		api.WriteToDebugLog("debug http client do")
		if err != nil {
			b, _ := ioutil.ReadAll(resp.Body)
			e := newAPIResponse(b)
			api.WriteToDebugLog(fmt.Sprintf("response message : %s", e.Message))
			api.WriteToDebugLog("response status code : " + fmt.Sprintf("%d", e.Code))
		}
		api.WriteToDebugLog("request headers : " + output.FormatItemAsPrettyJson(req.Header))
		api.WriteToDebugLog("request uri : " + req.URL.RequestURI())
		api.WriteToDebugLog("request method : " + req.Method)
	}
	if err != nil {
		if api.IsLoggingDebug() {
			api.WriteToDebugLog(err.Error())
			return nil, errors.New("debugging")
		}
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if api.IsLoggingDebug() {
			api.WriteToDebugLog(err.Error())
			return nil, errors.New("debugging")
		}
		return nil, err
	}

	defer resp.Body.Close()

	if api.IsLoggingDebug() {
		e := newAPIResponse(b)

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
	case 503:
		return nil, errors.New("service unavaliable... not sure what's going on")
	default:
		e := newAPIResponse(b)
		if api.IsLoggingDebug() {
			api.WriteToDebugLog(fmt.Sprintf("%s", e.Message))
			api.WriteToDebugLog(fmt.Sprintf("%s", e.Message))
		}
		return nil, fmt.Errorf("%v", e.Message)
	}
}
