package model

type ActionName string

var (
	ActionNameDatabaseOverview ActionName = "databases_overview"
	ActionNameCacheHitRate     ActionName = "cache_hit_rate"
	ActionNameCheckpointsStats ActionName = "checkpoints_stats"
	ActionNameWalActivity      ActionName = "wal_activity"
	ActionNameTablesInfo       ActionName = "tables_info"
	ActionNameLockingInfo      ActionName = "locking_info"
	ActionNameChangedSettings  ActionName = "changed_settings"
	ActionNameVersion          ActionName = "version"
	ActionNameIndexStats       ActionName = "index_stats"
	ActionNameActiveQueries    ActionName = "active_queries"
	ActionNameConnectionStats  ActionName = "connection_stats"
	ActionNameSlowQueries      ActionName = "slow_queries"
	ActionNameDatabaseSizes    ActionName = "database_sizes"
)
