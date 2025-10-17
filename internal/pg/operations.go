package pg

import (
	"context"
	"database/sql"
	"fmt"
)

// GetDatabaseOverview возвращает основную статистику по БД
func (c *Client) GetDatabaseOverview(ctx context.Context, dbName string) (*DatabaseOverview, error) {
	var stats DatabaseOverview

	err := c.db.QueryRowContext(ctx, SelectDatabaseOverview, dbName).Scan(
		&stats.XactCommit,
		&stats.XactRollback,
		&stats.BlksRead,
		&stats.BlksHit,
		&stats.TupReturned,
		&stats.TupFetched,
		&stats.TupInserted,
		&stats.TupUpdated,
		&stats.TupDeleted,
		&stats.Conflicts,
		&stats.TempFiles,
		&stats.TempBytes,
		&stats.Deadlocks,
		&stats.BlkReadTime,
		&stats.BlkWriteTime,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("database %s not found", dbName)
		}
		return nil, fmt.Errorf("failed to get database overview: %w", err)
	}

	return &stats, nil
}

// GetCacheHitRateGlobal возвращает общий cache hit rate по всем БД
func (c *Client) GetCacheHitRateGlobal(ctx context.Context) (*CacheHitRate, error) {
	var rate CacheHitRate

	err := c.db.QueryRowContext(ctx, SelectCacheHitRateGlobal).Scan(&rate.HitRate)
	if err != nil {
		return nil, fmt.Errorf("failed to get global cache hit rate: %w", err)
	}

	return &rate, nil
}

// GetCacheHitRateDB возвращает cache hit rate для конкретной БД
func (c *Client) GetCacheHitRateDB(ctx context.Context, dbName string) (*CacheHitRate, error) {
	var rate CacheHitRate

	err := c.db.QueryRowContext(ctx, SelectCacheHitRateDB, dbName).Scan(&rate.HitRate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("database %s not found", dbName)
		}
		return nil, fmt.Errorf("failed to get cache hit rate for database %s: %w", dbName, err)
	}

	return &rate, nil
}

// GetCheckpointsStats возвращает статистику чекпоинтов (version-aware)
func (c *Client) GetCheckpointsStats(ctx context.Context) (*CheckpointsStats, error) {
	var stats CheckpointsStats
	version := c.Version()

	if version == nil {
		return nil, fmt.Errorf("version not detected, call Connect() first")
	}

	// PG ≥17 использует pg_stat_checkpointer
	if version.SupportsCheckpointer() {
		err := c.db.QueryRowContext(ctx, SelectCheckpointsV17).Scan(
			&stats.CheckpointsTimed,
			&stats.CheckpointsReq,
			&stats.CheckpointWriteTime,
			&stats.CheckpointSyncTime,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get checkpoints stats (v17): %w", err)
		}
	} else {
		// PG ≤16 использует pg_stat_bgwriter
		err := c.db.QueryRowContext(ctx, SelectCheckpointsLegacy).Scan(
			&stats.CheckpointsTimed,
			&stats.CheckpointsReq,
			&stats.CheckpointWriteTime,
			&stats.CheckpointSyncTime,
			&stats.BuffersCheckpoint,
			&stats.BuffersBackend,
			&stats.BuffersBackendFsync,
			&stats.BuffersAlloc,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get checkpoints stats (legacy): %w", err)
		}
	}

	return &stats, nil
}

// GetWalActivity возвращает статистику WAL (только для PG ≥14)
func (c *Client) GetWalActivity(ctx context.Context) (*WalActivity, error) {
	version := c.Version()

	if version == nil {
		return nil, fmt.Errorf("version not detected, call Connect() first")
	}

	if !version.SupportsWalStats() {
		return nil, fmt.Errorf("WAL statistics not supported in PostgreSQL %d.%d (requires ≥14)",
			version.Major, version.Minor)
	}

	var stats WalActivity

	err := c.db.QueryRowContext(ctx, SelectWalActivity).Scan(
		&stats.WalRecords,
		&stats.WalFpi,
		&stats.WalBytes,
		&stats.WalBuffersFull,
		&stats.StatsReset,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get WAL activity: %w", err)
	}

	return &stats, nil
}

// GetTablesInfo возвращает статистику по таблицам
func (c *Client) GetTablesInfo(ctx context.Context, limit int) ([]TableInfo, error) {
	if limit <= 0 {
		limit = 200 // значение по умолчанию
	}

	rows, err := c.db.QueryContext(ctx, SelectTablesInfoLight, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables info: %w", err)
	}
	defer rows.Close()

	var tables []TableInfo

	for rows.Next() {
		var table TableInfo
		err := rows.Scan(
			&table.SchemaName,
			&table.TableName,
			&table.TotalBytes,
			&table.TableBytes,
			&table.IndexesBytes,
			&table.NLiveTup,
			&table.NDeadTup,
			&table.DeadRatio,
			&table.SeqScan,
			&table.IdxScan,
			&table.LastVacuum,
			&table.LastAutovacuum,
			&table.LastAnalyze,
			&table.LastAutoanalyze,
			&table.VacuumCount,
			&table.AutovacuumCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan table info: %w", err)
		}
		tables = append(tables, table)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tables: %w", err)
	}

	return tables, nil
}

// GetLockingInfo возвращает информацию о текущих блокировках
func (c *Client) GetLockingInfo(ctx context.Context, dbName string) ([]LockInfo, error) {
	rows, err := c.db.QueryContext(ctx, SelectLockingNow, dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to query locking info: %w", err)
	}
	defer rows.Close()

	var locks []LockInfo

	for rows.Next() {
		var lock LockInfo
		err := rows.Scan(
			&lock.PID,
			&lock.Username,
			&lock.Database,
			&lock.WaitEventType,
			&lock.WaitEvent,
			&lock.State,
			&lock.QueryStart,
			&lock.BlockingPids,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lock info: %w", err)
		}
		locks = append(locks, lock)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating locks: %w", err)
	}

	return locks, nil
}

// GetChangedSettings возвращает настройки, изменённые от дефолта
func (c *Client) GetChangedSettings(ctx context.Context) ([]SettingInfo, error) {
	rows, err := c.db.QueryContext(ctx, SelectCurrentSettingsChanged)
	if err != nil {
		return nil, fmt.Errorf("failed to query settings: %w", err)
	}
	defer rows.Close()

	var settings []SettingInfo

	for rows.Next() {
		var setting SettingInfo
		err := rows.Scan(
			&setting.Name,
			&setting.Setting,
			&setting.Unit,
			&setting.Source,
			&setting.PendingRestart,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan setting info: %w", err)
		}
		settings = append(settings, setting)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating settings: %w", err)
	}

	return settings, nil
}

// GetIndexStats возвращает статистику использования индексов
func (c *Client) GetIndexStats(ctx context.Context, limit int) ([]IndexStats, error) {
	if limit <= 0 {
		limit = 100 // значение по умолчанию
	}

	rows, err := c.db.QueryContext(ctx, SelectIndexStats, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query index stats: %w", err)
	}
	defer rows.Close()

	var indexes []IndexStats

	for rows.Next() {
		var index IndexStats
		err := rows.Scan(
			&index.SchemaName,
			&index.TableName,
			&index.IndexName,
			&index.IdxScan,
			&index.IdxTupRead,
			&index.IdxTupFetch,
			&index.SizeBytes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan index stats: %w", err)
		}
		indexes = append(indexes, index)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating indexes: %w", err)
	}

	return indexes, nil
}

// GetActiveQueries возвращает активные запросы с длительностью выше порога
func (c *Client) GetActiveQueries(ctx context.Context, dbName string, minDuration int) ([]ActiveQuery, error) {
	if minDuration <= 0 {
		minDuration = 5 // значение по умолчанию - 5 секунд
	}

	rows, err := c.db.QueryContext(ctx, SelectActiveQueries, dbName, minDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to query active queries: %w", err)
	}
	defer rows.Close()

	var queries []ActiveQuery

	for rows.Next() {
		var query ActiveQuery
		err := rows.Scan(
			&query.PID,
			&query.Username,
			&query.Database,
			&query.State,
			&query.DurationSeconds,
			&query.WaitEventType,
			&query.WaitEvent,
			&query.Query,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan active query: %w", err)
		}
		queries = append(queries, query)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating active queries: %w", err)
	}

	return queries, nil
}

// GetConnectionStats возвращает статистику соединений
func (c *Client) GetConnectionStats(ctx context.Context) (*ConnectionSummary, error) {
	var stats ConnectionSummary

	err := c.db.QueryRowContext(ctx, SelectConnectionStats).Scan(
		&stats.TotalConnections,
		&stats.Active,
		&stats.Idle,
		&stats.IdleInTransaction,
		&stats.Waiting,
		&stats.MaxConnections,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get connection stats: %w", err)
	}

	return &stats, nil
}

// GetSlowQueries возвращает топ медленных запросов из pg_stat_statements
// Требует установленного extension pg_stat_statements
func (c *Client) GetSlowQueries(ctx context.Context, limit int) ([]SlowQuery, error) {
	// Проверить наличие pg_stat_statements
	var exists bool
	err := c.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname = 'pg_stat_statements')").Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check pg_stat_statements: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("pg_stat_statements extension is not installed")
	}

	if limit <= 0 {
		limit = 20 // значение по умолчанию
	}

	rows, err := c.db.QueryContext(ctx, SelectSlowQueries, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query slow queries: %w", err)
	}
	defer rows.Close()

	var queries []SlowQuery

	for rows.Next() {
		var query SlowQuery
		err := rows.Scan(
			&query.Query,
			&query.Calls,
			&query.TotalExecTime,
			&query.MeanExecTime,
			&query.StddevExecTime,
			&query.Rows,
			&query.CacheHitPercent,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan slow query: %w", err)
		}
		queries = append(queries, query)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating slow queries: %w", err)
	}

	return queries, nil
}

// GetDatabaseSizes возвращает размеры всех баз данных
func (c *Client) GetDatabaseSizes(ctx context.Context) ([]DatabaseSize, error) {
	rows, err := c.db.QueryContext(ctx, SelectDatabaseSizes)
	if err != nil {
		return nil, fmt.Errorf("failed to query database sizes: %w", err)
	}
	defer rows.Close()

	var databases []DatabaseSize

	for rows.Next() {
		var db DatabaseSize
		err := rows.Scan(
			&db.DatabaseName,
			&db.SizeBytes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan database size: %w", err)
		}
		databases = append(databases, db)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating databases: %w", err)
	}

	return databases, nil
}
