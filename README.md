# goserve - Go Backend Architecture

[![Docker Compose CI](https://github.com/afteracademy/goserve-example-api-server-postgres/actions/workflows/docker_compose.yml/badge.svg)](https://github.com/afteracademy/goserve-example-api-server-postgres/actions/workflows/docker_compose.yml)
[![Architechture](https://img.shields.io/badge/Framework-blue?label=View&logo=go)](https://github.com/afteracademy/goserve)
[![Starter Project](https://img.shields.io/badge/Starter%20Project%20CLI-red?label=Get&logo=go)](https://github.com/afteracademy/goservegen)
[![Download](https://img.shields.io/badge/Download-Starter%20Project%20Postgres%20Zip-green.svg)](https://github.com/afteracademy/goservegen/raw/main/starter-project-postgres.zip)

![Banner](.extra/docs/goserve-banner.png)

## Create A Blog Service 

This project is a fully production-ready solution designed to implement best practices for building performant and secure backend REST API services. It provides a robust architectural framework to ensure consistency and maintain high code quality. The architecture emphasizes feature separation, facilitating easier unit and integration testing. It is built using the `goserve` framework, which offers essential functionalities such as authentication, authorization, database connectivity, and caching.

## Framework Libs
- Go
- goserve v2
- Gin
- jwt
- postgres
- pgx
- go-redis
- Validator
- Viper
- Crypto

**Highlights**
- REST API design
- goserve framework usage
- API key support
- Token based Authentication
- Role based Authorization
- Unit Tests
- Integration Tests
- Modular codebase

## Architecture
The goal is to make each API independent from one another and only share services among them. This will make code reusable and reduce conflicts while working in a team. 

The APIs will have separate directory based on the endpoint. Example `blog` and `blogs` will have seperate directory whereas `blog`, `blog/author`, and `blog/editor` will share common resources and will live inside same directory.

## Know More on goserve framework
- [GitHub - afteracademy/goserve](https://github.com/afteracademy/goserve) 

### Startup Flow
cmd/main → startup/server → module, postgres, redis, router → api/[feature]/middlewares → api/[feature]/controller -> api/[feature]/service, authentication, authorization → handlers → sender

### API Structure
```
Sample API
├── dto
│   └── create_sample.go
├── model
│   └── sample.go
├── controller.go
└── service.go
```

- Each feature API lives under `api` directory
- The request and response body is sent in the form of a DTO (Data Transfer Object) inside `dto` directory
- The database collection model lives inside `model` directory
- Controller is responsible for defining endpoints and corresponding handlers
- Service is the main logic component and handles data. Controller interact with a service to process a request. A service can also interact with other services.
 
## Project Directories
1. **api**: APIs code 
3. **cmd**: main function to start the program
4. **common**: code to be used in all the apis
5. **config**: load environment variables
6. **keys**: stores server pem files for token
7. **startup**: creates server and initializes database, redis, and router
8. **tests**: holds the integration tests
9. **utils**: contains utility functions

**Helper/Optional Directories**
1. **.extra**: postgres sql scripts for initialization inside docker, other web assets and documents
2. **.github**: CI for tests
3. **.tools**: api code, RSA key generator, and .env copier
4. **.vscode**: editor config and debug launch settings

## API Design
![Request-Response-Design](.extra/docs/request-flow.svg)

### API DOC
[![API Documentation](https://img.shields.io/badge/API%20Documentation-View%20Here-blue?style=for-the-badge)](https://documenter.getpostman.com/view/1552895/2sBXVihVLg)

## Installation Instructions
vscode is the recommended editor - dark theme 

**1. Get the repo**

```bash
git clone https://github.com/afteracademy/goserve-example-api-server-postgres.git
```

**2. Generate RSA Keys**
```
go run .tools/rsa/keygen.go
```

**3. Create .env files**
```
go run .tools/copy/envs.go 
```

**4. Run Docker Compose**
- Install Docker and Docker Compose. [Find Instructions Here](https://docs.docker.com/install/).

```bash
docker compose up --build
```
-  You will be able to access the api from http://localhost:8080

**5. Run Tests**
```bash
docker exec -t goserve_example_api_server_postgres go test -v ./...
```

If having any issue
- Make sure 8080 port is not occupied else change SERVER_PORT in **.env** file.
- Make sure 5432 port is not occupied else change DB_PORT in **.env** file.
- Make sure 6379 port is not occupied else change REDIS_PORT in **.env** file.

## Run on the local machine
```bash
go mod tidy
```

Keep the docker container for `postgres` and `redis` running and **stop** the `goserve_example_api_server_postgres` docker container

Change the following hosts in the **.env** and **.test.env**
- DB_HOST=localhost
- REDIS_HOST=localhost

Best way to run this project is to use the vscode `Run and Debug` button. Scripts are available for debugging and template generation on vscode.

### Optional - Running the app from terminal
```bash
go run cmd/main.go
```

## Template
New api creation can be done using command. `go run .tools/apigen.go [feature_name]`. This will create all the required skeleton files inside the directory api/[feature_name]

```bash
go run .tools/apigen.go sample
```

## Read the Article to understand this project
[How to Architect Good Go Backend REST API Services](https://afteracademy.com/article/how-to-architect-good-go-backend-rest-api-services)

## How to use this architecture in your project?
You can use [goservegen](https://github.com/afteracademy/goservegen) CLI to generate starter project for this architecture. 
> Check out the repo [github.com/afteracademy/goservegen](https://github.com/afteracademy/goservegen) for more information.

## Documentation
Information about the framework

### Model
`api/sample/model/sample.go`

```go
package model

import (
	"time"

	"github.com/google/uuid"
)

type Sample struct {
	ID        uuid.UUID  // id 
	Field     string     // field
	Status    bool       // status
	CreatedAt time.Time  // created_at
	UpdatedAt time.Time  // updated_at
}
```

### DTO
`api/sample/dto/create_sample.go`

```go
package dto

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/afteracademy/goserve/v2/utility"
)

type InfoSample struct {
	ID        uuid.UUID `json:"_id" binding:"required"`
	Field     string    `json:"field" binding:"required"`
	CreatedAt time.Time `json:"createdAt" binding:"required"`
}

func EmptyInfoSample() *InfoSample {
	return &InfoSample{}
}

func (d *InfoSample) GetValue() *InfoSample {
	return d
}

func (d *InfoSample) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utility.FormatValidationErrors(errs), nil
}
```

#### Notes: The DTO implements the interface 
`arch/network/interfaces.go`

```golang
type Dto[T any] interface {
  GetValue() *T
  ValidateErrors(errs validator.ValidationErrors) ([]string, error)
}
``` 

### Service
`api/sample/service.go`

```go
package sample

import (
	"context"

  "github.com/afteracademy/goserve-example-api-server-postgres/api/sample/dto"
	"github.com/afteracademy/goserve-example-api-server-postgres/api/Sample/model"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/afteracademy/goserve/v2/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"
)

type Service interface {
	FindSample(id uuid.UUID) (*model.Sample, error)
}

type service struct {
	network.BaseService
	db              *pgxpool.Pool
	infoSampleCache     redis.Cache[dto.InfoSample]
}

func NewService(db *pgxpool.Pool, store redis.Store) Service {
	return &service{
		BaseService:     network.NewBaseService(),
	  db:              db,
		infoSampleCache:     redis.NewCache[dto.InfoSample](store),
	}
}

func (s *service) FindSample(id uuid.UUID) (*model.Sample, error) {
  ctx := context.Background()
	
	query := `
		SELECT
			id,
			field,
			status,
			created_at,
			updated_at
		FROM samples
		WHERE id = $1
	`

	var m model.Sample

	err := s.db.QueryRow(ctx, query, id).
		Scan(
			&m.ID,
			&m.Field,
			&m.Status,
			&m.CreatedAt,
			&m.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &m, nil
}
```

#### Notes: The Service embeds the interface 
`github.com/afteracademy/goserve/v2/network/interfaces.go`

```golang
type BaseService interface {
  Context() context.Context
}
``` 

- Redis Cache: `redis.Cache[dto.InfoSample]` provide the methods to make common redis queries for the DTO `dto.InfoSample`

### Controller
`api/sample/controller.go`

```go
package sample

import (
	"github.com/afteracademy/goserve-example-api-server-postgres/api/sample/dto"
	"github.com/afteracademy/goserve-example-api-server-postgres/common"
	coredto "github.com/afteracademy/goserve/v2/dto"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/afteracademy/goserve/v2/utility"
	"github.com/gin-gonic/gin"
)

type controller struct {
	network.BaseController
	common.ContextPayload
	service Service
}

func NewController(
	authMFunc network.AuthenticationProvider,
	authorizeMFunc network.AuthorizationProvider,
	service Service,
) network.Controller {
	return &controller{
		BaseController: network.NewBaseController("/sample", authMFunc, authorizeMFunc),
		ContextPayload: common.NewContextPayload(),
		service:  service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/id/:id", c.getSampleHandler)
}

func (c *controller) getSampleHandler(ctx *gin.Context) {
	uuidParam, err := network.ReqParams(ctx, coredto.EmptyUUID())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	sample, err := c.service.FindSample(uuidParam.ID)
	if err != nil {
		c.Send(ctx).NotFoundError("sample not found", err)
		return
	}

	data, err := utility.MapTo[dto.InfoSample](sample)
	if err != nil {
		c.Send(ctx).InternalServerError("something went wrong", err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", data)
}
```

#### Notes: The Controller implements the interface 
`github.com/afteracademy/goserve/v2/network/interfaces.go`

```golang
type Controller interface {
  BaseController
  MountRoutes(group *gin.RouterGroup)
}

type BaseController interface {
  ResponseSender
  Path() string
  Authentication() gin.HandlerFunc
  Authorization(role string) gin.HandlerFunc
}

type ResponseSender interface {
  Debug() bool
  Send(ctx *gin.Context) SendResponse
}

type SendResponse interface {
  SuccessMsgResponse(message string)
  SuccessDataResponse(message string, data any)
  BadRequestError(message string, err error)
  ForbiddenError(message string, err error)
  UnauthorizedError(message string, err error)
  NotFoundError(message string, err error)
  InternalServerError(message string, err error)
  MixedError(err error)
}
``` 

### Enable Controller In Module
`startup/module.go`

```go
import (
  ...
  "github.com/afteracademy/goserve-example-api-server-postgres/api/sample"
)

...

func (m *module) Controllers() []network.Controller {
  return []network.Controller{
    ...
    sample.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), sample.NewService(m.DB, m.Store)),
  }
}
```

## Go Microservices Architecture using goserve
`goserve` also provides `micro` package to build REST API microservices. Find the microservices version of this blog service project at [github.com/afteracademy/gomicro](https://github.com/afteracademy/gomicro)

[Article - How to Create Microservices — A Practical Guide Using Go](https://afteracademy.com/article/how-to-create-microservices-a-practical-guide-using-go)

## Find this project useful ? :heart:
* Support it by clicking the :star: button on the upper right of this page. :v:

## More on YouTube channel - AfterAcademy
Subscribe to the YouTube channel `AfterAcademy` for understanding the concepts used in this project:

[![YouTube](https://img.shields.io/badge/YouTube-Subscribe-red?style=for-the-badge&logo=youtube&logoColor=white)](https://www.youtube.com/@afteracad)

## Contribution
Please feel free to fork it and open a PR.