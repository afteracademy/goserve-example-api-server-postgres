package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Helper function to capitalize the first letter of a string
func capitalizeFirstLetter(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(string(str[0])) + str[1:]
}

func generateFeature(featureTemplate string) error {
	if featureTemplate == "" {
		return errors.New("api name should be a non-empty string")
	}

	featureName := strings.ToLower(featureTemplate)
	featureDir := filepath.Join("api", featureName)
	if _, err := os.Stat(featureDir); err == nil {
		fmt.Println(featureName, "already exists")
		return nil
	}

	// Create api directory
	if err := os.MkdirAll(featureDir, os.ModePerm); err != nil {
		return err
	}

	if err := generateDto(featureDir, featureName); err != nil {
		return err
	}
	if err := generateModel(featureDir, featureName); err != nil {
		return err
	}
	if err := generateService(featureDir, featureName); err != nil {
		return err
	}
	if err := generateController(featureDir, featureName); err != nil {
		return err
	}
	return nil
}

func generateService(featureDir, featureName string) error {
	featureLower := strings.ToLower(featureName)
	featureCaps := capitalizeFirstLetter(featureName)
	servicePath := filepath.Join(featureDir, fmt.Sprintf("%sservice.go", ""))

	template := fmt.Sprintf(`package %s

import (
	"context"

  "github.com/afteracademy/goserve-example-api-server-postgres/api/%s/dto"
	"github.com/afteracademy/goserve-example-api-server-postgres/api/%s/model"
	"github.com/afteracademy/goserve/v2/redis"
	"github.com/afteracademy/goserve/v2/postgres"
	"github.com/google/uuid"
)

type Service interface {
	Find%s(id uuid.UUID) (*model.%s, error)
}

type service struct {
	db              postgres.Database
	info%sCache     redis.Cache[dto.Info%s]
}

func NewService(db postgres.Database, store redis.Store) Service {
	return &service{
	  db:              db,
		info%sCache:     redis.NewCache[dto.Info%s](store),
	}
}

func (s *service) Find%s(id uuid.UUID) (*model.%s, error) {
  ctx := context.Background()
	
	query := `+"`"+`
		SELECT
			id,
			field,
			status,
			created_at,
			updated_at
		FROM %ss
		WHERE id = $1
	`+"`"+`

	var m model.%s

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
`, featureLower, featureLower, featureCaps, featureCaps, featureCaps, featureCaps, featureCaps, featureCaps, featureCaps, featureCaps, featureCaps, featureLower, featureCaps)

	return os.WriteFile(servicePath, []byte(template), os.ModePerm)
}

func generateController(featureDir, featureName string) error {
	featureLower := strings.ToLower(featureName)
	featureCaps := capitalizeFirstLetter(featureName)
	controllerPath := filepath.Join(featureDir, fmt.Sprintf("%scontroller.go", ""))

	template := fmt.Sprintf(`package %s

import (
	"github.com/afteracademy/goserve-example-api-server-postgres/api/%s/dto"
	"github.com/afteracademy/goserve-example-api-server-postgres/common"
	coredto "github.com/afteracademy/goserve/v2/dto"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/afteracademy/goserve/v2/utility"
	"github.com/gin-gonic/gin"
)

type controller struct {
	network.Controller
	common.ContextPayload
	service Service
}

func NewController(
	authMFunc network.AuthenticationProvider,
	authorizeMFunc network.AuthorizationProvider,
	service Service,
) network.Controller {
	return &controller{
		Controller: network.NewController("/%s", authMFunc, authorizeMFunc),
		ContextPayload: common.NewContextPayload(),
		service:  service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/id/:id", c.get%sHandler)
}

func (c *controller) get%sHandler(ctx *gin.Context) {
	uuidParam, err := network.ReqParams[coredto.UUID](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	%s, err := c.service.Find%s(uuidParam.ID)
	if err != nil {
		network.SendNotFoundError(ctx, "%s not found", err)
		return
	}

	data, err := utility.MapTo[dto.Info%s](%s)
	if err != nil {
		network.SendInternalServerError(ctx, "something went wrong", err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", data)
}
`, featureLower, featureLower, featureLower, featureCaps, featureCaps, featureLower, featureCaps, featureLower, featureCaps, featureLower)

	return os.WriteFile(controllerPath, []byte(template), os.ModePerm)
}

func generateModel(featureDir, featureName string) error {
	modelDirPath := filepath.Join(featureDir, "model")
	if err := os.MkdirAll(modelDirPath, os.ModePerm); err != nil {
		return err
	}

	featureLower := strings.ToLower(featureName)
	featureCaps := capitalizeFirstLetter(featureName)
	modelPath := filepath.Join(featureDir, fmt.Sprintf("model/%s.go", featureLower))

	tStr := `package model

import (
	"time"

	"github.com/google/uuid"
)

type %s struct {
	ID        uuid.UUID  // id 
	Field     string     // field
	Status    bool       // status
	CreatedAt time.Time  // created_at
	UpdatedAt time.Time  // updated_at
}
`
	template := fmt.Sprintf(tStr, featureCaps)

	return os.WriteFile(modelPath, []byte(template), os.ModePerm)
}

func generateDto(featureDir, featureName string) error {
	dtoDirPath := filepath.Join(featureDir, "dto")
	if err := os.MkdirAll(dtoDirPath, os.ModePerm); err != nil {
		return err
	}

	featureLower := strings.ToLower(featureName)
	featureCaps := capitalizeFirstLetter(featureName)
	dtoPath := filepath.Join(featureDir, fmt.Sprintf("dto/create_%s.go", featureLower))

	tStr := `package dto

import (
	"time"

	"github.com/google/uuid"
)

type Info%s struct {
	ID        uuid.UUID ` + "`" + `json:"_id" binding:"required"` + "`" + `
	Field     string    ` + "`" + `json:"field" binding:"required"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"createdAt" binding:"required"` + "`" + `
}
`
	template := fmt.Sprintf(tStr, featureCaps)

	return os.WriteFile(dtoPath, []byte(template), os.ModePerm)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("api name should be non-empty string")
		return
	}

	featureName := os.Args[1]
	if err := generateFeature(featureName); err != nil {
		fmt.Println("Error:", err)
	}
}
