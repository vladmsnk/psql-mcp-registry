package pg

import (
	"database/sql"
)

// DatabaseOverview - основная статистика по БД
type DatabaseOverview struct {
	XactCommit   int64   `json:"xact_commit"`
	XactRollback int64   `json:"xact_rollback"`
	BlksRead     int64   `json:"blks_read"`
	BlksHit      int64   `json:"blks_hit"`
	TupReturned  int64   `json:"tup_returned"`
	TupFetched   int64   `json:"tup_fetched"`
	TupInserted  int64   `json:"tup_inserted"`
	TupUpdated   int64   `json:"tup_updated"`
	TupDeleted   int64   `json:"tup_deleted"`
	Conflicts    int64   `json:"conflicts"`
	TempFiles    int64   `json:"temp_files"`
	TempBytes    int64   `json:"temp_bytes"`
	Deadlocks    int64   `json:"deadlocks"`
	BlkReadTime  float64 `json:"blk_read_time"`
	BlkWriteTime float64 `json:"blk_write_time"`
}

// CacheHitRate - cache hit rate (процент попаданий в кэш)
type CacheHitRate struct {
	HitRate sql.NullFloat64 `json:"hit_rate"`
}

// CheckpointsStats - статистика чекпоинтов
type CheckpointsStats struct {
	CheckpointsTimed    int64         `json:"checkpoints_timed"`
	CheckpointsReq      int64         `json:"checkpoints_req"`
	CheckpointWriteTime float64       `json:"checkpoint_write_time"`
	CheckpointSyncTime  float64       `json:"checkpoint_sync_time"`
	BuffersCheckpoint   sql.NullInt64 `json:"buffers_checkpoint,omitempty"`    // только в legacy
	BuffersBackend      sql.NullInt64 `json:"buffers_backend,omitempty"`       // только в legacy
	BuffersBackendFsync sql.NullInt64 `json:"buffers_backend_fsync,omitempty"` // только в legacy
	BuffersAlloc        sql.NullInt64 `json:"buffers_alloc,omitempty"`         // только в legacy
}

// WalActivity - статистика WAL (доступна только в PG ≥14)
type WalActivity struct {
	WalRecords     int64        `json:"wal_records"`
	WalFpi         int64        `json:"wal_fpi"`
	WalBytes       int64        `json:"wal_bytes"`
	WalBuffersFull int64        `json:"wal_buffers_full"`
	StatsReset     sql.NullTime `json:"stats_reset"`
}

// TableInfo - статистика по таблице
type TableInfo struct {
	SchemaName      string          `json:"schema_name"`
	TableName       string          `json:"table_name"`
	TotalBytes      int64           `json:"total_bytes"`
	TableBytes      int64           `json:"table_bytes"`
	IndexesBytes    int64           `json:"indexes_bytes"`
	NLiveTup        sql.NullInt64   `json:"n_live_tup"`
	NDeadTup        sql.NullInt64   `json:"n_dead_tup"`
	DeadRatio       sql.NullFloat64 `json:"dead_ratio"`
	SeqScan         sql.NullInt64   `json:"seq_scan"`
	IdxScan         sql.NullInt64   `json:"idx_scan"`
	LastVacuum      sql.NullTime    `json:"last_vacuum"`
	LastAutovacuum  sql.NullTime    `json:"last_autovacuum"`
	LastAnalyze     sql.NullTime    `json:"last_analyze"`
	LastAutoanalyze sql.NullTime    `json:"last_autoanalyze"`
	VacuumCount     sql.NullInt64   `json:"vacuum_count"`
	AutovacuumCount sql.NullInt64   `json:"autovacuum_count"`
}

// LockInfo - информация о блокировках
type LockInfo struct {
	PID           int            `json:"pid"`
	Username      sql.NullString `json:"username"`
	Database      sql.NullString `json:"database"`
	WaitEventType sql.NullString `json:"wait_event_type"`
	WaitEvent     sql.NullString `json:"wait_event"`
	State         sql.NullString `json:"state"`
	QueryStart    sql.NullTime   `json:"query_start"`
	BlockingPids  string         `json:"blocking_pids"` // массив как строка, например "{123,456}"
}

// SettingInfo - информация о настройке PostgreSQL
type SettingInfo struct {
	Name           string         `json:"name"`
	Setting        string         `json:"setting"`
	Unit           sql.NullString `json:"unit"`
	Source         string         `json:"source"`
	PendingRestart bool           `json:"pending_restart"`
}

// Version - информация о версии PostgreSQL
type Version struct {
	Major      int    `json:"major"`       // основная версия (например, 14, 17)
	Minor      int    `json:"minor"`       // минорная версия
	Patch      int    `json:"patch"`       // патч версия
	FullString string `json:"full_string"` // полная строка версии
}

// SupportsWalStats проверяет, поддерживает ли версия pg_stat_wal (PG ≥14)
func (v *Version) SupportsWalStats() bool {
	return v.Major >= 14
}

// SupportsCheckpointer проверяет, поддерживает ли версия pg_stat_checkpointer (PG ≥17)
func (v *Version) SupportsCheckpointer() bool {
	return v.Major >= 17
}

// SupportsBlockingPids проверяет, поддерживает ли версия pg_blocking_pids (PG ≥9.6)
func (v *Version) SupportsBlockingPids() bool {
	return v.Major >= 10 || (v.Major == 9 && v.Minor >= 6)
}

// IndexStats - статистика по индексу
type IndexStats struct {
	SchemaName  string        `json:"schema_name"`
	TableName   string        `json:"table_name"`
	IndexName   string        `json:"index_name"`
	IdxScan     sql.NullInt64 `json:"idx_scan"`
	IdxTupRead  sql.NullInt64 `json:"idx_tup_read"`
	IdxTupFetch sql.NullInt64 `json:"idx_tup_fetch"`
	SizeBytes   int64         `json:"size_bytes"`
}

// ActiveQuery - информация об активном запросе
type ActiveQuery struct {
	PID             int             `json:"pid"`
	Username        sql.NullString  `json:"username"`
	Database        sql.NullString  `json:"database"`
	State           sql.NullString  `json:"state"`
	DurationSeconds sql.NullFloat64 `json:"duration_seconds"`
	WaitEventType   sql.NullString  `json:"wait_event_type"`
	WaitEvent       sql.NullString  `json:"wait_event"`
	Query           sql.NullString  `json:"query"`
}

// ConnectionSummary - сводная статистика соединений
type ConnectionSummary struct {
	TotalConnections  int `json:"total_connections"`
	Active            int `json:"active"`
	Idle              int `json:"idle"`
	IdleInTransaction int `json:"idle_in_transaction"`
	Waiting           int `json:"waiting"`
	MaxConnections    int `json:"max_connections"`
}

// SlowQuery - информация о медленном запросе из pg_stat_statements
type SlowQuery struct {
	Query           string          `json:"query"`
	Calls           int64           `json:"calls"`
	TotalExecTime   float64         `json:"total_exec_time"`
	MeanExecTime    float64         `json:"mean_exec_time"`
	StddevExecTime  float64         `json:"stddev_exec_time"`
	Rows            int64           `json:"rows"`
	CacheHitPercent sql.NullFloat64 `json:"cache_hit_percent"`
}

// DatabaseSize - информация о размере базы данных
type DatabaseSize struct {
	DatabaseName string `json:"database_name"`
	SizeBytes    int64  `json:"size_bytes"`
}
