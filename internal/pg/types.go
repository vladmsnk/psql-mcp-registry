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
	SchemaName      string        `json:"schema_name"`
	TableName       string        `json:"table_name"`
	TotalBytes      int64         `json:"total_bytes"`
	NLiveTup        sql.NullInt64 `json:"n_live_tup"`
	NDeadTup        sql.NullInt64 `json:"n_dead_tup"`
	SeqScan         sql.NullInt64 `json:"seq_scan"`
	IdxScan         sql.NullInt64 `json:"idx_scan"`
	LastVacuum      sql.NullTime  `json:"last_vacuum"`
	LastAutovacuum  sql.NullTime  `json:"last_autovacuum"`
	LastAnalyze     sql.NullTime  `json:"last_analyze"`
	LastAutoanalyze sql.NullTime  `json:"last_autoanalyze"`
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
