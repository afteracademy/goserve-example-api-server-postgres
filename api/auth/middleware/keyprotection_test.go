package middleware

import (
	"errors"
	"net/http"
	"testing"

	"github.com/afteracademy/goserve-example-api-server-postgres/api/auth"
	"github.com/afteracademy/goserve-example-api-server-postgres/api/auth/model"
	"github.com/afteracademy/goserve-example-api-server-postgres/common"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestKeyProtectionMiddleware_NoApiKey(t *testing.T) {
	mockAuthService := new(auth.MockService)

	rr := network.MockTestRootMiddleware(
		t,
		NewKeyProtection(mockAuthService),
		network.MockSuccessMsgHandler("success"),
		nil,
	)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"permission denied: missing x-api-key header"`)
}

func TestKeyProtectionMiddleware_WrongApiKey(t *testing.T) {
	mockAuthService := new(auth.MockService)
	key := "wrong"
	mockAuthService.On("FetchApiKey", key).Return(nil, errors.New(""))

	rr := network.MockTestRootMiddleware(
		t,
		NewKeyProtection(mockAuthService),
		network.MockSuccessMsgHandler("success"),
		map[string]string{network.ApiKeyHeader: key},
	)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"permission denied: invalid x-api-key"`)
}

func TestKeyProtectionMiddleware_CorrectApiKey(t *testing.T) {
	mockAuthService := new(auth.MockService)
	key := "correct"
	mockAuthService.On("FetchApiKey", key).Return(&model.ApiKey{Key: key}, nil)

	mockHandler := func(ctx *gin.Context) {
		assert.Equal(t, common.NewContextPayload().MustGetApiKey(ctx).Key, key)
		network.SendSuccessMsgResponse(ctx, "success")
	}

	rr := network.MockTestRootMiddleware(
		t,
		NewKeyProtection(mockAuthService),
		mockHandler,
		map[string]string{network.ApiKeyHeader: key},
	)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"success"`)
}
