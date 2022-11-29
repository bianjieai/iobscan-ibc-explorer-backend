package cache

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
)

type OpenApiKeyRepo struct {
	dbr repository.OpenApiKeyRepo
}

var _ repository.IOpenApiKeyRepo = new(OpenApiKeyRepo)

func (repo *OpenApiKeyRepo) FindByApiKey(apiKey string) (*entity.OpenApiKey, error) {
	fn := func() (interface{}, error) {
		return repo.dbr.FindByApiKey(apiKey)
	}

	var res entity.OpenApiKey
	err := rc.StringTemplateUnmarshal(fmt.Sprintf(openApiKey, apiKey), oneHour, fn, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
