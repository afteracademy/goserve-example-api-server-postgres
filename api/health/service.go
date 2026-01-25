package health

import (
	"time"

	"github.com/afteracademy/goserve-example-api-server-postgres/api/health/dto"
)

type Service interface {
	CheckHealth() (*dto.HealthCheck, error)
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (s *service) CheckHealth() (*dto.HealthCheck, error) {
	health := &dto.HealthCheck{
		Timestamp: time.Now(),
		Status:    "OK",
	}
	return health, nil
}
