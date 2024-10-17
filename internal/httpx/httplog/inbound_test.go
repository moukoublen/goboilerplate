package httplog

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/moukoublen/goboilerplate/pkg/testingx"
)

type serverTestCase struct {
	httpClient *http.Client
	testSrv    *httptest.Server
	logOutput  *bytes.Buffer
}

func (s *serverTestCase) Init(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) {
	t.Helper()
	s.logOutput = &bytes.Buffer{}
	logger := slog.New(slog.NewJSONHandler(s.logOutput, nil))
	il := NewHTTPLogger(
		WithLogger(logger),
		WithLogInLevel(slog.LevelInfo),
	)

	s.testSrv = httptest.NewServer(il.Handler(http.HandlerFunc(handler)))

	s.httpClient = NewHTTPClient(2 * time.Second)

	t.Cleanup(s.Close)
}

func (s *serverTestCase) Do(ctx context.Context, t *testing.T, method string, body io.Reader, assertReq func(*http.Request), assertRes func(*http.Response)) {
	t.Helper()

	req, err := http.NewRequestWithContext(ctx, method, s.testSrv.URL, body)
	if err != nil {
		t.Fatalf("error while creating http request: %s", err.Error())
	}

	req.Header.Set(`Content-Type`, `application/json`)
	res, err := s.httpClient.Do(req) //nolint:bodyclose // the asserter should close the body.
	if err != nil {
		t.Fatalf("error while performing http request: %s", err.Error())
	}
	if assertReq != nil {
		assertReq(req)
	}
	if assertRes != nil {
		assertRes(res)
	}
}

func (s *serverTestCase) Close() {
	s.testSrv.Close()
}

func TestInbound(t *testing.T) {
	srv := serverTestCase{}
	srv.Init(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody := testingx.ReadAndClose(t, r.Body)
		testingx.AssertEqual(t, requestBody, []byte(`"request body"`))

		w.Header().Set(`Content-Type`, `application/json`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`"response body"`))
	}))

	srv.Do(
		context.Background(),
		t,
		http.MethodGet,
		strings.NewReader(`"request body"`),
		func(_ *http.Request) {},
		func(r *http.Response) {
			body := testingx.ReadAndClose(t, r.Body)
			testingx.AssertEqual(t, body, []byte(`"response body"`))
		},
	)

	logsBytes := srv.logOutput.Bytes()
	logsObject := map[string]any{}
	err := json.Unmarshal(logsBytes, &logsObject)
	if err != nil {
		t.Fatal(err)
	}

	expectedLogs := map[string]any{
		// "duration": 44116,
		// "time":     "2024-10-05T14:40:01.323900808+03:00",
		"level": "INFO",
		"msg":   "http inbound",
		"request": map[string]any{
			// "host":          "127.0.0.1:39625",
			// "remoteAddr":    "127.0.0.1:57194",
			"body":          map[string]any{"size": float64(14), "value": `\"request body\"`},
			"contentLength": float64(14),
			"headers":       map[string]any{"Accept-Encoding": "gzip", "Content-Length": "14", "Content-Type": "application/json", "User-Agent": "Go-http-client/1.1"},
			"method":        "GET",
			"proto":         "HTTP/1.1",
			"requestUri":    "/",
			"url":           map[string]any{"fragment": "", "full": ":///", "host": "", "opaque": "", "path": "/", "scheme": ""},
		},
		"response": map[string]any{
			"body":    map[string]any{"size": float64(15), "value": `\"response body\"`},
			"headers": map[string]any{"Content-Type": "application/json"},
			"status":  map[string]any{"code": float64(200), "name": "OK"},
		},
	}

	testingx.AssertPartialEqualMap(t, logsObject, expectedLogs)
}
