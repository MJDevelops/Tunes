package services

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type ImportService struct {
	ctx context.Context
	db  *DbService
}

func NewImportService(dbService *DbService) *ImportService {
	return &ImportService{db: dbService}
}

func (s *ImportService) ServiceStartup(ctx context.Context, option application.ServiceOptions) error {
	s.ctx = ctx
	return nil
}
