package main

// This migration tool preserves historical Python analytics data while moving
// the production backend to the Go SQLite schema.
import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type accessLogRow struct {
	ID           int64
	IPAddress    string
	VisitorID    string
	UserAgent    string
	Path         string
	Method       string
	StatusCode   int64
	ResponseTime float64
	Referer      string
	DeviceType   string
	OS           string
	Browser      string
	Country      string
	Region       string
	CreatedAt    string
}

func main() {
	var (
		sourcePath string
		targetPath string
		backupDir  string
	)

	flag.StringVar(&sourcePath, "source", "", "Python analytics.db path")
	flag.StringVar(&targetPath, "target", "", "Go mdxy.db path")
	flag.StringVar(&backupDir, "backup-dir", "", "backup directory (default: alongside target)")
	flag.Parse()

	if sourcePath == "" || targetPath == "" {
		flag.Usage()
		os.Exit(2)
	}

	if err := migrate(sourcePath, targetPath, backupDir); err != nil {
		log.Fatalf("迁移失败: %v", err)
	}
}

func migrate(sourcePath, targetPath, backupDir string) error {
	sourcePath, err := filepath.Abs(sourcePath)
	if err != nil {
		return err
	}
	targetPath, err = filepath.Abs(targetPath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(sourcePath); err != nil {
		return fmt.Errorf("源数据库不存在: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return err
	}

	if backupDir == "" {
		backupDir = filepath.Join(filepath.Dir(targetPath), "migration-backups")
	}
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102-150405")
	if err := backupFile(sourcePath, filepath.Join(backupDir, "python-analytics-"+timestamp+".db")); err != nil {
		return err
	}
	if _, err := os.Stat(targetPath); err == nil {
		if err := backupFile(targetPath, filepath.Join(backupDir, "go-target-"+timestamp+".db")); err != nil {
			return err
		}
	}

	sourceDB, err := sql.Open("sqlite3", sourcePath)
	if err != nil {
		return err
	}
	defer sourceDB.Close()

	targetDB, err := sql.Open("sqlite3", targetPath)
	if err != nil {
		return err
	}
	defer targetDB.Close()

	if err := ensureAccessLogsSchema(targetDB); err != nil {
		return err
	}

	sourceCount, sourceMin, sourceMax, err := statsFor(sourceDB)
	if err != nil {
		return fmt.Errorf("读取源库统计失败: %w", err)
	}

	hasCountry, err := hasColumn(sourceDB, "access_logs", "country")
	if err != nil {
		return err
	}
	hasRegion, err := hasColumn(sourceDB, "access_logs", "region")
	if err != nil {
		return err
	}

	imported, err := copyAccessLogs(sourceDB, targetDB, hasCountry, hasRegion)
	if err != nil {
		return err
	}

	targetCount, targetMin, targetMax, err := statsFor(targetDB)
	if err != nil {
		return fmt.Errorf("读取目标库统计失败: %w", err)
	}

	log.Printf("源库 access_logs: count=%d, min=%s, max=%s", sourceCount, sourceMin, sourceMax)
	log.Printf("导入完成: inserted=%d", imported)
	log.Printf("目标库 access_logs: count=%d, min=%s, max=%s", targetCount, targetMin, targetMax)
	log.Printf("备份目录: %s", backupDir)
	return nil
}

func ensureAccessLogsSchema(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS access_logs (
			id integer PRIMARY KEY AUTOINCREMENT,
			ip_address text,
			visitor_id text,
			user_agent text,
			path text,
			method text,
			status_code integer,
			response_time real,
			referer text,
			device_type text,
			os text,
			browser text,
			country text,
			region text,
			created_at datetime
		)`,
		`CREATE INDEX IF NOT EXISTS idx_access_logs_created_at ON access_logs(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_access_logs_path ON access_logs(path)`,
		`CREATE INDEX IF NOT EXISTS idx_access_logs_visitor_id ON access_logs(visitor_id)`,
		`CREATE INDEX IF NOT EXISTS idx_access_logs_ip_address ON access_logs(ip_address)`,
	}

	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}

	for _, col := range []struct {
		name string
		typ  string
	}{
		{name: "country", typ: "TEXT"},
		{name: "region", typ: "TEXT"},
	} {
		exists, err := hasColumn(db, "access_logs", col.name)
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		if _, err := db.Exec(fmt.Sprintf("ALTER TABLE access_logs ADD COLUMN %s %s", col.name, col.typ)); err != nil {
			return err
		}
	}

	return nil
}

func copyAccessLogs(sourceDB, targetDB *sql.DB, hasCountry, hasRegion bool) (int64, error) {
	cols := []string{
		"id", "ip_address", "visitor_id", "user_agent", "path", "method",
		"status_code", "response_time", "referer", "device_type", "os",
		"browser",
	}
	if hasCountry {
		cols = append(cols, "country")
	} else {
		cols = append(cols, "'' AS country")
	}
	if hasRegion {
		cols = append(cols, "region")
	} else {
		cols = append(cols, "'' AS region")
	}
	cols = append(cols, "created_at")

	query := fmt.Sprintf("SELECT %s FROM access_logs ORDER BY id ASC", strings.Join(cols, ", "))
	rows, err := sourceDB.Query(query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	tx, err := targetDB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO access_logs (
			id, ip_address, visitor_id, user_agent, path, method,
			status_code, response_time, referer, device_type, os,
			browser, country, region, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var imported int64
	for rows.Next() {
		var row accessLogRow
		if err := rows.Scan(
			&row.ID,
			&row.IPAddress,
			&row.VisitorID,
			&row.UserAgent,
			&row.Path,
			&row.Method,
			&row.StatusCode,
			&row.ResponseTime,
			&row.Referer,
			&row.DeviceType,
			&row.OS,
			&row.Browser,
			&row.Country,
			&row.Region,
			&row.CreatedAt,
		); err != nil {
			return 0, err
		}

		result, err := stmt.Exec(
			row.ID,
			row.IPAddress,
			row.VisitorID,
			row.UserAgent,
			row.Path,
			row.Method,
			row.StatusCode,
			row.ResponseTime,
			row.Referer,
			row.DeviceType,
			row.OS,
			row.Browser,
			row.Country,
			row.Region,
			row.CreatedAt,
		)
		if err != nil {
			return 0, err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		imported += affected
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}

	if _, err := tx.Exec(`UPDATE sqlite_sequence SET seq = (SELECT COALESCE(MAX(id), 0) FROM access_logs) WHERE name = 'access_logs'`); err != nil {
		// sqlite_sequence 可能还不存在，忽略即可
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return imported, nil
}

func statsFor(db *sql.DB) (count int64, minCreatedAt, maxCreatedAt string, err error) {
	row := db.QueryRow(`SELECT COUNT(*), COALESCE(MIN(created_at), ''), COALESCE(MAX(created_at), '') FROM access_logs`)
	err = row.Scan(&count, &minCreatedAt, &maxCreatedAt)
	return
}

func hasColumn(db *sql.DB, tableName, columnName string) (bool, error) {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cid        int
			name       string
			colType    string
			notNull    int
			defaultV   sql.NullString
			primaryKey int
		)
		if err := rows.Scan(&cid, &name, &colType, &notNull, &defaultV, &primaryKey); err != nil {
			return false, err
		}
		if name == columnName {
			return true, nil
		}
	}
	return false, rows.Err()
}

func backupFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
