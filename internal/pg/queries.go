package pg

const (
	// SelectDatabaseOverview - основная статистика по БД (совместим с PG 10+)
	// Даёт «пульс» БД: транзакции, буферы, temp, deadlocks
	SelectDatabaseOverview = `
SELECT
  xact_commit,
  xact_rollback,
  blks_read,
  blks_hit,
  tup_returned,
  tup_fetched,
  tup_inserted,
  tup_updated,
  tup_deleted,
  conflicts,
  temp_files,
  temp_bytes,
  deadlocks,
  blk_read_time,   -- будет 0 если track_io_timing=off
  blk_write_time
FROM pg_stat_database
WHERE datname = $1;
`

	// SelectCacheHitRateGlobal - cache hit rate по всем БД (кроме template)
	SelectCacheHitRateGlobal = `
SELECT
  CASE WHEN sum(blks_hit + blks_read) = 0 THEN NULL
       ELSE sum(blks_hit)::float / NULLIF(sum(blks_hit + blks_read),0) END AS hit_rate
FROM pg_stat_database
WHERE datname NOT IN ('template0','template1');
`

	// SelectCacheHitRateDB - cache hit rate по конкретной БД
	SelectCacheHitRateDB = `
SELECT
  CASE WHEN (blks_hit + blks_read) = 0 THEN NULL
       ELSE blks_hit::float / NULLIF(blks_hit + blks_read,0) END AS hit_rate
FROM pg_stat_database
WHERE datname = $1;
`

	// SelectCheckpointsLegacy - статистика чекпоинтов для PG ≤16 (pg_stat_bgwriter)
	SelectCheckpointsLegacy = `
SELECT
  checkpoints_timed,
  checkpoints_req,
  checkpoint_write_time,
  checkpoint_sync_time,
  buffers_checkpoint,
  buffers_backend,
  buffers_backend_fsync,
  buffers_alloc
FROM pg_stat_bgwriter;
`

	// SelectCheckpointsV17 - статистика чекпоинтов для PG ≥17 (pg_stat_checkpointer)
	SelectCheckpointsV17 = `
SELECT
  checkpoints_timed,
  checkpoints_req,
  checkpoint_write_time,
  checkpoint_sync_time
FROM pg_stat_checkpointer;
`

	// SelectWalActivity - статистика WAL для PG ≥14
	// Очень информативно для диагностики WAL-нагрузки
	SelectWalActivity = `
SELECT
  wal_records,
  wal_fpi,
  wal_bytes,
  wal_buffers_full,
  stats_reset
FROM pg_stat_wal;
`

	// SelectTablesInfoLight - лёгкая статистика по таблицам
	// Показывает размер, сканы, мёртвые строки, последние VACUUM/ANALYZE
	// $1 - лимит (по умолчанию 200)
	SelectTablesInfoLight = `
SELECT
  schemaname,
  relname        AS table_name,
  pg_total_relation_size(relid) AS total_bytes,
  n_live_tup,
  n_dead_tup,
  seq_scan,
  idx_scan,
  last_vacuum,
  last_autovacuum,
  last_analyze,
  last_autoanalyze
FROM pg_stat_user_tables
WHERE (n_live_tup + COALESCE(seq_scan,0) + COALESCE(idx_scan,0)) > 0
ORDER BY total_bytes DESC
LIMIT COALESCE($1, 200);
`

	// SelectLockingNow - текущие блокировки и ожидания
	// Использует pg_blocking_pids для надёжного определения блокирующих процессов
	SelectLockingNow = `
SELECT
  a.pid,
  a.usename,
  a.datname,
  a.wait_event_type,
  a.wait_event,
  a.state,
  a.query_start,
  pg_blocking_pids(a.pid) AS blocking_pids
FROM pg_stat_activity a
WHERE a.datname = $1
  AND (a.wait_event IS NOT NULL OR cardinality(pg_blocking_pids(a.pid)) > 0)
ORDER BY a.query_start NULLS LAST;
`

	// SelectCurrentSettingsChanged - настройки, изменённые от дефолта
	// Полезно для быстрой диагностики конфигурации
	SelectCurrentSettingsChanged = `
SELECT
  name, 
  setting, 
  unit, 
  source,
  pending_restart
FROM pg_settings
WHERE source <> 'default'
ORDER BY name;
`
)
