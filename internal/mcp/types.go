package mcp

type DatabaseOverviewInput struct {
	InstanceName string `json:"instance_name" jsonschema:"name of the PostgreSQL instance,required"`
	DbName       string `json:"db_name,omitempty" jsonschema:"database name (default: postgres)"`
}
type CacheHitRateInput struct {
	InstanceName string `json:"instance_name" jsonschema:"name of the PostgreSQL instance,required"`
	DbName       string `json:"db_name,omitempty" jsonschema:"database name (optional, if not specified returns global stats)"`
}
type CheckpointsStatsInput struct {
	InstanceName string `json:"instance_name" jsonschema:"name of the PostgreSQL instance,required"`
}
type WalActivityInput struct {
	InstanceName string `json:"instance_name" jsonschema:"name of the PostgreSQL instance,required"`
}
type TablesInfoInput struct {
	InstanceName string `json:"instance_name" jsonschema:"name of the PostgreSQL instance,required"`
	Limit        int    `json:"limit,omitempty" jsonschema:"maximum number of tables to return (default: 200)"`
}
type LockingInfoInput struct {
	InstanceName string `json:"instance_name" jsonschema:"name of the PostgreSQL instance,required"`
	DbName       string `json:"db_name,omitempty" jsonschema:"database name (default: postgres)"`
}
type ChangedSettingsInput struct {
	InstanceName string `json:"instance_name" jsonschema:"name of the PostgreSQL instance,required"`
}
type VersionInput struct {
	InstanceName string `json:"instance_name" jsonschema:"name of the PostgreSQL instance,required"`
}
type ListInstancesInput struct {
}
type IndexStatsInput struct {
	InstanceName string `json:"instance_name" jsonschema:"name of the PostgreSQL instance,required"`
	Limit        int    `json:"limit,omitempty" jsonschema:"maximum number of indexes to return (default: 100)"`
}
type ActiveQueriesInput struct {
	InstanceName       string `json:"instance_name" jsonschema:"name of the PostgreSQL instance,required"`
	DbName             string `json:"db_name,omitempty" jsonschema:"database name (default: postgres)"`
	MinDurationSeconds int    `json:"min_duration_seconds,omitempty" jsonschema:"minimum query duration in seconds (default: 5)"`
}
type ConnectionStatsInput struct {
	InstanceName string `json:"instance_name" jsonschema:"name of the PostgreSQL instance,required"`
}
type SlowQueriesInput struct {
	InstanceName string `json:"instance_name" jsonschema:"name of the PostgreSQL instance,required"`
	Limit        int    `json:"limit,omitempty" jsonschema:"maximum number of slow queries to return (default: 20)"`
}
type DatabaseSizesInput struct {
	InstanceName string `json:"instance_name" jsonschema:"name of the PostgreSQL instance,required"`
}
