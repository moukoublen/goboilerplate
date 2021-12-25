package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPI_AboutRouteHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	AboutRouteHandler(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	responseBody := map[string]interface{}{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Contains(t, responseBody, "version")
	assert.Contains(t, responseBody, "branch")
	assert.Contains(t, responseBody, "commit")
	assert.Contains(t, responseBody, "commit_short")
	assert.Contains(t, responseBody, "tag")
}
