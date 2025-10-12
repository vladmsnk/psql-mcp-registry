package router

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"psql-mcp-registry/internal/model"
	"psql-mcp-registry/internal/pg"
	pgmocks "psql-mcp-registry/internal/pg/mocks"
	routermocks "psql-mcp-registry/internal/router/mocks"
)

func TestRouter_RouteQuery_DatabaseOverview_Success(t *testing.T) {
	ctx := context.Background()

	instance := model.Instance{
		Name:         "test-instance",
		DatabaseName: "testdb",
		Description:  "Test PostgreSQL instance",
		Status:       "active",
	}

	expectedOverview := &pg.DatabaseOverview{
		XactCommit:   1000,
		XactRollback: 10,
		BlksRead:     5000,
		BlksHit:      95000,
		TupReturned:  100000,
		TupFetched:   80000,
		TupInserted:  1000,
		TupUpdated:   500,
		TupDeleted:   100,
		Conflicts:    0,
		TempFiles:    5,
		TempBytes:    1024000,
		Deadlocks:    0,
		BlkReadTime:  100.5,
		BlkWriteTime: 50.3,
	}

	// Create mocks
	mockClient := pgmocks.NewClientInterface(t)
	mockRegistry := routermocks.NewRegistry(t)

	// Set expectations
	mockRegistry.On("GetInstanceClient", instance).Return(mockClient)
	mockClient.On("GetDatabaseOverview", ctx, "postgres").Return(expectedOverview, nil)

	router := New(mockRegistry)

	req := QueryRequest{
		InstanceName: instance.Name,
		Action:       model.ActionNameDatabaseOverview,
		Parameters: map[string]interface{}{
			"dbName": "postgres",
		},
	}

	response, err := router.RouteQuery(ctx, req, instance)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Success)
	assert.Equal(t, instance.Name, response.Instance)
	assert.Equal(t, model.ActionNameDatabaseOverview, response.Action)
	assert.Empty(t, response.Error)

	overview, ok := response.Data.(*pg.DatabaseOverview)
	assert.True(t, ok, "response.Data should be *pg.DatabaseOverview")
	assert.Equal(t, expectedOverview.XactCommit, overview.XactCommit)
	assert.Equal(t, expectedOverview.BlksRead, overview.BlksRead)
	assert.Equal(t, expectedOverview.BlksHit, overview.BlksHit)
}
