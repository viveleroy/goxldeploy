package goxldeploy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	libraryVersion = "0.0.1"
	basePath       = ""
	userAgent      = "goxldeploy " + libraryVersion
	mediaType      = "application/json"
	format         = "json"
)

// Config holds configuration for the client
type Config struct {
	User     string
	Password string
	Host     string
	Port     int
	Context  string
	Scheme   string
}

// APIError defines the MailChimp API response error structure.
type APIError struct {
	Detail string
}

type ClientError struct {
	Detail string
}

//Error method for APIError to satisfy the Error type interface
func (e APIError) Error() string {
	return fmt.Sprintf("xldeploy: API Error: %s", e.Detail)
}

//Error method for APIError to satisfy the Error type interface
func (e ClientError) Error() string {
	return fmt.Sprintf("xldeploy: Client Error: %s", e.Detail)
}

// A Client manages communication with XL-Deploy
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.
	baseURL *url.URL

	// User agent
	UserAgent string

	// Client configuration
	Config *Config

	// Services
	Metadata   MetadataService
	Repository RepositoryService
}

//NewClient returns a new functional client struct
func NewClient(config *Config) *Client {
	// create the base url out of the stuff given
	var baseURL url.URL
	finalHost := config.Host + ":" + strconv.Itoa(config.Port)
	baseURL.Host = finalHost
	baseURL.Path = basePath
	baseURL.Scheme = config.Scheme

	c := &Client{client: http.DefaultClient, baseURL: &baseURL, UserAgent: userAgent, Config: config}

	c.Metadata = MetadataService{client: c}
	c.Repository = RepositoryService{client: c}

	return c
}

// New just returns a NewClient
func New(config *Config) *Client {
	return NewClient(config)
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(urlStr string, method string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(rel)
	buf := new(bytes.Buffer)

	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Config.User, c.Config.Password)
	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// Do executes request and returns response
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {

	//Execute the request
	resp, err := c.client.Do(req)

	//Handle programmetical errors
	if err != nil {
		return nil, err
	}

	// Handle API error.
	// Borrowed this code from MailChimp
	if resp.StatusCode >= 400 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		apiErr := APIError{
			Detail: buf.String(),
		}

		return nil, apiErr
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return nil, err
	}

	return resp, err
}

// Connected checks if the client is connected to XL-Deploy
func (c *Client) Connected() bool {

	req, err := c.NewRequest("deployit/server/info", "GET", nil)
	resp, err := c.client.Do(req)

	if err == nil && resp.StatusCode == 200 {
		return true
	}

	return false
}
