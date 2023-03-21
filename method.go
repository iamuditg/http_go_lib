package http_go_lib

import (
	"context"
	"net/http"
)

// Get sends an HTTP GET request to the specified URL with the given options.
func Get(ctx context.Context, url string, options *RequestOptions) (*http.Response, error) {
	return makeRequest(ctx, http.MethodGet, url, options)
}

// Post sends an HTTP POST request to the specified URL with the given options.
func Post(ctx context.Context, url string, options *RequestOptions) (*http.Response, error) {
	return makeRequest(ctx, http.MethodPost, url, options)
}

// Put sends an HTTP PUT request to the specified URL with the given options.
func Put(ctx context.Context, url string, options *RequestOptions) (*http.Response, error) {
	return makeRequest(ctx, http.MethodPut, url, options)
}

// Delete sends an HTTP DELETE request to the specified URL with the given options.
func Delete(ctx context.Context, url string, options *RequestOptions) (*http.Response, error) {
	return makeRequest(ctx, http.MethodDelete, url, options)
}

// Head sends an HTTP HEAD request to the specified URL with the given options.
func Head(ctx context.Context, url string, options *RequestOptions) (*http.Response, error) {
	return makeRequest(ctx, http.MethodHead, url, options)
}

// Patch sends an HTTP PATCH request with the given context, URL, and options.
func Patch(ctx context.Context, url string, options *RequestOptions) (*http.Response, error) {
	return makeRequest(ctx, http.MethodPatch, url, options)
}

// Options sends an HTTP OPTIONS request with the given context, URL, and options.
func Options(ctx context.Context, url string, options *RequestOptions) (*http.Response, error) {
	return makeRequest(ctx, http.MethodOptions, url, options)
}
