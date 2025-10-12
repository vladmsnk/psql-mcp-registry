package pg

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
)

// DetectVersion определяет версию PostgreSQL
func DetectVersion(ctx context.Context, db *sql.DB) (*Version, error) {
	var versionString string
	err := db.QueryRowContext(ctx, "SELECT version()").Scan(&versionString)
	if err != nil {
		return nil, fmt.Errorf("failed to query PostgreSQL version: %w", err)
	}

	version, err := parseVersion(versionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PostgreSQL version: %w", err)
	}

	version.FullString = versionString
	return version, nil
}

// parseVersion парсит строку версии PostgreSQL
// Примеры:
// - "PostgreSQL 14.5 (Debian 14.5-1.pgdg110+1) on x86_64-pc-linux-gnu..."
// - "PostgreSQL 17.0 on x86_64-apple-darwin23.6.0..."
func parseVersion(versionString string) (*Version, error) {
	// Регулярное выражение для извлечения версии
	re := regexp.MustCompile(`PostgreSQL (\d+)\.(\d+)(?:\.(\d+))?`)
	matches := re.FindStringSubmatch(versionString)

	if len(matches) < 3 {
		return nil, fmt.Errorf("unable to parse version from: %s", versionString)
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", matches[1])
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", matches[2])
	}

	patch := 0
	if len(matches) >= 4 && matches[3] != "" {
		patch, err = strconv.Atoi(matches[3])
		if err != nil {
			return nil, fmt.Errorf("invalid patch version: %s", matches[3])
		}
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

// GetServerVersionNum возвращает числовое значение версии сервера
// Это альтернативный способ определения версии через SHOW server_version_num
func GetServerVersionNum(ctx context.Context, db *sql.DB) (int, error) {
	var versionNum int
	err := db.QueryRowContext(ctx, "SHOW server_version_num").Scan(&versionNum)
	if err != nil {
		return 0, fmt.Errorf("failed to get server_version_num: %w", err)
	}
	return versionNum, nil
}
