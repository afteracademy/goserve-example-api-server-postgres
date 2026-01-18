package startup

import (
	"net/http/httptest"

	"github.com/afteracademy/goserve-example-api-server-postgres/config"
	"github.com/afteracademy/goserve/v2/network"
)

type Teardown = func()

func TestServer() (network.Router, Module, Teardown) {
	env := config.NewEnv("../.test.env", false)
	router, module, shutdown := create(env)
	ts := httptest.NewServer(router.GetEngine())
	teardown := func() {
		ts.Close()
		shutdown()
	}
	return router, module, teardown
}
