package mcp

type DatabaseOverviewInput struct {
	InstanceName string `json:"instance_name" jsonschema:"required,description=Name of the PostgreSQL instance"`
	DbName       string `json:"db_name,omitempty" jsonschema:"description=Database name (default: postgres)"`
}
type CacheHitRateInput struct {
	InstanceName string `json:"instance_name" jsonschema:"required,description=Name of the PostgreSQL instance"`
	DbName       string `json:"db_name,omitempty" jsonschema:"description=Database name (optional, if not specified returns global stats)"`
}
type CheckpointsStatsInput struct {
	InstanceName string `json:"instance_name" jsonschema:"required,description=Name of the PostgreSQL instance"`
}
type WalActivityInput struct {
	InstanceName string `json:"instance_name" jsonschema:"required,description=Name of the PostgreSQL instance"`
}
type TablesInfoInput struct {
	InstanceName string `json:"instance_name" jsonschema:"required,description=Name of the PostgreSQL instance"`
	Limit        int    `json:"limit,omitempty" jsonschema:"description=Maximum number of tables to return (default: 200)"`
}
type LockingInfoInput struct {
	InstanceName string `json:"instance_name" jsonschema:"required,description=Name of the PostgreSQL instance"`
	DbName       string `json:"db_name,omitempty" jsonschema:"description=Database name (default: postgres)"`
}
type ChangedSettingsInput struct {
	InstanceName string `json:"instance_name" jsonschema:"required,description=Name of the PostgreSQL instance"`
}
type VersionInput struct {
	InstanceName string `json:"instance_name" jsonschema:"required,description=Name of the PostgreSQL instance"`
}
type ListInstancesInput struct {
}
