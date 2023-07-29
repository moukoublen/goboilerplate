package httpx

import (
	"compress/flate"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

func DrainAndCloseResponse(res *http.Response) {
	if res == nil || res.Body == nil {
		return
	}

	_, _ = io.Copy(io.Discard, res.Body)
	_ = res.Body.Close()
}

type InnerClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	InnerClient
}

// DoAndDecode performs the request (req) and tries to json decodes the response to output. It also handles gzip and flate compression.
func (c *Client) DoAndDecode(ctx context.Context, req *http.Request, output any) error {
	start := time.Now()

	res, err := c.Do(req)
	if err != nil {
		return err
	}

	zerolog.Ctx(ctx).Debug().
		Dur("duration", time.Since(start)).
		Object("request", &httpRequestMarshalZerologObject{Request: req}).
		Object("response", &httpResponseMarshalZerologObject{Response: res}).
		Msg("outbound traffic")

	if res.StatusCode >= http.StatusBadRequest {
		defer DrainAndCloseResponse(res)
		return NewStatusCodeError(res.StatusCode)
	}

	defer DrainAndCloseResponse(res)

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
				Timeout:   30 * time.Second, //nolint:gomnd
				KeepAlive: 30 * time.Second, //nolint:gomnd
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,              //nolint:gomnd
			IdleConnTimeout:       90 * time.Second, //nolint:gomnd
			TLSHandshakeTimeout:   10 * time.Second, //nolint:gomnd
			ExpectContinueTimeout: 1 * time.Second,  //nolint:gomnd
		},
	}
}
