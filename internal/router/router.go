package router

import (
	"context"
	"fmt"

	"psql-mcp-registry/internal/model"
	"psql-mcp-registry/internal/pg"
)

type Router struct {
	registry Registry
}

//go:generate mockery --case snake --name Registry
type Registry interface {
	AddInstanceToRegistry(instance model.Instance) error
	GetInstanceClient(instance model.Instance) pg.ClientInterface
}

func New(registry Registry) *Router {
	return &Router{registry: registry}
}

func (r *Router) RouteQuery(ctx context.Context, req QueryRequest, instance model.Instance) (*QueryResponse, error) {
	client := r.registry.GetInstanceClient(instance)
	if client == nil {
		return &QueryResponse{
			Instance: instance.Name,
			Action:   req.Action,
			Success:  false,
			Error:    fmt.Sprintf("client not found for instance: %s", instance.Name),
		}, fmt.Errorf("client not found for instance: %s", instance.Name)
	}

	response := &QueryResponse{
		Instance: instance.Name,
		Action:   req.Action,
		Success:  false,
	}

	var err error
	var data interface{}

	switch req.Action {
	case model.ActionNameDatabaseOverview:
		dbName := getStringParam(req.Parameters, "dbName", "postgres")
		data, err = client.GetDatabaseOverview(ctx, dbName)

	case model.ActionNameCacheHitRate:
		if dbName, exists := req.Parameters["dbName"]; exists && dbName != "" {
			data, err = client.GetCacheHitRateDB(ctx, dbName.(string))
		} else {
			data, err = client.GetCacheHitRateGlobal(ctx)
		}

	case model.ActionNameCheckpointsStats:
		data, err = client.GetCheckpointsStats(ctx)

	case model.ActionNameWalActivity:
		data, err = client.GetWalActivity(ctx)

	case model.ActionNameTablesInfo:
		limit := getIntParam(req.Parameters, "limit", 200)
		data, err = client.GetTablesInfo(ctx, limit)

	case model.ActionNameLockingInfo:
		dbName := getStringParam(req.Parameters, "dbName", "postgres")
		data, err = client.GetLockingInfo(ctx, dbName)

	case model.ActionNameChangedSettings:
		data, err = client.GetChangedSettings(ctx)

	case model.ActionNameVersion:
		version := client.Version()
		if version == nil {
			err = fmt.Errorf("version not detected, call Connect() first")
		} else {
			data = version
		}

	default:
		err = fmt.Errorf("unsupported action: %s", req.Action)
	}

	if err != nil {
		response.Error = err.Error()
		return response, err
	}

	response.Success = true
	response.Data = data
	return response, nil
}

func getStringParam(params map[string]interface{}, key, defaultValue string) string {
	if val, exists := params[key]; exists {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return defaultValue
}

func getIntParam(params map[string]interface{}, key string, defaultValue int) int {
	if val, exists := params[key]; exists {
		switch v := val.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case float64:
			return int(v)
		}
	}
	return defaultValue
}
