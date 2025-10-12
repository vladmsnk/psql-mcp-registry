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
			&table.NLiveTup,
			&table.NDeadTup,
			&table.SeqScan,
			&table.IdxScan,
			&table.LastVacuum,
			&table.LastAutovacuum,
			&table.LastAnalyze,
			&table.LastAutoanalyze,
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
