package httpclient

//****************************************************/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : utils/httpclient.go
// Author      : ningzhong.zeng
// Revision    : 2016-01-8 20:55:18
// Description :
//****************************************************/

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"
)

// Supported http methods
const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
)

// Request headers
const (
	HeaderAccept            = "Accept"
	HeaderAcceptCharset     = "Accept-Charset"
	HeaderAcceptEncoding    = "Accept-Encoding"
	HeaderAcceptLanguage    = "Accept-Language"
	HeaderAcceptDatetime    = "Accept-Datetime"
	HeaderAuthorization     = "Authorization"
	HeaderCacheControl      = "Cache-Control"
	HeaderConnection        = "Connection"
	HeaderCookie            = "Cookie"
	HeaderContentType       = "Content-Type"
	HeaderDate              = "Date"
	HeaderIfMatch           = "If-Match"
	HeaderIfModifiedSince   = "If-Modified-Since"
	HeaderIfNoneMatch       = "If-None-Match"
	HeaderIfRange           = "If-Range"
	HeaderIfUnmodifiedSince = "If-Unmodified-Since"
	HeaderMaxForwards       = "Max-Forwards"
	HeaderOrigin            = "Origin"
	HeaderPragma            = "Pragma"
	HeaderRange             = "Range"
	HeaderReferer           = "Referer"
	HeaderUserAgent         = "User-Agent"
	HeaderUpgrade           = "Upgrade"
	HeaderWarning           = "Warning"
)

// default values
var (
	defaultReadTimeout         = 30
	defaultTLSHandshaleTimeout = 30
	defaultKeepAliveTimeout    = 65
	defaultVerify              = true
	defaultUserAgentHeader     = fmt.Sprintf("%s-%s-%s", runtime.Version(), runtime.GOARCH, runtime.GOOS)
	defaultAcceptHeader        = "application/json"
	defaultContentType         = "application/json"
)

// request is structure that holds info about ongoing request.
type request struct {
	method string
	path   string
	query  map[string]string
	body   interface{}
}

// Client holds information about endpoint and communication configuration.
type Client struct {
	Endpoint            string
	Headers             map[string]string
	ReadTimeout         time.Duration
	TLSHandshakeTimeout time.Duration
	KeepAlive           time.Duration
	Verify              bool
	http                *http.Client
	request             *request
}

// Options is used to provide values of initialization of client.
type Options struct {
	Endpoint            string
	Headers             map[string]string
	ReadTimeout         time.Duration
	TLSHandshakeTimeout time.Duration
	KeepAlive           time.Duration
	Verify              bool
	UserAgent           string
}

// New creates and returns instance of Client with all default values.
func New() (client *Client) {
	readTimeout := time.Duration(defaultReadTimeout) * time.Second
	keepAlive := time.Duration(defaultKeepAliveTimeout) * time.Second
	tlsHandshaketimeout := time.Duration(defaultTLSHandshaleTimeout) * time.Second
	return &Client{
		Endpoint:            "",
		Headers:             make(map[string]string),
		ReadTimeout:         readTimeout,
		TLSHandshakeTimeout: tlsHandshaketimeout,
		KeepAlive:           keepAlive,
		Verify:              defaultVerify,
	}
}

// NewWithOptions creates instance of Client by using provided option values.
func NewWithOptions(options Options) (client *Client) {
	var (
		headers             map[string]string
		readTimeout         time.Duration
		tlsHandshateTimeout time.Duration
		keepAlive           time.Duration
	)

	if options.Headers != nil {
		headers = options.Headers
	} else {
		headers = make(map[string]string)
	}

	if options.ReadTimeout != 0 {
		readTimeout = time.Duration(options.ReadTimeout) * time.Second
	} else {
		readTimeout = time.Duration(defaultReadTimeout) * time.Second
	}

	if options.TLSHandshakeTimeout != 0 {
		tlsHandshateTimeout = time.Duration(options.TLSHandshakeTimeout) * time.Second
	} else {
		tlsHandshateTimeout = time.Duration(defaultTLSHandshaleTimeout) * time.Second
	}

	if options.KeepAlive != 0 {
		keepAlive = time.Duration(options.KeepAlive) * time.Second
	} else {
		keepAlive = time.Duration(defaultKeepAliveTimeout) * time.Second
	}

	cl := &Client{
		Endpoint:            options.Endpoint,
		Headers:             headers,
		ReadTimeout:         readTimeout,
		TLSHandshakeTimeout: tlsHandshateTimeout,
		KeepAlive:           keepAlive,
		Verify:              options.Verify,
	}
	if options.UserAgent != "" {
		cl.UserAgent(options.UserAgent)
	}
	return cl
}

// Config functions

// Endpoint sets root endpoint for client.
func (cl *Client) WithEndpoint(endpoint string) *Client {
	cl.Endpoint = endpoint
	return cl
}

// Timeout sets read timeout for requests made by this client.
func (cl *Client) Timeout(timeout int) *Client {
	cl.ReadTimeout = time.Duration(timeout) * time.Second
	return cl
}

// TLSTimeout sets TLS handshake timeout for requests made by this client.
func (cl *Client) TLSTimeout(timeout int) *Client {
	cl.TLSHandshakeTimeout = time.Duration(timeout) * time.Second
	return cl
}

// KeepAliveTimeout sets timeout for connection keep-alive.
func (cl *Client) KeepAliveTimeout(timeout int) *Client {
	cl.KeepAlive = time.Duration(timeout) * time.Second
	return cl
}

// Verify sets bool flag indicating if TLS certificates should be verified.
func (cl *Client) WithVerify(verify bool) *Client {
	cl.Verify = verify
	return cl
}

// Header sets header value for all requests set using this client.
func (cl *Client) Header(key, value string) *Client {
	cl.Headers[key] = value
	return cl
}

// Query sets query key - value pair to be used in next request.
func (cl *Client) Query(key, value string) *Client {
	cl.updateRequest("", "", map[string]string{key: value}, nil)
	return cl
}

// Body sets body that will be send in next request by this client.
func (cl *Client) Body(body interface{}) *Client {
	cl.updateRequest("", "", nil, body)
	return cl
}

// utility header setting functions

// UserAgent sets user agent string for next request made with this client.
func (cl *Client) UserAgent(useragent string) *Client {
	cl.Headers[HeaderUserAgent] = useragent
	return cl
}

// Accept sets Accept header for next request made with this client.
func (cl *Client) Accept(accept string) *Client {
	cl.Headers[HeaderAccept] = accept
	return cl
}

// ContentType sets Content-Type header value for next request made with
// this client.
func (cl *Client) ContentType(contentType string) *Client {
	cl.Headers[HeaderContentType] = contentType
	return cl
}

// request definition functions

// Get creates new request inside this client and markes it as GET request.
// Provided URL should be relative and will be appended to clients Endpoint.
func (cl *Client) Get(url string) *Client {
	cl.updateRequest(GET, url, nil, nil)
	// cl.request = &request{GET, url, nil}
	return cl
}

// Post creates new request inside this client and markes it as POST request.
// Provided URL should be relative and will be appended to clients Endpoint.
func (cl *Client) Post(url string) *Client {
	cl.updateRequest(POST, url, nil, nil)
	// cl.request = &request{POST, url, nil}
	return cl
}

// Put creates new request inside this client and markes it as PUT request.
// Provided URL should be relative and will be appended to clients Endpoint.
func (cl *Client) Put(url string) *Client {
	cl.updateRequest(PUT, url, nil, nil)
	// cl.request = &request{PUT, url, nil}
	return cl
}

// Delete creates new request inside this client and markes it as DELETE request.
// Provided URL should be relative and will be appended to clients Endpoint.
func (cl *Client) Delete(url string) *Client {
	cl.updateRequest(DELETE, url, nil, nil)
	// cl.request = &request{DELETE, url, nil}
	return cl
}

// PATCH creates new request inside this client and markes it as PATCH request.
// Provided URL should be relative and will be appended to clients Endpoint.
func (cl *Client) Patch(url string) *Client {
	cl.updateRequest(PATCH, url, nil, nil)
	// cl.request = &request{PATCH, url, nil}
	return cl
}

// Head creates new request inside this client and markes it as HEAD request.
// Provided URL should be relative and will be appended to clients Endpoint.
func (cl *Client) Head(url string) *Client {
	cl.updateRequest(HEAD, url, nil, nil)
	// cl.request = &request{HEAD, url, nil}
	return cl
}

// Options creates new request inside this client and markes it as OPTIONS request.
// Provided URL should be relative and will be appended to clients Endpoint.
func (cl *Client) Options(url string) *Client {
	cl.updateRequest(OPTIONS, url, nil, nil)
	// cl.request = &request{OPTIONS, url, nil}
	return cl
}

// utility functions

// updateRequest creates new request if there is no existing one for this client
// and initializes it with provided data. If request already exists it will
// be updated if provided data are not zero values.
func (cl *Client) updateRequest(method, path string, query map[string]string, body interface{}) {
	if cl.request == nil {
		cl.request = &request{
			method: method,
			path:   path,
			query:  query,
			body:   body,
		}
	} else {
		if method != "" {
			cl.request.method = method
		}
		if path != "" {
			cl.request.path = path
		}
		if query != nil {
			if cl.request.query == nil {
				cl.request.query = make(map[string]string)
			}
			for k, v := range query {
				cl.request.query[k] = v
			}
		}
		if body != nil {
			cl.request.body = body
		}
	}
}

// createUrl creates instance of URL for request in provided client.
func createUrl(cl *Client) (parsedUrl *url.URL, err error) {
	u, err := url.Parse(cl.Endpoint + cl.request.path)
	if err != nil {
		return nil, err
	}

	if cl.request.query != nil {
		q := u.Query()
		for key, val := range cl.request.query {
			q.Set(key, val)
		}
		u.RawQuery = q.Encode()
	}
	fmt.Println("Generated URL: ", u.String())
	return u, nil
}

// createRequest creates instance of http.Request containing all config
// values and URL from current client and provided body.
func createRequest(cl *Client, data io.Reader) (r *http.Request, err error) {
	url, err := createUrl(cl)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	if data != nil {
		if req, err = http.NewRequest(cl.request.method, url.String(), data); err != nil {
			return nil, err
		}
	} else {
		if req, err = http.NewRequest(cl.request.method, url.String(), nil); err != nil {
			return nil, err
		}
	}

	for key, val := range cl.Headers {
		req.Header.Set(key, val)
	}

	if req.Header.Get(HeaderUserAgent) == "" {
		req.Header.Set(HeaderUserAgent, defaultUserAgentHeader)
	}
	if req.Header.Get(HeaderContentType) == "" {
		req.Header.Set(HeaderContentType, defaultContentType)
	}
	if req.Header.Get(HeaderAccept) == "" {
		req.Header.Set(HeaderAccept, defaultAcceptHeader)
	}

	return req, nil
}

// createHttpClient creates instance of http.Client with all config data
// in provided client.
func createHttpClient(cl *Client) *http.Client {
	return &http.Client{
		Timeout: cl.ReadTimeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   cl.ReadTimeout,
				KeepAlive: cl.KeepAlive,
			}).Dial,
			TLSHandshakeTimeout: cl.TLSHandshakeTimeout,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: !cl.Verify},
		},
	}
}

// handleError checks HTTP response status code and returns error if response is
// not is range from 200 to 300.
func handleError(resp *http.Response) (err error) {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		message := fmt.Sprintf(
			"Http request resulted in error with status code %d and message %s",
			resp.StatusCode,
			resp.Status,
		)
		return errors.New(message)
	}
	return nil
}

// ending functions

// end issues request to remote host and returns response.
// end should not be called if cl.request is nil. At the end, this function
// will set request to nil to reset all values for new request.
func (cl *Client) end() (response *http.Response, err error) {
	// check if method function is called
	if cl.request == nil {
		return nil, errors.New("Method function not called (Get, Post, Put...).")
	}

	var bodyReader io.Reader
	if cl.request.body != nil {
		body, err := json.Marshal(cl.request.body)
		if err != nil {
			return nil, err
		}
		bodyReader = strings.NewReader(string(body))
	} else {
		bodyReader = nil
	}

	req, err := createRequest(cl, bodyReader)
	if err != nil {
		return nil, err
	}

	httpClient := createHttpClient(cl)

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// reset, if same client instance is reuesd
	cl.request = nil

	return resp, nil
}

// AsString returns response of server as string.
func (cl *Client) AsString() (b string, err error) {
	resp, err := cl.Raw()
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// Json returns response from server decoded from JSON into generic interface.
func (cl *Client) Json() (j interface{}, err error) {
	resp, err := cl.end()
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	var response interface{}
	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&response); err != nil {
		if errt, ok := err.(*json.SyntaxError); ok {
			return nil, fmt.Errorf("JSON Decoder: %s, offset: %d", errt, errt.Offset)
		}
		return nil, err
	}
	return response, nil
}

// JsonDecode decodes JSON response from server into provided object.
func (cl *Client) JsonDecode(response interface{}) (err error) {
	resp, err := cl.end()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(response); err != nil {
		if errt, ok := err.(*json.SyntaxError); ok {
			return fmt.Errorf("JSON Decoder: %s, offset %d", errt, errt.Offset)
		}
		return err
	}
	return nil
}

// Raw returns response from server as raw byte slice.
func (cl *Client) Raw() (response []byte, err error) {
	resp, err := cl.end()
	defer resp.Body.Close()

	if err != nil {
		return []byte{}, err
	}

	buff := new(bytes.Buffer)
	buff.ReadFrom(resp.Body)
	return buff.Bytes(), nil
}
