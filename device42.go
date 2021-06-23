package device42

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	apiVersion     = "1.0"
	defaultTimeout = 60
)

// The API client
type Api struct {
	Username   string
	Password   string
	BaseUrl    string
	TimeOut    int
	httpClient *http.Client
}

// Creates a API client that uses basic auth
func NewApiBasicAuth(username string, password string, baseUrl string, ignoreTlsErrors bool, timeOut int) (*Api, error) {
	if username == "" || password == "" {
		return nil, errors.New(ErrorEmptyCredentials)
	}

	api := &Api{
		BaseUrl:    baseUrl + "/api/" + apiVersion + "/",
		Username:   username,
		Password:   password,
		httpClient: &http.Client{},
	}

	if ignoreTlsErrors {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		api.httpClient.Transport = tr
	}

	if timeOut == 0 {
		timeOut = defaultTimeout
	}

	api.TimeOut = timeOut
	api.httpClient.Timeout = time.Duration(api.TimeOut) * time.Second

	return api, nil
}

func (api *Api) Do(method, url string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(api.Username, api.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return body, nil
}
