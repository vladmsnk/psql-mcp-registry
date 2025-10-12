package router

import "psql-mcp-registry/internal/model"

type QueryRequest struct {
	InstanceName string                 `json:"instance_name"`
	Action       model.ActionName       `json:"action"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
}

type QueryResponse struct {
	Instance string           `json:"instance"`
	Action   model.ActionName `json:"action"`
	Success  bool             `json:"success"`
	Data     interface{}      `json:"data,omitempty"`
	Error    string           `json:"error,omitempty"`
}
