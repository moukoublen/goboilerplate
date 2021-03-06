package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPI_AboutRouteHandler(t *testing.T) {
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// http call
	AboutHandler(resp, req)

	// verify
	assert.Equal(t, http.StatusOK, resp.Code)
	responseBody := map[string]interface{}{}
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Contains(t, responseBody, "version")
	assert.Contains(t, responseBody, "branch")
	assert.Contains(t, responseBody, "commit")
	assert.Contains(t, responseBody, "commit_short")
	assert.Contains(t, responseBody, "tag")
}
