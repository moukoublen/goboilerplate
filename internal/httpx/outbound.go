package httpx

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// DrainAndCloseResponse can be used (most probably with defer) from the client side to ensure that the http response body is consumed til the end and closed.
func DrainAndCloseResponse(res *http.Response, errOut *error) {
	if res == nil || res.Body == nil || res.Body == http.NoBody {
		return
	}

	_, discardErr := io.Copy(io.Discard, res.Body)
	closeErr := res.Body.Close()

	if discardErr != nil || closeErr != nil {
		*errOut = errors.Join(*errOut, discardErr, closeErr)
	}
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// DoAndDecode performs the request (req) and tries to json decodes the response to output, it handles gzip and flate compression.
func DoAndDecode(c HTTPClient, req *http.Request, output any) (er error) {
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer DrainAndCloseResponse(res, &er)

	if res.StatusCode >= http.StatusBadRequest {
		return NewStatusCodeError(res.StatusCode)
	}

	reader := res.Body

	var cer error
	switch res.Header.Get("Content-Encoding") {
	// case "br": // TODO
	case "gzip":
		reader, cer = gzip.NewReader(res.Body)
	case "deflate":
		reader = flate.NewReader(res.Body)
	}
	if cer != nil {
		return cer
	}

	return json.NewDecoder(reader).Decode(output)
}

func NewStatusCodeError(statusCode int) *StatusCodeError {
	return &StatusCodeError{statusCode: statusCode}
}

type StatusCodeError struct {
	statusCode int // e.g. 200
}

func (s *StatusCodeError) Error() string {
	return fmt.Sprintf("http status code %d", s.statusCode)
}

func (s *StatusCodeError) StatusCode() int {
	return s.statusCode
}

func (s *StatusCodeError) Is(target error) bool {
	//nolint:errorlint
	if other, is := target.(*StatusCodeError); is {
		return s.statusCode == other.statusCode
	}

	return false
}

// NewHTTPClient returns a new default http client.
func NewHTTPClient(globalTimeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: globalTimeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second, //nolint:mnd,gomnd
				KeepAlive: 30 * time.Second, //nolint:mnd,gomnd
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,              //nolint:mnd,gomnd
			IdleConnTimeout:       90 * time.Second, //nolint:mnd,gomnd
			TLSHandshakeTimeout:   10 * time.Second, //nolint:mnd,gomnd
			ExpectContinueTimeout: 1 * time.Second,  //nolint:mnd,gomnd
		},
	}
}
