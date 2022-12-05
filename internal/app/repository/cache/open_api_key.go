package cache

import (
	"encoding/json"
	"fmt"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
)

type OpenApiKeyRepo struct {
	dbr repository.OpenApiKeyRepo
}

var _ repository.IOpenApiKeyRepo = new(OpenApiKeyRepo)

func (repo *OpenApiKeyRepo) FindByApiKey(apiKey string) (entity.OpenApiKey, error) {
	//fn := func() (interface{}, error) {
	//	return repo.dbr.FindByApiKey(apiKey)
	//}
	//
	//var res entity.OpenApiKey
	//err := rc.StringTemplateUnmarshal(fmt.Sprintf(openApiKey, apiKey), oneHour, fn, &res)
	//if err != nil {
	//	return res, err
	//}
	var res entity.OpenApiKey
	key := fmt.Sprintf(openApiKey, apiKey)
	get, err := rc.Get(key)
	if err == nil {
		_ = json.Unmarshal([]byte(get), &res)
		return res, nil
	}

	dbrRes, err := repo.dbr.FindByApiKey(apiKey)
	if err == nil {
		val, _ := json.Marshal(dbrRes)
		_ = rc.Set(key, val, oneHour)
		return dbrRes, nil
	}

	return res, err
}
