package http_go_lib

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

// LogLevel is an enumeration to represent the logging level.
type LogLevel int

const (
	LogLevelNone  LogLevel = iota // LogLevelNone disables logging.
	LogLevelBasic                 // LogLevelBasic logs basic request and response information (headers, status, etc)
	LogLevelBody                  // LogLevelBody logs request and response bodies in addition to basic information.
)

// loggingRoundTripper is a struct that wraps an http.RoundTripper and includes a LogLevel.
type loggingRoundTripper struct {
	transport http.RoundTripper // The Original transport to be wrapped for logging.
	logLevel  LogLevel          // The level of logging to be used for this round tripper
}

// RoundTrip is a method for the loggingRoundTripper struct that logs the request and response details
// based on the specified logLevel and then calls the original RoundTrip method of the wrapped transport.
func (lrt *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// Log the request based on the log level
	if lrt.logLevel >= LogLevelBasic {
		reqDump, err := httputil.DumpRequestOut(req, lrt.logLevel == LogLevelBody)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error dumping request: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "Request: \n%s\n", string(reqDump))
		}
	}

	// call the roundTrip method of the wrapped transport
	resp, err := lrt.transport.RoundTrip(req)

	// Log the response based on the log Level.
	if err == nil && lrt.logLevel >= LogLevelBasic {
		respDump, err := httputil.DumpResponse(resp, lrt.logLevel == LogLevelBody)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error dumping response: %v\n", err)
		} else {
			elapsed := time.Since(start)
			fmt.Fprintf(os.Stderr, "Response (in %s): \n%s\n", elapsed, string(respDump))
		}
	}
	return resp, err
}

// LoggingMiddleware is a function that takes an http.RoundTripper and a LogLevel and returns a new
// http.RoundTripper that wraps the original with logging functionality based on the specified LogLevel.
func LoggingMiddleware(rt http.RoundTripper, logLevel LogLevel) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return &loggingRoundTripper{transport: rt, logLevel: logLevel}
}
