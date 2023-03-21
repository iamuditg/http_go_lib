package http_go_lib

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

// RequestOptions is struct that holds options for an HTTP request.
type RequestOptions struct {
	Headers      map[string]string      // custom headers to be added to the request
	Body         io.Reader              // request body (for methods like POST,PUT and PATCH)
	MaxRetries   int                    // Maximum number of retries for the request
	RetryWait    time.Duration          // Time to wait between retries
	Middlewares  []Middleware           // List of middleware to apply to the request
	Timeout      time.Duration          // Timeout for the request
	QueryParams  map[string]interface{} // Query parameters to be added to the request URL
	LogLevel     LogLevel
	LogTransport bool
}

// Middleware is a function that takes an http.RoundTripper and returns a modified http.RoundTripper.
type Middleware func(tripper http.RoundTripper) http.RoundTripper

// makeRequest is a private helper function that creates and sends an http request using the specified method, URL, and options.
func makeRequest(ctx context.Context, method, urlStr string, options *RequestOptions) (*http.Response, error) {
	// Parse the url and add query parameters
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	for key, value := range options.QueryParams {
		query.Set(key, value.(string))
	}
	u.RawQuery = query.Encode()

	// create a new HTTP request with the specified method, URL, and body
	req, err := http.NewRequest(method, u.String(), options.Body)
	if err != nil {
		return nil, err
	}
	// Add the context to the request
	req = req.WithContext(ctx)

	// set custom headers
	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	// create and HTTP client and apply middleware to its transport
	client := &http.Client{
		Timeout: options.Timeout,
	}
	for _, middleware := range options.Middlewares {
		client.Transport = middleware(client.Transport)
	}

	// wrap the transport with logging middleware if logTransport is true
	if options.LogTransport {
		client.Transport = LoggingMiddleware(client.Transport, options.LogLevel)
	}

	// perform the request with retries
	var resp *http.Response
	for i := 0; i <= options.MaxRetries; i++ {
		resp, err = client.Do(req)
		if err == nil {
			return resp, nil
		}

		select {
		case <-ctx.Done():
			return nil, errors.New("request canceled")
		case <-time.After(options.RetryWait):
			continue
		}
	}

	return nil, err
}
