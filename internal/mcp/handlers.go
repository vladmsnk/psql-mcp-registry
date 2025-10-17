package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"psql-mcp-registry/internal/model"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *MCPServer) handleDatabaseOverview(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input DatabaseOverviewInput,
) (*mcp.CallToolResult, interface{}, error) {
	params := make(map[string]interface{})
	if input.DbName != "" {
		params["dbName"] = input.DbName
	}

	data, err := s.executeRouterQuery(ctx, input.InstanceName, model.ActionNameDatabaseOverview, params)
	if err != nil {
		return nil, nil, err
	}

	return nil, data, nil
}

func (s *MCPServer) handleCacheHitRate(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input CacheHitRateInput,
) (*mcp.CallToolResult, interface{}, error) {
	params := make(map[string]interface{})
	if input.DbName != "" {
		params["dbName"] = input.DbName
	}

	data, err := s.executeRouterQuery(ctx, input.InstanceName, model.ActionNameCacheHitRate, params)
	if err != nil {
		return nil, nil, err
	}

	return nil, data, nil
}

func (s *MCPServer) handleCheckpointsStats(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input CheckpointsStatsInput,
) (*mcp.CallToolResult, interface{}, error) {
	data, err := s.executeRouterQuery(ctx, input.InstanceName, model.ActionNameCheckpointsStats, nil)
	if err != nil {
		return nil, nil, err
	}

	return nil, data, nil
}

func (s *MCPServer) handleWalActivity(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input WalActivityInput,
) (*mcp.CallToolResult, interface{}, error) {
	data, err := s.executeRouterQuery(ctx, input.InstanceName, model.ActionNameWalActivity, nil)
	if err != nil {
		return nil, nil, err
	}

	return nil, data, nil
}

func (s *MCPServer) handleTablesInfo(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input TablesInfoInput,
) (*mcp.CallToolResult, interface{}, error) {
	params := make(map[string]interface{})
	if input.Limit > 0 {
		params["limit"] = input.Limit
	}

	data, err := s.executeRouterQuery(ctx, input.InstanceName, model.ActionNameTablesInfo, params)
	if err != nil {
		return nil, nil, err
	}

	return nil, data, nil
}

func (s *MCPServer) handleLockingInfo(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input LockingInfoInput,
) (*mcp.CallToolResult, interface{}, error) {
	params := make(map[string]interface{})
	if input.DbName != "" {
		params["dbName"] = input.DbName
	}

	data, err := s.executeRouterQuery(ctx, input.InstanceName, model.ActionNameLockingInfo, params)
	if err != nil {
		return nil, nil, err
	}

	return nil, data, nil
}

func (s *MCPServer) handleChangedSettings(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ChangedSettingsInput,
) (*mcp.CallToolResult, interface{}, error) {
	data, err := s.executeRouterQuery(ctx, input.InstanceName, model.ActionNameChangedSettings, nil)
	if err != nil {
		return nil, nil, err
	}

	return nil, data, nil
}

func (s *MCPServer) handleVersion(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input VersionInput,
) (*mcp.CallToolResult, interface{}, error) {
	data, err := s.executeRouterQuery(ctx, input.InstanceName, model.ActionNameVersion, nil)
	if err != nil {
		return nil, nil, err
	}

	return nil, data, nil
}

func (s *MCPServer) handleListInstancesResource(
	ctx context.Context,
	req *mcp.ReadResourceRequest,
) (*mcp.ReadResourceResult, error) {
	instances, err := s.manager.ListInstances(ctx)
	if err != nil {
		return nil, err
	}
	var instanceData []map[string]interface{}
	for _, inst := range instances {
		instanceData = append(instanceData, map[string]interface{}{
			"name":          inst.Name,
			"database_name": inst.DatabaseName,
			"description":   inst.Description,
			"status":        inst.Status,
			"created_at":    inst.CreatedAt,
			"updated_at":    inst.UpdatedAt,
		})
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "instances://list",
				MIMEType: "application/json",
				Text:     formatJSON(instanceData),
			},
		},
	}, nil
}

func formatJSON(data interface{}) string {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "failed to marshal JSON: %s"}`, err.Error())
	}
	return string(jsonBytes)
}

func (s *MCPServer) handleIndexStats(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input IndexStatsInput,
) (*mcp.CallToolResult, interface{}, error) {
	params := make(map[string]interface{})
	if input.Limit > 0 {
		params["limit"] = input.Limit
	}

	data, err := s.executeRouterQuery(ctx, input.InstanceName, model.ActionNameIndexStats, params)
	if err != nil {
		return nil, nil, err
	}

	return nil, data, nil
}

func (s *MCPServer) handleActiveQueries(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ActiveQueriesInput,
) (*mcp.CallToolResult, interface{}, error) {
	params := make(map[string]interface{})
	if input.DbName != "" {
		params["dbName"] = input.DbName
	}
	if input.MinDurationSeconds > 0 {
		params["minDuration"] = input.MinDurationSeconds
	}

	data, err := s.executeRouterQuery(ctx, input.InstanceName, model.ActionNameActiveQueries, params)
	if err != nil {
		return nil, nil, err
	}

	return nil, data, nil
}

func (s *MCPServer) handleConnectionStats(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ConnectionStatsInput,
) (*mcp.CallToolResult, interface{}, error) {
	data, err := s.executeRouterQuery(ctx, input.InstanceName, model.ActionNameConnectionStats, nil)
	if err != nil {
		return nil, nil, err
	}

	return nil, data, nil
}

func (s *MCPServer) handleSlowQueries(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SlowQueriesInput,
) (*mcp.CallToolResult, interface{}, error) {
	params := make(map[string]interface{})
	if input.Limit > 0 {
		params["limit"] = input.Limit
	}

	data, err := s.executeRouterQuery(ctx, input.InstanceName, model.ActionNameSlowQueries, params)
	if err != nil {
		return nil, nil, err
	}

	return nil, data, nil
}

func (s *MCPServer) handleDatabaseSizes(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input DatabaseSizesInput,
) (*mcp.CallToolResult, interface{}, error) {
	data, err := s.executeRouterQuery(ctx, input.InstanceName, model.ActionNameDatabaseSizes, nil)
	if err != nil {
		return nil, nil, err
	}

	return nil, data, nil
}
