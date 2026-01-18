package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afteracademy/goserve-example-api-server-postgres/api/auth/model"
	roleModel "github.com/afteracademy/goserve-example-api-server-postgres/api/user/model"
	"github.com/afteracademy/goserve-example-api-server-postgres/startup"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/afteracademy/goserve/v2/utility"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationAuthController_SignupSuccess(t *testing.T) {
	router, module, shutdown := startup.TestServer()
	var role *roleModel.Role
	var apikey *model.ApiKey
	defer shutdown()

	t.Cleanup(func() {
		if apikey != nil {
			module.GetInstance().AuthService.DeleteApiKey(apikey)
		}
	})

	t.Cleanup(func() {
		if role != nil {
			module.GetInstance().UserService.DeleteRole(role)
		}
	})

	t.Cleanup(func() {
		module.GetInstance().UserService.RemoveUserByEmail("test@abc.com")
	})

	key, err := utility.GenerateRandomString(6)
	if err != nil {
		t.Fatalf("could not create key: %v", err)
	}

	apikey, err = module.GetInstance().AuthService.CreateApiKey(key, 1, []model.Permission{"test"}, []string{"comment"})
	if err != nil {
		t.Fatalf("could not create apikey: %v", err)
	}

	role, err = module.GetInstance().UserService.CreateRole(roleModel.RoleCodeLearner)
	if err != nil {
		t.Fatalf("could not create role: %v", err)
	}

	body := `{"email":"test@abc.com","password":"123456","name":"test name"}`

	req, err := http.NewRequest("POST", "/auth/signup/basic", bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add(network.ApiKeyHeader, apikey.Key)

	rr := httptest.NewRecorder()
	router.GetEngine().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"success"`)
	assert.Contains(t, rr.Body.String(), `"data"`)
	assert.Contains(t, rr.Body.String(), `"user"`)
	assert.Contains(t, rr.Body.String(), `"roles"`)
	assert.Contains(t, rr.Body.String(), `"tokens"`)

}
