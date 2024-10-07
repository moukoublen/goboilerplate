package httplog

import (
	"net"
	"net/http"
	"time"
)

// NewHTTPClient returns a new default http client.
func NewHTTPClient(globalTimeout time.Duration) *http.Client {
	// http.DefaultClient
	// http.DefaultTransport
	// http.defaultTransportDialContext
	return &http.Client{
		Timeout: globalTimeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

// func TestOutbound(t *testing.T) {
// 	ml := NewMockSlogHandler(slog.NewJSONHandler(os.Stdout, nil), slog.LevelDebug)

// 	il := NewHTTPLogger(
// 		WithLogger(slog.New(ml)),
// 		WithLogInLevel(slog.LevelInfo),
// 	)

// 	client := NewHTTPClient(10 * time.Second)
// 	client.Transport = il.LoggerRoundTripper(client.Transport)

// 	_, _ = client.Get("https://api.myip.com")
// }
